package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

// JWTMaker is make token
type JWTMaker struct {
	secretkey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("Invalid Key size: must be at least %d character", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *PayLoad, error) {
	payload, err := NewPayLoad(username, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretkey))
	return token, payload, err
}

// Verify check if token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*PayLoad, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, InvalidToken
		}
		return []byte(maker.secretkey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &PayLoad{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, InvalidToken
	}
	payload, ok := jwtToken.Claims.(*PayLoad)
	if !ok {
		return nil, InvalidToken
	}
	return payload, nil
}
