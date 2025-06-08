package db_test

import (
	"context"
	"testing"
	"time"

	db "github.com/namph-hanoi/fiddle-golang-restful/db/sqlc"
	"github.com/namph-hanoi/fiddle-golang-restful/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) (db.User, db.CreateUserParams) {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := db.CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.CreateRandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	return user, arg
}

func TestCreateUser(t *testing.T) {
	user, arg := createRandomUser(t)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)
}
func TestGetUser(t *testing.T) {
	user, _ := createRandomUser(t)
	response, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, response.FullName, user.FullName)
	require.Equal(t, response.HashedPassword, user.HashedPassword)
	require.Equal(t, response.FullName, user.FullName)
	require.WithinDuration(t, user.CreatedAt, response.CreatedAt, 3*time.Second)
}
