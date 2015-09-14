package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	file, err := ioutil.ReadFile("./text.txt")
	check(err)
	r := regexp.MustCompile("- ")
	text := strings.TrimSpace(string(file))
	text = r.ReplaceAllString(text, "")
	fmt.Println(Summarize(text))
}

func Summarize(text string) string {
	// Create list of common English words
	// to exclude from frequency map
	commonWords := commonWords()
	// make map of frequently occuring words
	frequentWords := frequency(strings.ToLower(text), commonWords)
	sentences := regexp.MustCompile("\n|\\. ").Split(text, -1)
	sentenceScores := scoreSentences(sentences, frequentWords, commonWords)
	byRank := byRankSentences(sentenceScores)
	byAppearance := byAppearanceSentences(byRank, sentences)
	return byAppearance
}

func byAppearanceSentences(byRank []string, sentences []string) string {
	byAppearance := []string{}
	for _, sent := range sentences {
		if contains(byRank, sent) {
			byAppearance = append(byAppearance, sent)
		}
	}
	return strings.Join(byAppearance, ". ")
}

func byRankSentences(sentenceScores map[string]float64) []string {
	topRanked := []string{}
	scores := []float64{}
	for _, score := range sentenceScores {
		scores = append(scores, score)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(scores)))

	for _, score := range scores {
		for sent, s := range sentenceScores {
			if score == s {
				topRanked = append(topRanked, sent)
			}
		}
	}
	topTen := int(float64(len(topRanked)) * 0.05)
	topRanked = topRanked[:topTen]
	return topRanked
}

func commonWords() []string {
	var commonWords []string
	file, err := os.Open("./commonWords.txt")
	check(err)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commonWords = append(commonWords, scanner.Text())
	}
	return commonWords
}

func frequency(text string, common []string) []string {
	frequentWords := make(map[string]int)
	r := regexp.MustCompile("\\.|\\,|;|\n|- ")
	words := strings.Split(r.ReplaceAllString(text, ""), " ")
	for _, word := range words {
		if !contains(common, word) {
			if _, ok := frequentWords[word]; ok {
				frequentWords[word]++
			} else {
				frequentWords[word] = 1
			}
		}
	}

	// Get top 10% of frequent words
	var occs []int
	var topFrequent []string

	for _, num := range frequentWords {
		occs = append(occs, num)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(occs)))
	topTen := int64(math.Floor(float64(len(occs))*0.1 + 0.5))
	occs = occs[:topTen]

	for _, val := range occs {
		for word, occ := range frequentWords {
			if occ == val {
				if !contains(topFrequent, word) {
					topFrequent = append(topFrequent, word)
				}
			}
		}
	}
	return topFrequent
}

func scoreSentences(sentences []string, freqWords []string, common []string) map[string]float64 {
	sentenceScores := make(map[string]float64)
	for _, sentence := range sentences {
		var freqCount, beg, end int
		var begSet bool
		r := regexp.MustCompile("\\,|;|- ")
		s := r.ReplaceAllString(strings.ToLower(sentence), "")
		words := strings.Split(s, " ")
		for i, word := range words {
			if contains(freqWords, word) {
				if !begSet {
					beg = i
					begSet = true
				}
				freqCount++
				end = i
			}
		}
		sentenceSlice := words[beg : end+1]
		final := []string{}
		for _, word := range sentenceSlice {
			if !contains(common, word) {
				final = append(final, word)
			}
		}
		sentenceScores[sentence] = math.Pow(float64(freqCount), 2) / float64(len(final))
	}
	return sentenceScores
}

func contains(list []string, str string) bool {
	for _, elem := range list {
		if elem == str {
			return true
		}
	}
	return false
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
