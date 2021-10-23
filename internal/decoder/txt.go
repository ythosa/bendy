package decoder

import (
	"bufio"
	"io"
)

type TXTDecoder struct {
	scanner *bufio.Scanner
}

func NewTXTDecoder(file io.Reader) *TXTDecoder {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	return &TXTDecoder{scanner: scanner}
}

func (d *TXTDecoder) DecodeNext() (string, bool) {
	if d.scanner.Scan() {
		return d.scanner.Text(), true
	}

	return "", false
}
