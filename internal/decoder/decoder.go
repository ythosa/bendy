package decoder

type Decoder interface {
	DecodeNext() (string, bool)
}
