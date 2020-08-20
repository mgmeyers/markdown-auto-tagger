package keywords

import (
	"sort"
	"strings"
	"unicode"
)

//
type TextRaker struct {
	lang *Lang
}

//
func NewTextRaker(lang *Lang) *TextRaker {
	return &TextRaker{lang}
}

// Score : (Word, Score) pair
type score struct {
	word  string
	score float64
}

type byScore []score

func (s byScore) Len() int {
	return len(s)
}

func (s byScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byScore) Less(i, j int) bool {
	return s[i].score > s[j].score
}

func (r *TextRaker) splitIntoWords(text string) []string {
	words := []string{}
	splitWords := strings.Fields(text)
	for _, word := range splitWords {
		_, currentWord := r.lang.FindRootWord(strings.TrimSpace(word))
		if currentWord != "" {
			words = append(words, currentWord)
		}
	}
	return words
}

func (r *TextRaker) generateCandidatePhrases(text string) []string {
	words := r.splitIntoWords(text)
	acceptedWords := []string{}
	for _, word := range words {
		if !r.lang.IsStopWord(word) {
			acceptedWords = append(acceptedWords, word)
		} else {
			acceptedWords = append(acceptedWords, "|")
		}
	}

	phraseList := []string{}
	phrase := ""
	for _, word := range acceptedWords {
		if word == "|" {
			phraseList = append(phraseList, phrase)
			phrase = ""
		} else {
			phrase = phrase + " " + word
		}
	}
	return phraseList
}

func (r *TextRaker) splitIntoSentences(text string) []string {
	splitFunc := func(c rune) bool {
		return unicode.IsPunct(c)
	}
	return strings.FieldsFunc(text, splitFunc)
}

func (r *TextRaker) combineScores(phraseList []string, scores map[string]float64) map[string]float64 {
	candidateScores := map[string]float64{}
	for _, phrase := range phraseList {
		words := r.splitIntoWords(phrase)
		candidateScore := float64(0.0)

		for _, word := range words {
			candidateScore += scores[word]
		}
		candidateScores[phrase] = candidateScore
	}
	return candidateScores
}

func (r *TextRaker) calculateWordScores(phraseList []string) map[string]float64 {
	frequencies := map[string]int{}
	degrees := map[string]int{}
	for _, phrase := range phraseList {
		words := r.splitIntoWords(phrase)
		length := len(words)
		degree := length - 1

		for _, word := range words {
			frequencies[word]++
			degrees[word] += degree
		}
	}
	for key := range frequencies {
		degrees[key] = degrees[key] + frequencies[key]
	}

	score := map[string]float64{}

	for key := range frequencies {
		score[key] += (float64(degrees[key]) / float64(frequencies[key]))
	}

	return score
}

func (r *TextRaker) sortScores(scores map[string]float64, topN int) []score {
	rakeScores := []score{}
	for k, v := range scores {
		rakeScores = append(rakeScores, score{k, v})
	}
	sort.Sort(byScore(rakeScores))
	if topN < len(rakeScores) && topN > 0 {
		return rakeScores[0:topN]
	}
	return rakeScores
}

func (r *TextRaker) rake(text string, topN int) map[string]float64 {
	sentences := r.splitIntoSentences(text)
	phraseList := []string{}
	for _, sentence := range sentences {
		phraseList = append(phraseList, r.generateCandidatePhrases(sentence)...)
	}
	wordScores := r.calculateWordScores(phraseList)
	candidateScores := r.combineScores(phraseList, wordScores)
	sortedScores := r.sortScores(candidateScores, topN)
	scoreDict := make(map[string]float64)
	for _, score := range sortedScores {
		scoreDict[strings.TrimSpace(score.word)] = score.score
	}
	return scoreDict
}

// RakeText : Run rake directly from text
func (r *TextRaker) RakeText(text string) map[string]float64 {
	return r.rake(text, 10)
}

// RakeTextN : Run rake directly from text and return the top N results
func (r *TextRaker) RakeTextN(text string, topN int) map[string]float64 {
	return r.rake(text, topN)
}
