package token

import (
	"encoding/base64"
	"math/rand"
)

type SizedGenerator struct {
	size int
}

func NewGenerator(size int) *SizedGenerator {
	if size == 0 {
		size = defaultTokenSize
	}

	return &SizedGenerator{size: size}
}

func (g *SizedGenerator) Token() (string, error) {
	src := make([]byte, g.size)
	if _, err := rand.Read(src); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(src), nil
}
