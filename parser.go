package keywords

import (
	"sort"

	textrank "github.com/DavidBelicza/TextRank"
	"github.com/DavidBelicza/TextRank/convert"
	"github.com/DavidBelicza/TextRank/parse"
	"github.com/DavidBelicza/TextRank/rank"
)

//
type Keywords struct {
	Words   []string
	Phrases []string
}

//
type KeywordParser struct {
	charFilter parse.Rule
	algorithm  rank.Algorithm
	lang       convert.Language
}

//
func NewKeywordParser() (*KeywordParser, error) {
	charFilter := NewCharFilter()
	algo := textrank.NewDefaultAlgorithm()
	en, err := NewEnglish()

	if err != nil {
		return nil, err
	}

	return &KeywordParser{
		charFilter: charFilter,
		algorithm:  algo,
		lang:       en,
	}, nil
}

//
func (p *KeywordParser) GetKeywords(text string, minScore float32) ([]string, []string) {
	rank := textrank.NewTextRank()
	rank.Populate(text, p.lang, p.charFilter)
	rank.Ranking(p.algorithm)

	keywords := textrank.FindSingleWords(rank)
	filtered := []string{}

	for _, keyword := range keywords {
		if keyword.Weight < minScore {
			break
		}

		filtered = append(filtered, keyword.Word)
	}

	keyPhrases := textrank.FindPhrases(rank)

	filteredPhrases := []string{}

	for _, keyPhrase := range keyPhrases {
		if keyPhrase.Weight < minScore {
			break
		}

		filteredPhrases = append(filteredPhrases, keyPhrase.Left+" "+keyPhrase.Right)
	}

	return filtered, filteredPhrases
}

//
type rakeResult struct {
	Phrase string
	Score  float64
}

//
func (p *KeywordParser) GetRakePhrases(text string, minScore float64) []string {
	rake := NewTextRaker(p.lang.(*Lang))
	result := rake.RakeText(text)

	sortRef := []rakeResult{}
	keyPhrases := []string{}

	for phrase, score := range result {
		if score < minScore {
			continue
		}

		sortRef = append(sortRef, rakeResult{
			Phrase: phrase,
			Score:  score,
		})

		keyPhrases = append(keyPhrases, phrase)
	}

	sort.SliceStable(keyPhrases, func(i, j int) bool {
		return sortRef[i].Score > sortRef[j].Score
	})

	return keyPhrases
}
