// Задача L2.12 - Утилита grep

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	MessageFileNotFound = "grep: %s: No such file or directory\n"

	delimiter = "--"
)

type Matcher struct {
	n    int
	fn   func(string) Match
	opts Options
}

func NewMatcher(query string, opts Options) *Matcher {
	var re *regexp.Regexp

	if opts.LikeString {
		if opts.IgnoreCase {
			query = strings.ToLower(query)
		}
	} else {
		pattern := query
		if opts.IgnoreCase {
			pattern = "(?i)" + query
		}

		var err error

		re, err = regexp.Compile(pattern)
		if err != nil {
			opts.LikeString = true
			if opts.IgnoreCase {
				query = strings.ToLower(query)
			}
		}
	}

	fn := func(line string) Match {
		var found bool
		var index = -1

		if opts.LikeString {
			original := line
			if opts.IgnoreCase {
				original = strings.ToLower(line)
			}
			index = strings.Index(original, query)
			found = index != -1
		} else {
			loc := re.FindStringIndex(line)
			if loc != nil {
				index = loc[0]
				found = true
			}
		}

		return Match{Line: line, Index: index, Found: found}
	}

	return &Matcher{fn: fn, opts: opts}
}

func (m *Matcher) Match(line string) Match {
	m.n++
	match := m.fn(line)
	match.LineN = m.n
	return match
}

type Match struct {
	LineN int
	Line  string
	Index int
	Found bool
}

type Options struct {
	LinesBefore int
	LinesAfter  int
	LinesAround int
	OnlyNumber  bool
	IgnoreCase  bool
	Invert      bool
	LikeString  bool
	LineNumber  bool
}

type Grep struct {
	m    *Matcher
	opts Options
}

func NewGrep(query string, opts Options) Grep {
	return Grep{
		m:    NewMatcher(query, opts),
		opts: opts,
	}
}

func (g Grep) ContainsIndex(query, s string) (int, bool) {
	var contains bool
	var index int

	if g.opts.LikeString {
		index = strings.Index(s, query)
		contains = true
	} else {
		re, err := regexp.Compile(query)
		if err != nil {
			index = strings.Index(s, query)
			contains = true
		} else {
			loc := re.FindStringIndex(s)
			if loc == nil {
				index = -1
			} else {
				index = loc[0]
				contains = true
			}
		}
	}

	if g.opts.Invert {
		return index, !contains
	}

	return index, contains
}

func (g Grep) PrintCount(r io.Reader) {
	scan := bufio.NewScanner(r)

	var total int64

	for scan.Scan() {
		line := scan.Text()

		match := g.m.Match(line)

		isCounted := (match.Found && !g.opts.Invert) ||
			(!match.Found && g.opts.Invert)

		if isCounted {
			total++
		}
	}

	fmt.Println(total)
}

func (g Grep) Run(r io.Reader) {
	scan := bufio.NewScanner(r)

	if g.opts.OnlyNumber {
		g.PrintCount(r)
		return
	}

	var (
		before []Match

		afterCount = 0
		lastLine   = 0
	)

	if g.opts.LinesAround > 0 {
		g.opts.LinesAfter = g.opts.LinesAround
		g.opts.LinesBefore = g.opts.LinesAround
	}

	for scan.Scan() {
		line := scan.Text()

		match := g.m.Match(line)

		if g.opts.Invert {
			match.Found = !match.Found
		}

		if match.Found {
			first := match.LineN

			if len(before) > 0 {
				first = before[0].LineN
			}

			printDelim := lastLine > 0 && first > lastLine+1 &&
				(g.opts.LinesAfter > 0 || g.opts.LinesBefore > 0)

			if printDelim {
				fmt.Println(delimiter)
			}

			for _, m := range before {
				g.Print(m)
			}

			before = nil

			g.Print(match)
			lastLine = match.LineN

			afterCount = g.opts.LinesAfter
		} else {
			if afterCount > 0 {
				g.Print(match)
				lastLine = match.LineN
				afterCount--
			} else if g.opts.LinesBefore > 0 {
				before = append(before, match)
				if len(before) > g.opts.LinesBefore {
					before = before[1:]
				}
			}
		}
	}
}

func (g Grep) Print(m Match) {
	if g.opts.LineNumber {
		if m.Found {
			fmt.Printf("%d:%s\n", m.LineN, m.Line)
		} else {
			fmt.Printf("%d-%s\n", m.LineN, m.Line)
		}
	} else {
		fmt.Println(m.Line)
	}
}

func parseOptions() Options {
	var opts Options

	flag.IntVar(&opts.LinesAfter, "A", 0, "Number of lines after")
	flag.IntVar(&opts.LinesBefore, "B", 0, "Number of lines before")
	flag.IntVar(&opts.LinesAround, "C", 0, "Number of lines around")
	flag.BoolVar(&opts.OnlyNumber, "c", false, "Print only line number")
	flag.BoolVar(&opts.IgnoreCase, "i", false, "")
	flag.BoolVar(&opts.Invert, "v", false, "")
	flag.BoolVar(&opts.LikeString, "F", false, "")
	flag.BoolVar(&opts.LineNumber, "n", false, "")

	flag.Parse()

	return opts
}

func main() {
	opts := parseOptions()

	var r io.Reader
	var query string

	if filename := flag.Arg(1); filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, MessageFileNotFound, filename)
			os.Exit(1)
		}
		defer file.Close()

		r = file
	} else {
		r = os.Stdin
	}

	query = flag.Arg(0)

	grep := NewGrep(query, opts)

	grep.Run(r)
}
