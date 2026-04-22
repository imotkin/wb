package main

import (
	"fmt"
	"slices"
	"strings"
)

type Set[K comparable, V any] map[K][]V

func findAnagrams(words []string) Set[string, string] {
	anagrams := make(Set[string, string])
	first := make(map[string]string)

	for _, word := range words {
		word = strings.ToLower(word)

		runes := []rune(word)
		slices.Sort(runes)

		key := string(runes)

		if _, ok := first[key]; !ok {
			first[key] = word
		}

		anagrams[key] = append(anagrams[key], word)
	}

	set := make(Set[string, string])

	for key, list := range anagrams {
		if len(list) == 1 {
			continue
		}

		slices.Sort(list)

		firstWord := first[key]
		set[firstWord] = list
	}

	return set
}

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

	anagrams := findAnagrams(words)

	for first, list := range anagrams {
		fmt.Println(first, list)
	}
}
