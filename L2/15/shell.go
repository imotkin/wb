package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const (
	prompt = "> "

	OperatorSuccess  = "&&"
	OperatorFailed   = "||"
	OperatorPipeline = "|"
	OperatorTo       = ">"
	OperatorFrom     = "<"
)

func closer(c io.Closer) error {
	if c != nil {
		return c.Close()
	}

	return nil
}

type Operator string

type Stream struct {
	in  io.ReadCloser
	out io.WriteCloser
}

func (s Stream) Read(p []byte) (n int, err error) {
	return s.in.Read(p)
}

func (s Stream) Write(p []byte) (n int, err error) {
	return s.out.Write(p)
}

func (s Stream) Close() error {
	return errors.Join(closer(s.in), closer(s.out))
}

func (s Stream) hasIn() bool {
	return s.in != nil
}

func (s Stream) hasOut() bool {
	return s.out != nil
}

func (s Stream) isEmpty() bool {
	return s.in == nil && s.out == nil
}

type Shell struct {
	w       io.Writer
	builtin map[string]func([]string, Stream) error
}

func NewShell() (*Shell, error) {
	shell := &Shell{w: os.Stdout}

	shell.builtin = map[string]func([]string, Stream) error{
		"cd":   shell.cd,
		"pwd":  shell.pwd,
		"echo": shell.echo,
		"kill": shell.kill,
		"ps":   shell.ps,
	}

	return shell, nil
}

func (s *Shell) Run() {
	scan := bufio.NewScanner(os.Stdin)

	for {
		if !scan.Scan() {
			break
		}

		line := scan.Text()
		if line == "exit" {
			break
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		s.Execute(ctx, os.ExpandEnv(line))
		cancel()
	}
}

func (s *Shell) Execute(ctx context.Context, line string) bool {
	idx, op := s.getOperator(line)

	if idx != -1 {
		left := line[:idx]
		right := line[idx+len(op):]

		ok := s.Execute(ctx, left)

		switch op {
		case OperatorSuccess:
			if ok {
				return s.Execute(ctx, right)
			}
			return ok
		case OperatorFailed:
			if !ok {
				return s.Execute(ctx, right)
			}
			return ok
		}
	}

	err := s.exec(ctx, line)

	if err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			s.println()
		} else if _, ok := errors.AsType[*exec.ExitError](err); !ok {
			s.printf("mini-shell: %v\n", err)
		}

		return false
	}

	return true
}

func (s *Shell) println(a ...any) {
	fmt.Fprintln(os.Stdout, a...)
}

func (s *Shell) printf(format string, a ...any) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func (s *Shell) print(a ...any) {
	fmt.Fprint(os.Stdout, a...)
}

func (s *Shell) exec(ctx context.Context, line string) error {
	line = strings.TrimSpace(line)

	if strings.Contains(line, OperatorPipeline) {
		return s.execPipeline(ctx, line)
	}

	return s.execCmd(ctx, line)
}

func (s *Shell) execPipeline(ctx context.Context, line string) error {
	parts := strings.Split(line, OperatorPipeline)
	cmds := make([]*exec.Cmd, len(parts))

	var streams []Stream

	defer func() {
		for _, s := range streams {
			s.Close()
		}
	}()

	for i, p := range parts {
		fields := strings.Fields(p)

		args, stream, err := s.parseRedirect(fields)
		if err != nil {
			return err
		}

		cmds[i] = exec.CommandContext(ctx, args[0], args[1:]...)

		cmds[i].Stderr = os.Stderr

		if stream.hasIn() {
			cmds[i].Stdin = stream.in
		}

		if stream.hasOut() {
			cmds[i].Stdout = stream.out
		} else if i == (len(parts) - 1) {
			cmds[i].Stdout = os.Stdout
		}

		if !stream.isEmpty() {
			streams = append(streams, stream)
		}
	}

	for i := 0; i < len(cmds)-1; i++ {
		if cmds[i+1].Stdin == nil {
			out, _ := cmds[i].StdoutPipe()
			cmds[i+1].Stdin = out
		}
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	var errs []error

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (s *Shell) pwd(_ []string, stream Stream) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	if stream.hasOut() {
		fmt.Fprintln(stream, dir)
	} else {
		s.println(dir)
	}

	return nil
}

func (s *Shell) cd(args []string, _ Stream) error {
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		return os.Chdir(home)
	}

	return os.Chdir(args[0])
}

func (s *Shell) ps(_ []string, stream Stream) error {
	cmd := exec.Command("ps")

	if stream.hasOut() {
		cmd.Stdout = stream.out
	} else {
		cmd.Stdout = s.w
	}

	return cmd.Run()
}

func (s *Shell) echo(args []string, stream Stream) error {
	message := strings.Trim(strings.Join(args, " "), `"'`)

	if stream.hasOut() {
		fmt.Fprintln(stream, message)
	} else {
		s.println(message)
	}

	return nil
}

func (s *Shell) kill(args []string, _ Stream) error {
	if len(args) == 0 {
		return errors.New("kill: not enough arguments")
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("kill: invalid pid: %v", args)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("kill %s failed: no such process", args)
	}

	return process.Signal(syscall.SIGTERM)
}

func (s *Shell) execBuiltin(args []string, stream Stream) error {
	return s.builtin[args[0]](args[1:], stream)
}

func (s *Shell) execCmd(ctx context.Context, text string) error {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return nil
	}

	args, stream, err := s.parseRedirect(parts)
	if err != nil {
		return err
	}
	defer stream.Close()

	if _, ok := s.builtin[args[0]]; ok {
		return s.execBuiltin(args, stream)
	}

	command := exec.CommandContext(ctx, args[0], args[1:]...)

	if stream.hasIn() {
		command.Stdin = stream.in
	}
	if stream.hasOut() {
		command.Stdout = stream.out
	} else {
		command.Stdout = s.w
	}

	command.Stderr = os.Stderr

	return command.Run()
}

func (s *Shell) parseRedirect(parts []string) (args []string, stream Stream, err error) {
	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case OperatorTo:
			if i+1 < len(parts) {
				stream.out, err = os.OpenFile(parts[i+1], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return nil, Stream{}, err
				}
				i++
			}
		case OperatorFrom:
			if i+1 < len(parts) {
				stream.in, err = os.Open(parts[i+1])
				if err != nil {
					return nil, Stream{}, err
				}
				i++
			}
		default:
			args = append(args, strings.Trim(parts[i], `'"`))
		}
	}

	return args, stream, nil
}

func (s *Shell) getOperator(line string) (int, Operator) {
	var (
		quote = false
		idx   = -1
		op    Operator
	)

	for i, r := range line {
		if r == '"' {
			quote = !quote
			continue
		}
		if quote {
			continue
		}

		if strings.HasPrefix(line[i:], string(OperatorSuccess)) {
			idx = i
			op = OperatorSuccess
		} else if strings.HasPrefix(line[i:], string(OperatorFailed)) {
			idx = i
			op = OperatorFailed
		}
	}

	return idx, op
}

func main() {
	shell, err := NewShell()
	if err != nil {
		log.Fatal(err)
	}

	shell.Run()
}
