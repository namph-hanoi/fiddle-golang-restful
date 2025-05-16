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

func TestUpdateAccount(t *testing.T) {
	accountBefore, arg := createRandomAccount(t)
	updateParams := UpdateAccountParams{
		ID:      accountBefore.ID,
		Balance: arg.Balance + 1,
	}
	accountAfter, err := testQueries.UpdateAccount(context.Background(), updateParams)
	require.NotEqual(t, accountAfter.Balance, accountBefore.Balance)
	require.NoError(t, err)
	require.WithinDuration(t, accountBefore.CreatedAt, accountAfter.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account, _ := createRandomAccount(t)

	testQueries.DeleteAccount(context.Background(), account.ID)

	account, _err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, _err)
	require.EqualError(t, _err, "sql: no rows in result set")
	require.Empty(t, account)

}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountParams{
		Offset: 5,
		Limit:  5,
	}
	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}

func TestCreatingTransfer(t *testing.T) {
	account1, _ := createRandomAccount(t)
	account2, _ := createRandomAccount(t)
	amount := int64(10)
	transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	})
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.FromAccountID, account1.ID)
	require.Equal(t, transfer.ToAccountID, account2.ID)
	require.Equal(t, transfer.Amount, amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	_, err = testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
}
