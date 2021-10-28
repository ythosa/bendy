package filestorage

import (
	"container/list"
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

func (i *IndexImpl) Get() (map[string]*list.List, error) {
	rawData, err := ioutil.ReadFile(i.indexStoragePath)
	if err != nil {
		if _, err := os.Create(i.indexStoragePath); err != nil {
			logrus.Error(err)

			return nil, ErrCreatingDataFile
		}

		return make(map[string]*list.List), nil
	}

	var indexSlices map[string][]index.DocID

	if err := json.Unmarshal(rawData, &indexSlices); err != nil {
		return nil, ErrUnmarshallingDataFile
	}

	return index.MapOnSlicesToMapOnLists(indexSlices), nil
}

func (i *IndexImpl) Set(idx map[string]*list.List) error {
	file, err := os.Create(i.indexStoragePath)
	if err != nil {
		return ErrCreatingDataFile
	}

	indexOnSlices := index.MapOnListsToMapOnSlices(idx)

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
