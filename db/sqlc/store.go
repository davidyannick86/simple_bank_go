package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store provides all functions to execute SQL queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
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

// trasnferstx perfoms money transfer from account1 to account2
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// update accounts' balance in transaction in ID order to avoid deadlocks
		// fetch current account balances
		fromAccount, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		toAccount, err := q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      fromAccount.ID,
				Balance: fromAccount.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}
			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      toAccount.ID,
				Balance: toAccount.Balance + arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      toAccount.ID,
				Balance: toAccount.Balance + arg.Amount,
			})
			if err != nil {
				return err
			}
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      fromAccount.ID,
				Balance: fromAccount.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}
