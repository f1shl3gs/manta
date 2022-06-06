package token

const defaultTokenSize = 64

type Generator interface {
	Token() (string, error)
}
