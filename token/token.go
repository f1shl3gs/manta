package token

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/f1shl3gs/manta"
)

const defaultTokenSize = 64

type tokenGenerator struct {
	size int
}

func NewTokenGenerator(size int) manta.TokenGenerator {
	if size == 0 {
		size = defaultTokenSize
	}

	return &tokenGenerator{size: size}
}

func (t *tokenGenerator) Token() (string, error) {
	b := make([]byte, t.size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
