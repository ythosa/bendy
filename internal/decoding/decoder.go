package decoding

import (
	"fmt"
	"os"
	"path/filepath"
)

type Decoder interface {
	GetDecoder(file *os.File) (Strategy, error)
}

type DecoderImpl struct{}

func NewDecoderImpl() *DecoderImpl {
	return &DecoderImpl{}
}

func (d DecoderImpl) GetDecoder(file *os.File) (Strategy, error) {
	ext := filepath.Ext(file.Name())

	switch ext {
	case ".txt":
		return newTXTDecoder(file), nil
	default:
		return nil, fmt.Errorf("undefined file extension")
	}
}
