package keywords

//
type CharFilter struct {
	wordSeparators     [25]string
	sentenceSeparators [4]string
}

//
func NewCharFilter() *CharFilter {
	return &CharFilter{
		[25]string{
			":",
			" ",
			",",
			"'",
			"‘",
			"’",
			"\"",
			"“",
			"”",
			")",
			"(",
			"[",
			"]",
			"{",
			"}",
			"\"",
			";",
			">",
			"<",
			"%",
			"@",
			"&",
			"=",
			"#",
			"—",
		},
		[4]string{"!", ".", "?", "\n"},
	}
}

//
func (r *CharFilter) IsWordSeparator(rn rune) bool {
	chr := string(rn)

	for _, val := range r.wordSeparators {
		if chr == val {
			return true
		}
	}

	return r.IsSentenceSeparator(rn)
}

//
func (r *CharFilter) IsSentenceSeparator(rn rune) bool {
	chr := string(rn)

	for _, val := range r.sentenceSeparators {
		if chr == val {
			return true
		}
	}

	return false
}
