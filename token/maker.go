package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *PayLoad, error)

	//Verify check if token is valid or not
	VerifyToken(token string) (*PayLoad, error)
}
