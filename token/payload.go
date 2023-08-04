package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Error Token
var (
	InvalidToken    = errors.New("token has is invalid")
	ErrExpiredToken = errors.New("token has been expired")
)

// Payload data of token
type PayLoad struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssueAt   time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expried_ar"`
}

func NewPayLoad(username string, duration time.Duration) (*PayLoad, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &PayLoad{
		ID:        tokenID,
		Username:  username,
		IssueAt:   time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// VALID CHECK TOKEN
func (payLoad *PayLoad) Valid() error {
	if time.Now().After(payLoad.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
