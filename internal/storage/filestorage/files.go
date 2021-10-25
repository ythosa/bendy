package filestorage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

type FilesImpl struct {
	indexingFilesFilenamesPath string
}

func NewFilesImpl(indexingFilesFilenamesPath string) *FilesImpl {
	return &FilesImpl{indexingFilesFilenamesPath: indexingFilesFilenamesPath}
}

func (f *FilesImpl) Get() ([]string, error) {
	rawData, err := ioutil.ReadFile(f.indexingFilesFilenamesPath)
	if err != nil {
		logrus.Error(err)

		return nil, ErrOpeningDataFile
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
	if err := os.Truncate(f.indexingFilesFilenamesPath, 0); err != nil {
		logrus.Errorf("failed to truncate: %v", err)

		return ErrTruncateDataFile
	}

	s, err := json.Marshal(filenames)
	if err != nil {
		logrus.Error(err)

		return ErrMarshallingData
	}

	file, err := os.Open(f.indexingFilesFilenamesPath)
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
