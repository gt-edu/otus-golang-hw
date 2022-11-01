package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	reader := bufio.NewReader(r)

	for {
		line, _, err := reader.ReadLine()

		if errors.Is(err, io.EOF) {
			break
		}

		var user User
		if err := user.UnmarshalJSON(line); err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
