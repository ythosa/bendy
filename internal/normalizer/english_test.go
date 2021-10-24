package normalizer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ythosa/bendy/internal/normalizer"
)

func TestEnglishNormalizer_Normalize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		str      string
		expected string
	}{
		{
			str:      "",
			expected: "",
		},
		{
			str:      "Hi!",
			expected: "hi",
		},
		{
			str:      "I'm",
			expected: "i'm",
		},
		{
			str:      "./heLlO!",
			expected: "hello",
		},
	}

	englishNormalizer := normalizer.NewEnglishNormalizer()
	for _, tc := range testCases {
		assert.Equal(t, englishNormalizer.Normalize(tc.str), tc.expected)
	}
}
