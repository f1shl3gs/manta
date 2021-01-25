package token

import (
	"crypto/rand"
	"encoding/base64"
)

const defaultTokenSize = 64

type Generator struct {
	size int
}

func NewGenerator(size int) *Generator {
	if size == 0 {
		size = defaultTokenSize
	}

	return &Generator{size: size}
}

func (g *Generator) Token() (string, error) {
	src := make([]byte, g.size)
	if _, err := rand.Read(src); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(src), nil
}
