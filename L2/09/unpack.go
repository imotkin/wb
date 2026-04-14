// Задача L2.9 - Распаковка строки

package main

import (
	"errors"
	"strings"
	"unicode"
)

const slash = rune(92)

var ErrInvalidString = errors.New("invalid input string")

func unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var (
		sb strings.Builder

		valid, slashed bool

		runes = []rune(s)
	)

	for i, r := range s {
		switch {
		case unicode.IsDigit(r) && i > 0:
			digit := int(r - 49)
			previous := runes[i-1]

			if unicode.IsLetter(previous) || unicode.IsDigit(previous) && slashed {
				sb.WriteString(strings.Repeat(string(previous), digit))
			} else if previous == slash {
				sb.WriteRune(r)
			}
		case r == slash:
			valid, slashed = true, true
		case unicode.IsLetter(r):
			valid = true
			sb.WriteRune(r)
		default:
			sb.WriteRune(r)
		}
	}

	unpacked := sb.String()

	if len(unpacked) == 0 {
		return "", ErrInvalidString
	}

	if !valid {
		return "", ErrInvalidString
	}

	return unpacked, nil
}
