package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inp string) (string, error) {
	outBuilder := strings.Builder{}
	var lastSymbol rune
	runes := []rune(inp)

	for pos := 0; pos < len(runes); pos++ {
		c := runes[pos]

		switch {
		case unicode.IsLetter(c):
			handleLastSymbol(lastSymbol, &outBuilder)
			lastSymbol = c
		case unicode.IsDigit(c):
			if lastSymbol == 0 {
				return "", ErrInvalidString
			}
			repCount, err := strconv.Atoi(string(c))
			if err != nil {
				return "", err
			}

			if repCount > 0 {
				outBuilder.WriteString(strings.Repeat(string(lastSymbol), repCount))
			}

			lastSymbol = 0
		case string(c) == `\`:
			handleLastSymbol(lastSymbol, &outBuilder)

			pos++
			if pos >= len(runes) {
				return "", ErrInvalidString
			}
			lastSymbol = runes[pos]
			if string(lastSymbol) != `\` && !unicode.IsDigit(lastSymbol) {
				return "", ErrInvalidString
			}
		default:
			return "", ErrInvalidString
		}
	}

	handleLastSymbol(lastSymbol, &outBuilder)

	return outBuilder.String(), nil
}

func handleLastSymbol(lastSymbol rune, outBuilder *strings.Builder) {
	if lastSymbol != 0 {
		outBuilder.WriteRune(lastSymbol)
	}
}
