package filestorage

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/ythosa/bendy/internal/indexer"
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
		logrus.Error(err)

		return nil, ErrOpeningDataFile
	}

	var indexSlices map[string][]indexer.DocID

	if err := json.Unmarshal(rawData, &indexSlices); err != nil {
		return nil, ErrUnmarshallingDataFile
	}

	return indexer.MapOnSlicesToMapOnLists(indexSlices), nil
}

func (i *IndexImpl) Update(index map[string]*list.List) error {
	if err := os.Truncate(i.indexStoragePath, 0); err != nil {
		logrus.Errorf("failed to truncate: %v", err)

		return ErrTruncateDataFile
	}

	indexOnSlices := indexer.MapOnListsToMapOnSlices(index)

	s, err := json.Marshal(indexOnSlices)
	if err != nil {
		logrus.Error(err)

		return ErrMarshallingData
	}

	file, err := os.Open(i.indexStoragePath)
	if err != nil {
		logrus.Error(err)

		return ErrOpeningDataFile
	}

	if _, err := file.Write(s); err != nil {
		logrus.Error(err)

		return ErrWritingToDataFile
	}

	_ = file.Close()

	return nil
}
