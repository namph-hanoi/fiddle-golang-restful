package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	fmt.Println(store)
	account1, _ := createRandomAccount(t)
	account2, _ := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx-%d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			fmt.Println(store)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		require.Equal(t, account1.Balance-int64(i+1)*amount, fromAccount.Balance)
		require.NotZero(t, fromAccount.UpdatedAt)
		require.NotZero(t, fromAccount.CreatedAt)
		require.NotEqual(t, account1.Balance, fromAccount.Balance)
		require.NotEqual(t, account1.Balance-int64(i+1)*amount, account1.Balance)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		// require.Equal(t, account2.Balance+int64(i+1)*amount, toAccount.Balance)
		require.NotZero(t, toAccount.UpdatedAt)
		require.NotZero(t, toAccount.CreatedAt)
		require.NotEqual(t, account2.Balance, toAccount.Balance)
		require.NotEqual(t, account2.Balance+int64(i+1)*amount, account2.Balance)
	}
}
