package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder
	runes := []rune(input)
	length := len(runes)

	for i := 0; i < length; i++ {
		char := runes[i]

		if char == '\\' {
			if i+1 >= length {
				return "", ErrInvalidString
			}
			char = runes[i+1]
			if char != '\\' && !unicode.IsDigit(char) {
				return "", ErrInvalidString
			}
			i++
		} else if unicode.IsDigit(char) {
			return "", ErrInvalidString
		}

		if i+1 < length && unicode.IsDigit(runes[i+1]) {
			count := int(runes[i+1] - '0')
			result.WriteString(strings.Repeat(string(char), count))
			i++
		} else {
			result.WriteRune(char)
		}
	}

	return result.String(), nil
}
