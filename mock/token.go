package mock

type Token struct {
	Current string
}

func (t *Token) Token() (string, error) {
	return t.Current, nil
}
