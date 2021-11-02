package normalizing

type Strategy interface {
	Normalize(string) string
}
