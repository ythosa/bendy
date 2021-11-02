package filestorage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/ythosa/bendy/internal/index"
)

type IndexImpl struct {
	indexStoragePath string
}

func NewIndexImpl(indexStoragePath string) *IndexImpl {
	return &IndexImpl{indexStoragePath: indexStoragePath}
}

func (i *IndexImpl) Get() (index.InvertIndex, error) {
	rawData, err := ioutil.ReadFile(i.indexStoragePath)
	if err != nil {
		if _, err := os.Create(i.indexStoragePath); err != nil {
			logrus.Error(err)

			return nil, ErrCreatingDataFile
		}

		return make(index.InvertIndex), nil
	}

	var indexSlices map[string][]index.DocID

	if err := json.Unmarshal(rawData, &indexSlices); err != nil {
		return nil, ErrUnmarshallingDataFile
	}

	return invertIndexFromDecoded(indexSlices), nil
}

func (i *IndexImpl) Set(idx index.InvertIndex) error {
	file, err := os.Create(i.indexStoragePath)
	if err != nil {
		return ErrCreatingDataFile
	}

	indexOnSlices := invertIndexToEncoded(idx)

	s, err := json.Marshal(indexOnSlices)
	if err != nil {
		logrus.Error(err)

		return ErrMarshallingData
	}

	if _, err := file.Write(s); err != nil {
		logrus.Error(err)

		return ErrWritingToDataFile
	}

	_ = file.Close()

	return nil
}

func invertIndexToEncoded(i index.InvertIndex) map[string][]index.DocID {
	result := make(map[string][]index.DocID)
	for k, v := range i {
		result[k] = index.ListToSlice(v.DocIDs)
	}

	return result
}

func invertIndexFromDecoded(i map[string][]index.DocID) index.InvertIndex {
	result := make(index.InvertIndex)
	for k, v := range i {
		result[k] = index.NewIndex(index.SliceToList(v))
	}

	return result
}
