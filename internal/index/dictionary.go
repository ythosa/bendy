package index

import (
	"encoding/gob"
	"fmt"
	"os"
)

func encodeDictionaryToFile(dictionary []string, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error while creating file: %w", err)
	}

	e := gob.NewEncoder(f)
	if err := e.Encode(dictionary); err != nil {
		return fmt.Errorf("error while encoding: %w", err)
	}

	_ = f.Close()

	return nil
}

func decodeDictionaryFromFile(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while opening file: %w", err)
	}

	var terms []string

	d := gob.NewDecoder(f)
	if err := d.Decode(&terms); err != nil {
		return nil, fmt.Errorf("error while decoding terms from file: %w", err)
	}

	_ = f.Close()

	return terms, nil
}
