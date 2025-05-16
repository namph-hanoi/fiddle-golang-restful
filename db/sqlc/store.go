package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Future will replace the nil with &sql.TxOptions{}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction error: %v, rollbackError: %v", err, rollbackErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "creat transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "creat entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "creat entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get source account")
		sourceAccount, err := q.GetAccountForUpdate(ctx, args.FromAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "update source account")
		destAccount, err := q.GetAccountForUpdate(ctx, args.ToAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "get destination account")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      args.FromAccountID,
			Balance: sourceAccount.Balance - args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "update destination account")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      args.ToAccountID,
			Balance: destAccount.Balance + args.Amount,
		})
		if err != nil {
			return err
		}

		return nil

	})
	return result, err
}
