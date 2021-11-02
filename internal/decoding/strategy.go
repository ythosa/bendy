package decoding

type Strategy interface {
	DecodeNext() (string, bool)
}
