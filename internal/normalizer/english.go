package normalizer

import (
	"strings"
	"unicode"
)

type EnglishNormalizer struct{}

func NewEnglishNormalizer() *EnglishNormalizer {
	return &EnglishNormalizer{}
}

func (n *EnglishNormalizer) Normalize(s string) string {
	var normalized strings.Builder

	for _, c := range s {
		switch {
		case unicode.IsLetter(c):
			normalized.WriteRune(unicode.ToLower(c))
		case c == '\'':
			normalized.WriteRune(c)
		}
	}

	return normalized.String()
}
