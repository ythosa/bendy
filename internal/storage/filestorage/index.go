package filestorage

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/ythosa/bendy/internal/indexer"
	"github.com/ythosa/bendy/pkg/utils"
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

	var indexSlices map[string][]indexer.DocID

	if err := json.Unmarshal(rawData, &indexSlices); err != nil {
		return nil, ErrUnmarshallingDataFile
	}

	return indexer.MapOnSlicesToMapOnLists(indexSlices), nil
}

func (i *IndexImpl) Update(index map[string]*list.List) error {
	file, err := os.Create(i.indexStoragePath)
	if err != nil {
		if file, err = os.Create(i.indexStoragePath); err != nil {
			logrus.Error(err)

			return ErrCreatingDataFile
		}
	}

	indexOnSlices := indexer.MapOnListsToMapOnSlices(index)

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

func (i *IndexImpl) isDataFileExists() error {
	return utils.CheckIsFilePathsValid([]string{i.indexStoragePath})
}

func (i *IndexImpl) createEmptyDataFile() error {
	file, err := os.Create(i.indexStoragePath)
	if err != nil {
		logrus.Error(err)

		return ErrCreatingDataFile
	}

	emptyData := make(map[string][]indexer.DocID)

	s, err := json.Marshal(emptyData)
	if err != nil {
		logrus.Error(err)

		return ErrMarshallingData
	}

	if _, err := file.WriteString(string(s)); err != nil {
		logrus.Error(err)

		return ErrWritingToDataFile
	}

	_ = file.Close()

	return nil
}
