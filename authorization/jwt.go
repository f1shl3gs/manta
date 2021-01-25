package authorization

import (
	"context"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/f1shl3gs/manta"
)

type Token struct {
	jwt.StandardClaims

	// KeyID is the identifier of the key used to sign the token
	KeyID string `json:"kid"`

	// UserID is the identifier of the Owner
	UserID manta.ID `json:"uid"`

	Permissions []manta.Permission `json:"-"`
}

func (token *Token) PermissionSet() manta.PermissionSet {
	return token.Permissions
}

func (token *Token) Identifier() manta.ID {
	id, err := manta.IDFromString(token.Id)
	if err != nil || id == nil {
		return manta.ID(1)
	}

	return *id
}

func (token *Token) GetUserID() manta.ID {
	return token.UserID
}

func (token *Token) Kind() string {
	return "jwt"
}

type TokenParser struct {
	keyring manta.Keyring
	parser  *jwt.Parser
}

func (t *TokenParser) Parse(v string) (*Token, error) {
	jt, err := t.parser.ParseWithClaims(v, &Token{}, func(jt *jwt.Token) (interface{}, error) {
		token, ok := jt.Claims.(*Token)
		if !ok {
			return nil, errors.New("missing kid in token claims")
		}

		kid, err := manta.IDFromString(token.KeyID)
		if err != nil {
			return nil, err
		}

		return t.keyring.Key(context.Background(), *kid)
	})

	if err != nil {
		return nil, err
	}

	token, ok := jt.Claims.(*Token)
	if !ok {
		return nil, errors.New("token is unexpected type")
	}

	return token, nil
}
