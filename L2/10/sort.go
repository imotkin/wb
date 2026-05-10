// Задача L2.10 - Утилита sort

package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Parser struct {
	opts Options
}

func (p Parser) Parse(row string) string {
	cols := strings.Fields(row)

	if (p.opts.ColumnIndex - 1) >= len(cols) {
		return ""
	} else if p.opts.ColumnIndex > 0 {
		if p.opts.IgnoreLeading {
			return strings.TrimLeft(cols[p.opts.ColumnIndex-1], " \t")
		}

		return cols[p.opts.ColumnIndex-1]
	} else {
		return row
	}
}

type Comparator struct {
	p        Parser
	opts     Options
	months   map[string]byte
	suffixes map[string]int64
}

func NewComparator(opts Options) Comparator {
	return Comparator{
		p:    Parser{opts: opts},
		opts: opts,
		months: map[string]byte{
			"jan": 1, "feb": 2, "mar": 3,
			"apr": 4, "may": 5, "jun": 6,
			"jul": 7, "aug": 8, "sep": 9,
			"oct": 10, "nov": 11, "dec": 12,
		},
		suffixes: map[string]int64{
			"K": 1 << 10, "k": 1 << 10,
			"M": 1 << 20, "m": 1 << 20,
			"G": 1 << 30, "g": 1 << 30,
			"T": 1 << 40, "t": 1 << 40,
			"P": 1 << 50, "p": 1 << 50,
			"E": 1 << 60, "e": 1 << 60,
		},
	}
}

func (c Comparator) compareSuffix(a, b string) int {
	var suffa, suffb string

	for suffix := range c.suffixes {
		if strings.HasSuffix(a, suffix) {
			suffa = suffix
		}
		if strings.HasSuffix(b, suffix) {
			suffb = suffix
		}
	}

	valueA := strings.TrimSuffix(a, suffa)
	valueB := strings.TrimSuffix(b, suffb)

	va, erra := strconv.ParseInt(valueA, 10, 64)
	vb, errb := strconv.ParseInt(valueB, 10, 64)

	switch {
	case erra != nil && errb != nil:
		return strings.Compare(a, b)
	case erra != nil:
		return -1
	case errb != nil:
		return 1
	}

	var multA int64 = 1
	if m, ok := c.suffixes[suffa]; ok {
		multA = m
	}

	var multB int64 = 1
	if m, ok := c.suffixes[suffb]; ok {
		multB = m
	}

	return cmp.Compare(va*multA, vb*multB)
}

func (c Comparator) compareMonth(a, b string) int {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	monthA, oka := c.months[a]
	monthB, okb := c.months[b]

	switch {
	case !oka && !okb:
		return cmp.Compare(a, b)
	case !oka:
		return -1
	case !okb:
		return 1
	default:
		return cmp.Compare(monthA, monthB)
	}
}

func (c Comparator) compareNumber(a, b string) int {
	na, erra := strconv.Atoi(a)
	nb, errb := strconv.Atoi(b)

	switch {
	case erra != nil && errb != nil:
		return cmp.Compare(a, b)
	case erra != nil:
		return -1
	case errb != nil:
		return 1
	default:
		return cmp.Compare(na, nb)
	}
}

func (c Comparator) Compare(a, b string) int {
	keyA := c.p.Parse(a)
	keyB := c.p.Parse(b)

	var result int

	switch {
	case c.opts.LikeNumbers:
		result = c.compareNumber(keyA, keyB)
	case c.opts.LikeMonths:
		result = c.compareMonth(keyA, keyB)
	case c.opts.LikeSuffixes:
		result = c.compareSuffix(keyA, keyB)
	default:
		result = cmp.Compare(keyA, keyB)
	}

	return result
}

type Options struct {
	ColumnIndex int

	LikeNumbers  bool
	LikeSuffixes bool
	LikeMonths   bool

	Reverse       bool
	Unique        bool
	IgnoreLeading bool
	CheckSort     bool
}

type Sort struct {
	comp Comparator
	opts Options
}

func NewSort(opts Options) Sort {
	return Sort{
		comp: NewComparator(opts),
		opts: opts,
	}
}

func (s Sort) Sort(r io.Reader, src string) {
	var rows []string

	if s.opts.CheckSort {
		s.Check(r, src)
		return
	}

	scan := bufio.NewScanner(r)

	for scan.Scan() {
		rows = append(rows, scan.Text())
	}

	s.sortRows(rows)
	s.printRows(rows)
}

func (s Sort) sortRows(rows []string) {
	slices.SortStableFunc(rows, s.comp.Compare)
}

func (s Sort) Check(r io.Reader, src string) {
	var previous string
	var n int64 = 1

	scan := bufio.NewScanner(r)

	if scan.Scan() {
		previous = scan.Text()
	}

	for scan.Scan() {
		n++
		current := scan.Text()

		result := s.comp.Compare(previous, current)

		invalid := (!s.opts.Reverse && result > 0) ||
			(s.opts.Reverse && result < 0) ||
			(s.opts.Unique && result == 0)

		if invalid {
			fmt.Printf("sort: %s:%d: disorder: %s\n", src, n, current)
			os.Exit(1)
		}

		previous = current
	}
}

func (s Sort) printRows(rows []string) {
	switch {
	case s.opts.Reverse && s.opts.Unique:
		fmt.Println(rows[len(rows)-1])
		for i := len(rows) - 2; i >= 0; i-- {
			if rows[i] != rows[i+1] {
				fmt.Println(rows[i])
			}
		}
	case s.opts.Reverse && !s.opts.Unique:
		for i := len(rows) - 1; i >= 0; i-- {
			fmt.Println(rows[i])
		}
	case !s.opts.Reverse && s.opts.Unique:
		fmt.Println(rows[0])
		for i := 1; i < len(rows); i++ {
			if rows[i] != rows[i-1] {
				fmt.Println(rows[i])
			}
		}
	case !s.opts.Reverse && !s.opts.Unique:
		for _, row := range rows {
			fmt.Println(row)
		}
	}
}

func parseOptions() Options {
	var opts Options

	flag.IntVar(&opts.ColumnIndex, "k", 0, "Sort by column")
	flag.BoolVar(&opts.LikeNumbers, "n", false, "Sort like numbers")
	flag.BoolVar(&opts.Reverse, "r", false, "Sort in reverse order")
	flag.BoolVar(&opts.Unique, "u", false, "Print only unique rows")
	flag.BoolVar(&opts.LikeMonths, "M", false, "Sort by month names")
	flag.BoolVar(&opts.IgnoreLeading, "b", false, "Ignore leading blanks")
	flag.BoolVar(&opts.CheckSort, "c", false, "Check if rows are sorted")
	flag.BoolVar(&opts.LikeSuffixes, "h", false, "Sort like numbers with suffix")

	flag.Parse()

	return opts
}

func main() {
	opts := parseOptions()
	s := NewSort(opts)

	var r io.Reader
	var src string

	if filename := flag.Arg(0); filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("sort: no such file or directory")
			os.Exit(1)
		}
		defer file.Close()

		r = file
		src = filename
	} else {
		r = os.Stdin
		src = "-"
	}

	s.Sort(r, src)
}
