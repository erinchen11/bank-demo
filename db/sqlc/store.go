package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// mock DB to test HTTP API, need to create an interface
// Server struct has db.Store to connect to DB
// so, our interface should be included methods of db.Store
// and TransferTx
type Store interface{
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) 
}

// Store struct provides all functions to execute db queries and transactions
// ------ after mock DB-----
// Store change to SQLStore
type SQLStore struct {
	*Queries
	db *sql.DB
}


// NewStore creates a new Store
// func NewStore(db *sql.DB) *Store {
// 	return &Store{
// 		db:      db,
// 		Queries: New(db),
// 	}
// }
//----- after mock DB------
//NewStore() function should not return a pointer, 
//but just a Store interface. And inside, 
//it should return the real DB implementation of the interface, 
//which is SQLStore.
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}


// execTx executes a function within a database transaction
//----- after mock DB------ Store change to SQLStore
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Call store.db.BiginTx() for start a new db transaction
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	// call New() with created transaction, and get back a new Queries object

	q := New(tx)
	// Now we have the Queries that runs within transaction
	// we can call the input function with that queries, and get back an error
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	// all operation successful, commit the transaction
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json: "from_account_id"`
	ToAccountID   int64 `json: "to_account_id"`
	Amount        int64 `json: "amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfers `json: "transfer"`
	FromAccount Accounts  `json: "from_account"`
	ToAccount   Accounts  `json: "to_account"`
	FromEntry   Entries   `json: "from_entry"`
	ToEntry     Entries   `json: "to_entry"`
}
var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer recorde, add account entries,
// and update accounts' balance with a singlge database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	// empty TransferTxresult
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// we can use the Queries object to call any individual CRUD function that it provides.
		// the Queries object is created from 1 single database transaction
		// so all of its provided methods that we call will be run within that transaction

		// 把編號後的tx帶入
		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")
        
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		//add account entries
		fmt.Println(txName, "create entry1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			// because money is moving out of this account
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry2")

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			// because money is moving out of this account
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		// transfer record and 2 account entries are created
		// TODO: update accounts' balance
		// get account -> update its balance
		fmt.Println(txName, "get account 1")
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(txName, "update account 1")

		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: account1.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "get account 2")

		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(txName, "update account ")

		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.ToAccountID,
		// 	Balance: account2.Balance + arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err

}
