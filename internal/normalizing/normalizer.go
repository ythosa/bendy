package normalizing

import "sync"

type Normalizer interface {
	GetNormalizer(string) (Strategy, error)
}

type NormalizerImpl struct{}

func NewNormalizerImpl() *NormalizerImpl {
	return &NormalizerImpl{}
}

func (n *NormalizerImpl) GetNormalizer(_ string) (Strategy, error) {
	var (
		once sync.Once
		eng  Strategy
	)

	once.Do(func() {
		eng = newEnglishNormalizer()
	})

	return eng, nil
}
