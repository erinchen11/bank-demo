package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	// we will send the money from account1 to account2
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(" >>> before: ", account1.Balance, account2.Balance)
	// write database transaction is something we must always be vary carefule with
	// must handle the concurrency carefully.
	// the best way to make the transaction well is to run it with several concurrent go routines.
	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// have to assign a name for each transaction and pass it into the TransferTx() function
   		// via the context argument.
		txName := fmt.Sprintf("tx %d", i + 1)
		go func() {
			// inside the go routine
			// we call store.TransferTx() function
			// this function will return a result or an err
			// we cannot just use testify require to check them right here
			// because this function is running inside a different go routine
			// from the one that our TestTransferTx function is running on.
			// So there's no guarantee that it will stop the whole test if a condition is not satisfied
			// The correct way to verify the error and result is to send them back to the main goroutine
			// that our test is running on, and check them from there.
			// To do that, we can use channels
			// Channel is designed to connect concurrent goroutines,
			// and allow them to safely share data with each other without explicit locking.

			// In this case, we need one channel to receive err,
			// and one other channel to receive the TransferTxResult.

			// we will pass in a new context with the transaction name.
			// WithValue : key-value, value is the transaction name
			// in store.go, we declare a txKey, later we will have to use this key to get transaction name
			// from the input context of the TransferTx() function.
			ctx := context.WithValue(context.Background(),txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			// Now, inside the goroutine we can send error to the errors channel
			errs <- err
			// inside the goroutine we can send result to the result channel
			results <- result

		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer after the result is not empty
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//TODO :  check accounts' balance
		fmt.Println(" >>> tx: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n) // n is the number of transactions
		// to check those we need a map existed
		// we check that the existed map should not contain k
		require.NotContains(t, existed, k)
		// after that we set existed[k] to true
		existed[k] = true
	}

	// check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(" >>> after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	// error : 
	// before:  564 48
 	//>>> tx:  554 58
 	//>>> tx:  544 68
 	//>>> tx:  544 78
	//第三筆資料沒有-10，回顧一下 account.sql中的 getAccount
	// SELECT * FROM accounts
	// WHERE id = $1 LIMIT 1;
	// it doesn't block other transactions from reading the same Account record.
	// Thesefore, 2 concurrent transactions can get the same value of the account1
	// 就是其中一個transaction拿到的還是舊資料
	// 我們要在 account.sql 裡面新增加一個 GetAccountForUpdate
	// 在store.go中 就要改用 GetAccountForUpdate
}
