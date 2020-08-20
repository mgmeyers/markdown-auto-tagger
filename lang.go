package keywords

import (
	"strings"
	"unicode/utf8"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
)

//
type Lang struct {
	lemmatizer *golem.Lemmatizer
	stopwords  map[string]string
}

//
func NewEnglish() (*Lang, error) {
	lemma, err := golem.New(en.New())

	if err != nil {
		return nil, err
	}

	return &Lang{
		lemmatizer: lemma,
		stopwords:  english,
	}, nil
}

//
func (lang *Lang) IsStopWord(word string) bool {
	if utf8.RuneCountInString(word) <= 2 {
		return true
	}

	_, exists := lang.stopwords[strings.ToLower(word)]

	return exists
}

//
func (lang *Lang) FindRootWord(word string) (bool, string) {
	return true, lang.lemmatizer.Lemma(word)
}

//
func (lang *Lang) SetActiveLanguage(code string) {}

//
func (lang *Lang) SetWords(code string, words []string) {}
