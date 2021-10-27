package filestorage

import "fmt"

var (
	ErrCreatingDataFile      = fmt.Errorf("error occurred while creating data file")
	ErrReadingDataFile       = fmt.Errorf("error occurred while reading data file")
	ErrUnmarshallingDataFile = fmt.Errorf("error occurred while unmarshalling data file")
	ErrMarshallingData       = fmt.Errorf("error occurred while marshal file")
	ErrTruncateDataFile      = fmt.Errorf("error occurred while trancating data file")
	ErrWritingToDataFile     = fmt.Errorf("error occurred while writing to data file")
)
