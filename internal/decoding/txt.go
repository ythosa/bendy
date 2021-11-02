package decoding

import (
	"bufio"
	"io"
)

type txtStrategy struct {
	scanner *bufio.Scanner
}

func newTXTDecoder(file io.Reader) *txtStrategy {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	return &txtStrategy{scanner: scanner}
}

func (d *txtStrategy) DecodeNext() (string, bool) {
	if d.scanner.Scan() {
		return d.scanner.Text(), true
	}

	return "", false
}
