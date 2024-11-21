package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile(`\s+`)

type wordFreq struct {
	word  string
	count int
}

func Top10(text string) []string {
	words := re.Split(text, -1)
	frequency := make(map[string]int)

	for _, word := range words {
		word = strings.Trim(word, " \t\n.,;:!?()[]\"'`")
		word = strings.ToLower(word)
		if word != "" && word != "-" {
			frequency[word]++
		}
	}

	wordFreqs := make([]wordFreq, 0, len(frequency))
	for word, count := range frequency {
		wordFreqs = append(wordFreqs, wordFreq{word, count})
	}

	sort.Slice(wordFreqs, func(i, j int) bool {
		if wordFreqs[i].count == wordFreqs[j].count {
			return wordFreqs[i].word < wordFreqs[j].word
		}
		return wordFreqs[i].count > wordFreqs[j].count
	})

	result := make([]string, 0, 10)
	for i := 0; i < len(wordFreqs) && i < 10; i++ {
		result = append(result, wordFreqs[i].word)
	}

	return result
}
