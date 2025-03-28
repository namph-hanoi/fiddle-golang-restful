package db

import (
	"context"
	"testing"
	"time"

	"github.com/namph-hanoi/fiddle-golang-restful/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) (Account, CreateAccountParams) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	return account, arg
}

func TestCreateAccount(t *testing.T) {
	account, arg := createRandomAccount(t)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}
func TestGetAccount(t *testing.T) {
	account, arg := createRandomAccount(t)
	response, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, response.Balance, arg.Balance)
	require.Equal(t, response.Currency, arg.Currency)
	require.Equal(t, response.Owner, arg.Owner)
	require.WithinDuration(t, account.CreatedAt, response.CreatedAt, 3*time.Second)
}
