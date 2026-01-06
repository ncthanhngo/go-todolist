package tokenprovider

import (
	"errors"
	"todolist/common"
)

// provider
type Provider interface {
	Generate(data TokenPayload, expiry int) (Token, error)
	Validate(token string) (TokenPayload, error)
	SecretKey() string
}

type TokenPayload interface {
	UserId() int
	Role() string
}

type Token interface {
	GetToken() string
}

// Loi thuong co
var (
	ErrNotFound = common.NewCustomError(
		errors.New("Token not found"),
		"Token not found",
		"ErrNotFound")
	ErrEndcodingToken = common.NewCustomError(
		errors.New("error encoding token"),
		"Err endconding token",
		"ErrEndcodingToken")
	ErrInvalidToken = common.NewCustomError(
		errors.New("Err invalid Token"),
		"Err invalid token",
		"ErrInvalidToken")
)
