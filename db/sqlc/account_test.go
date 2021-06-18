package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bank-demo/util"
	"github.com/stretchr/testify/require"
)

// 根據 account.sql.go
// 要測的函數有
// CreateAccount
// GetAccount
// DeleteAccount
// ListAccounts
// UpdateAccount
// Make sure that they are independent from each o

// func TestCreatAccount(t *testing.T) {
// 	// according account.sql.go file to write unit test
// 	arg := CreateAccountParams{
// 		//Owner:    "Erin", // randomly generated ?
// 		Owner: util.RandomOwner(),
// 		//Balance:  1000,
// 		Balance: util.RandomMoney(),
// 		//Currency: "USD",
// 		Currency: util.RandomCurrency(),
// 	}

// 	// then call the testQueried.CreateAccount
// 	// 去找 *Queries有哪些method需要測試
// 	account, err := testQueries.CreateAccount(context.Background(), arg)
// 	// use testify to check the test result
// 	require.NoError(t, err)
// 	require.NotEmpty(t, account)

// 	require.Equal(t, arg.Owner, account.Owner)
// 	require.Equal(t, arg.Balance, account.Balance)
// 	require.Equal(t, arg.Currency, account.Currency)

// 	require.NotZero(t, account.ID)
// 	require.NotZero(t, account.CreatedAt)

// }

func createRandomAccount(t *testing.T) Accounts {
	// according account.sql.go file to write unit test
	arg := CreateAccountParams{
		//Owner:    "Erin", // randomly generated ?
		Owner: util.RandomOwner(),
		//Balance:  1000,
		Balance: util.RandomMoney(),
		//Currency: "USD",
		Currency: util.RandomCurrency(),
	}

	// then call the testQueried.CreateAccount
	// 去找 *Queries有哪些method需要測試
	account, err := testQueries.CreateAccount(context.Background(), arg)
	// use testify to check the test result
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreatAccount(t *testing.T) {
	// just call createRandomAccount
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	//先用隨機成帳戶並存入table
	// 再用Get Account去找出該帳戶來測試

	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	// WithinDuration to check that 2 timestamps are different
	// by at most some delta duration.
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestListAccounts(t *testing.T) {
	// because it select multiple records
	// we need to create several account
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	// WithinDuration to check that 2 timestamps are different
	// by at most some delta duration.
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	// 刪除table中的account1
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	// 成功刪除 account1 理論上不會有 error存在
	require.NoError(t, err)

	//再用GetAccount去查 account1是否還存在table
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	// 因為找不到account1 所以理論上會有Error返回
	require.Error(t, err)
	// check that the error should be sql.ErrNoRows
	//
	require.EqualError(t, err, sql.ErrNoRows.Error())
	// check thate the account2 object should be empty
	require.Empty(t, account2)

}
