package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"todolist/common"
	"todolist/component/tokenprovider"
)

type jwtProvider struct {
	prefix string
	secret string
}

func NewTokenJWTProvider(prefix string, secret string) *jwtProvider {
	return &jwtProvider{prefix: prefix,
		secret: secret}
}

type myClaim struct {
	PayLoad common.TokenPayload `json:"payload"`
	jwt.RegisteredClaims
}
type token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

func (t *token) GetToken() string {
	return t.Token

}
func (j *jwtProvider) SecretKey() string {
	return j.secret
}

// method
func (j *jwtProvider) Generate(data tokenprovider.TokenPayload, expiry int) (tokenprovider.Token, error) {
	//Generate the JWT
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaim{
		common.TokenPayload{
			UId:   data.UserId(),
			URole: data.Role(),
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(expiry))),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d", now.UnixNano()),
		},
	})
	myToken, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}
	//return token
	return &token{
		Token:   myToken,
		Created: now,
		Expiry:  expiry,
	}, nil
}

func (j *jwtProvider) Validate(myToken string) (tokenprovider.TokenPayload, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, tokenprovider.ErrInvalidToken
	}
	//validate token
	if !res.Valid {
		return nil, tokenprovider.ErrInvalidToken
	}
	claims, ok := res.Claims.(*myClaim)
	if !ok {
		return nil, tokenprovider.ErrInvalidToken
	}
	return &claims.PayLoad, nil
}
