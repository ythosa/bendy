package normalizing

import (
	"strings"
	"unicode"
)

type englishNormalizer struct{}

func newEnglishNormalizer() *englishNormalizer {
	return &englishNormalizer{}
}

func (n *englishNormalizer) Normalize(s string) string {
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
