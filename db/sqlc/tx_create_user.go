package db

import (
	"context"
)

// CreateUserTx performs a money to other
type CreateUserTxParams struct {
	CreateUserParams
	//exec after user add => commit or rollback
	AfterCreate func(user Users) error
}

type CreateUserTxResult struct {
	User Users
}

func (Store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var res CreateUserTxResult
	err := Store.execTx(ctx, func(q *Queries) error {
		var err error

		res.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(res.User)
	})
	return res, err
}
