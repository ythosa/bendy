package decoder_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/ythosa/bendy/internal/decoder"
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
		tc := tc
		txtDecoder := decoder.NewTXTDecoder(tc.file)

		for i := 0; i < len(tc.expectedResult); i++ {
			decoded, ok := txtDecoder.DecodeNext()
			if !ok && tc.expectedResult != nil {
				t.Fatalf("decoded file part haven't next")
			}

			if decoded != tc.expectedResult[i] {
				t.Errorf("value is not expected. got: \"%s\", want: \"%s\"", decoded, tc.expectedResult[i])
			}
		}
	}
}
