package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Cut struct {
	opts Options
}

func (c *Cut) parseFields(s string) []int {
	if s == "" {
		return nil
	}

	var cols []int

	fields := strings.SplitSeq(s, ",")

	for field := range fields {
		if strings.Contains(field, "-") {
			nums := strings.Split(field, "-")

			if len(nums) != 2 {
				continue
			}

			start, errs := strconv.Atoi(nums[0])
			end, erre := strconv.Atoi(nums[1])

			if errs == nil && erre == nil && start > 0 && start <= end {
				for i := start; i <= end; i++ {
					cols = append(cols, i)
				}
			}
		} else {
			i, err := strconv.Atoi(field)
			if err == nil && i > 0 {
				cols = append(cols, i)
			}
		}
	}

	if len(cols) == 0 {
		return nil
	}

	m := make(map[int]struct{})

	for _, c := range cols {
		m[c] = struct{}{}
	}

	unique := slices.Collect(maps.Keys(m))

	slices.Sort(unique)

	return unique
}

func (c *Cut) printRow(rowText string, indexes []int) {
	cols := strings.Split(rowText, c.opts.Delimiter)

	if len(cols) == 1 {
		if !c.opts.OnlyDelimited {
			fmt.Println(rowText)
		}

		return
	}

	selected := make([]string, 0, len(indexes))

	for _, idx := range indexes {
		if idx <= len(cols) {
			selected = append(selected, cols[idx-1])
		}
	}

	fmt.Println(strings.Join(selected, c.opts.Delimiter))
}

func (c *Cut) Run(args []string) error {
	var r io.Reader = os.Stdin

	if len(args) > 0 {
		filename := args[0]

		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		r = file
	}

	fields := c.parseFields(c.opts.Fields)

	scan := bufio.NewScanner(r)

	for scan.Scan() {
		c.printRow(scan.Text(), fields)
	}

	return scan.Err()
}

type Options struct {
	Fields        string
	Delimiter     string
	OnlyDelimited bool
}

func run() error {
	var c Cut

	flag.StringVar(&c.opts.Fields, "f", "", "Number of fields")
	flag.StringVar(&c.opts.Delimiter, "d", "\t", "Selected column delimiter")
	flag.BoolVar(&c.opts.OnlyDelimited, "s", false, "Print only lines containing delimiter")

	flag.Parse()

	return c.Run(flag.Args())
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("cut: %v\n", err)
	}
}
