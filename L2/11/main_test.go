package main

import (
	"reflect"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	cases := []struct {
		words []string
		set   Set[string, string]
	}{
		{
			[]string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"},
			Set[string, string]{
				"пятак":  []string{"пятак", "пятка", "тяпка"},
				"листок": []string{"листок", "слиток", "столик"},
			},
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			set := findAnagrams(tt.words)
			if !reflect.DeepEqual(set, tt.set) {
				t.Errorf("got: %v, expected: %v\n", set, tt.set)
			}
		})
	}
}
