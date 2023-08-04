package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomUser(t *testing.T) Users {
	hashedPassword, err := util.HashPassWord(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
func TestGetUser(t *testing.T) {
	account1 := createRandomUser(t)
	account2, err := testQueries.GetUser(context.Background(), account1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, account1.HashedPassword, account2.HashedPassword)
	require.Equal(t, account1.FullName, account2.FullName)
	require.Equal(t, account1.Email, account2.Email)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}
