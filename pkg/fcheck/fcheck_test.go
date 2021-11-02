package fcheck_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ythosa/bendy/pkg/fcheck"
)

func TestCheckIsFilePathsValid(t *testing.T) {
	t.Parallel()

	const (
		validFilePath   = "./valid_file_path"
		invalidFilePath = "./invalid_file_path"
	)

	_, _ = os.Create(validFilePath)

	testCases := []struct {
		filePaths     []string
		expectedError bool
	}{
		{
			filePaths:     []string{validFilePath},
			expectedError: false,
		},
		{
			filePaths:     []string{invalidFilePath},
			expectedError: true,
		},
		{
			filePaths:     []string{validFilePath, invalidFilePath},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		if tc.expectedError {
			assert.NotNil(t, fcheck.CheckIsFilePathsValid(tc.filePaths))
		} else {
			assert.Nil(t, fcheck.CheckIsFilePathsValid(tc.filePaths))
		}
	}

	_ = os.Remove(validFilePath)
}
