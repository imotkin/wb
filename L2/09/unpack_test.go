package main

import "testing"

func TestUnpack(t *testing.T) {
	cases := []struct {
		s, want  string
		hasError bool
	}{
		{s: "a4bc2d5e", want: "aaaabccddddde"},
		{s: "abcd", want: "abcd"},
		{s: "45", hasError: true},
		{s: "", hasError: true},
		{s: `qwe\4\5`, want: "qwe45"},
		{s: `qwe\45`, want: "qwe44444"},
		{s: `qwe5`, want: "qweeeee"},
		{s: `\11`, want: "1"},
		{s: `\1`, want: "1"},
		{s: `абв`, want: "абв"},
	}

	for _, c := range cases {
		s, err := unpack(c.s)
		if err != nil {
			if !c.hasError {
				t.Errorf("unpack(%q): %v", c.s, err)
			}

			continue
		}

		if s != c.want {
			t.Errorf("got %q, expected %q", s, c.want)
		}
	}
}
