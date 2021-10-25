package filestorage

import "fmt"

var (
	ErrOpeningDataFile       = fmt.Errorf("error occurred while opening data file")
	ErrUnmarshallingDataFile = fmt.Errorf("error occurred while unmarshalling data file")
	ErrMarshallingData       = fmt.Errorf("error occurred while marshal file")
	ErrTruncateDataFile      = fmt.Errorf("error occurred while trancating data file")
	ErrWritingToDataFile     = fmt.Errorf("error occurred while writing to data file")
)
