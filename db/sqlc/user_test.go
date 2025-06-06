package db

import (
	"context"
	"testing"
	"time"

	"github.com/namph-hanoi/fiddle-golang-restful/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) (User, CreateUserParams) {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
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
