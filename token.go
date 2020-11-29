package manta

type TokenGenerator interface {
	Token() (string, error)
}
