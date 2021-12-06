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
	var lastSymbol byte
	for pos := 0; pos < len(inp); pos++ {
		c := inp[pos]
		// fmt.Printf("character %c (%d) starts at byte position %d\n", c, c, pos)

		switch {
		case unicode.IsLetter(rune(c)):
			if lastSymbol != 0 {
				outBuilder.WriteByte(lastSymbol)
			}
			lastSymbol = c
		case unicode.IsDigit(rune(c)):
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
			if lastSymbol != 0 {
				outBuilder.WriteByte(lastSymbol)
			}
			pos++
			if pos >= len(inp) {
				return "", ErrInvalidString
			}
			lastSymbol = inp[pos]
			if string(lastSymbol) != `\` && !unicode.IsDigit(rune(lastSymbol)) {
				return "", ErrInvalidString
			}
		default:
			return "", ErrInvalidString
		}
	}

	if lastSymbol != 0 {
		outBuilder.WriteByte(lastSymbol)
	}

	return outBuilder.String(), nil
}
