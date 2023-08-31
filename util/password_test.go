package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassW(t *testing.T) {
	password := RandomString(6)

	hashedPassword1, err := HashPassWord(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassWord(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassWord := RandomString(6)
	err = CheckPassWord(wrongPassWord, hashedPassword1)
	require.Error(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassWord(password)
	require.NoError(t, err)
	require.NotEqual(t, hashedPassword1, hashedPassword2)

}
