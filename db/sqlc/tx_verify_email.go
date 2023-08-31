package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// CreateUserTx performs a money to other
type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        Users
	VerifyEmail VerifyEmails
}

func (Store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var res VerifyEmailTxResult
	err := Store.execTx(ctx, func(q *Queries) error {
		var err error
		res.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}
		user, err := q.UpdateUser(ctx, UpdateUserParams{
			Username: res.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}
		res.User = user
		return err
	})
	return res, err
}
