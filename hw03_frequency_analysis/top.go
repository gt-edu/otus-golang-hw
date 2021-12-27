package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

type FrequencyStringSlice struct {
	sort.StringSlice
	uniqWordsCount map[string]int
}

func (x FrequencyStringSlice) Less(i, j int) bool {
	if x.uniqWordsCount[x.StringSlice[i]] == x.uniqWordsCount[x.StringSlice[j]] {
		return x.StringSlice[i] < x.StringSlice[j]
	}
	return x.uniqWordsCount[x.StringSlice[i]] > x.uniqWordsCount[x.StringSlice[j]]
}

func (x *FrequencyStringSlice) StringExists(word string) bool {
	_, ok := x.uniqWordsCount[word]
	return ok
}

func (x *FrequencyStringSlice) AddWord(word string) {
	x.uniqWordsCount[word] = 0
	x.StringSlice = append(x.StringSlice, word)
}

func (x *FrequencyStringSlice) IncreaseWordCount(word string) {
	x.uniqWordsCount[word]++
}

func NewFrequencyStringSlice() *FrequencyStringSlice {
	return &FrequencyStringSlice{
		StringSlice:    make([]string, 0),
		uniqWordsCount: make(map[string]int),
	}
}

func Top10(inp string) []string {
	if len(inp) == 0 {
		return nil
	}

	freqStrSl := NewFrequencyStringSlice()
	words := strings.Fields(inp)
	for _, w := range words {
		ok := freqStrSl.StringExists(w)
		if !ok {
			if len(w) == 0 {
				continue
			}
			freqStrSl.AddWord(w)
		}

		freqStrSl.IncreaseWordCount(w)
	}

	sort.Sort(freqStrSl)

	if freqStrSl.Len() == 0 {
		return nil
	}

	min := 10
	if freqStrSl.Len() < 10 {
		min = freqStrSl.Len()
	}

	return freqStrSl.StringSlice[:min]
}

func PrintSliceValues(checkSlice []string, countMap map[string]int) {
	for _, w := range checkSlice {
		fmt.Printf("%s: %d\n", w, countMap[w])
	}
}
