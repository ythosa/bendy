package decoding

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTXTDecoder_DecodeNext(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		file           io.Reader
		expectedResult []string
	}{
		{
			file:           bytes.NewBufferString(""),
			expectedResult: nil,
		},
		{
			file:           bytes.NewBufferString("(t e s t)"),
			expectedResult: []string{"(t", "e", "s", "t)"},
		},
		{
			file:           bytes.NewBufferString("hello ! i am txt decoder. i'm rly working."),
			expectedResult: []string{"hello", "!", "i", "am", "txt", "decoder.", "i'm", "rly", "working."},
		},
	}

	for _, tc := range testCases {
		txtDecoder := newTXTDecoder(tc.file)

		for i := 0; i < len(tc.expectedResult); i++ {
			decoded, ok := txtDecoder.DecodeNext()
			assert.Equal(t, ok, true)
			assert.Equal(t, tc.expectedResult[i], decoded)
		}

		// there is no words to decode
		decoded, ok := txtDecoder.DecodeNext()
		assert.Equal(t, ok, false)
		assert.Equal(t, decoded, "")
	}
}
