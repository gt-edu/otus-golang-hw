package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

func Top10(inp string) []string {
	if len(inp) == 0 {
		return nil
	}

	words := strings.Fields(inp)
	uniqWordsCount := make(map[string]int)
	var uniqWords []string
	for _, w := range words {
		_, ok := uniqWordsCount[w]
		if !ok {
			if len(w) == 0 {
				continue
			}
			uniqWordsCount[w] = 0
			uniqWords = append(uniqWords, w)
		}

		uniqWordsCount[w]++
	}

	sort.Slice(uniqWords, func(i, j int) bool {
		if uniqWordsCount[uniqWords[i]] == uniqWordsCount[uniqWords[j]] {
			return uniqWords[i] < uniqWords[j]
		}
		return uniqWordsCount[uniqWords[i]] > uniqWordsCount[uniqWords[j]]
	})

	if len(uniqWords) == 0 {
		return nil
	}

	min := 10
	if len(uniqWords) < 10 {
		min = len(uniqWords)
	}

	return uniqWords[:min]
}

func PrintSliceValues(checkSlice []string, countMap map[string]int) {
	for _, w := range checkSlice {
		fmt.Printf("%s: %d\n", w, countMap[w])
	}
}
