package filestorage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/ythosa/bendy/pkg/utils"
)

type FilesImpl struct {
	indexingFilesFilenamesPath string
}

func NewFilesImpl(indexingFilesFilenamesPath string) *FilesImpl {
	return &FilesImpl{indexingFilesFilenamesPath: indexingFilesFilenamesPath}
}

func (f *FilesImpl) Get() ([]string, error) {
	if err := f.isDataFileExists(); err != nil {
		if err := f.createEmptyDataFile(); err != nil {
			return nil, err
		}
	}

	rawData, err := ioutil.ReadFile(f.indexingFilesFilenamesPath)
	if err != nil {
		return nil, ErrReadingDataFile
	}

	var files []string

	if err := json.Unmarshal(rawData, &files); err != nil {
		return nil, ErrUnmarshallingDataFile
	}

	return files, nil
}

func (f *FilesImpl) Put(filename string) error {
	files, err := f.Get()
	if err != nil {
		return err
	}

	files = append(files, filename)

	return f.update(files)
}

func (f *FilesImpl) Delete(filename string) error {
	files, err := f.Get()
	if err != nil {
		return err
	}

	updatedFiles := remove(files, filename)

	return f.update(updatedFiles)
}

func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}

	return l
}

func (f *FilesImpl) update(filenames []string) error {
	file, err := os.Create(f.indexingFilesFilenamesPath)
	if err != nil {
		logrus.Error(err)

		return ErrCreatingDataFile
	}

	s, err := json.Marshal(filenames)
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

func (f *FilesImpl) isDataFileExists() error {
	return utils.CheckIsFilePathsValid([]string{f.indexingFilesFilenamesPath})
}

func (f *FilesImpl) createEmptyDataFile() error {
	file, err := os.Create(f.indexingFilesFilenamesPath)
	if err != nil {
		logrus.Error(err)

		return ErrCreatingDataFile
	}

	emptyData := make([]string, 0)

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
