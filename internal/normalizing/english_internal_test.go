package normalizing

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	englishNormalizer := newEnglishNormalizer()
	for _, tc := range testCases {
		assert.Equal(t, englishNormalizer.Normalize(tc.str), tc.expected)
	}
}
