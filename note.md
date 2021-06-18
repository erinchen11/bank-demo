## Generate CRUD Golang code from SQL

### step
- write config in account.sql
- run "make sqlc" to genrate go code
### Create operation

```
-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;
```
- Explaination
Labe "one" means it should return 1 single Account object.
By the init_schema.up.sql, we don't need to provide the id,
because it's an auto increment column.
Everytime a new record is inserted,
the db will automatically increase the account id sequence number,
and use it as the value of the id column.
The "created_at" column will also be automatically filled with
the default value, which is the time when the record is created.
So we only need to provide values for the owner, balance, and currency.
there 3 columns, we need to pass 3 arguments into the VALUES
RETURNING * is used to tell postgresql to return the value of all columns,
after inserting the record into accounts table, this is important,
because after the account is created,we will alwayts want to return its ID to the client.
then run "make sqlc" in terminal.
we can see three .go file in the dir sqlc
in the models.go we can see json struct,
because we set "emit_exact_table_names: true" in sqlc.yaml


### READ operation

```
-- name: GetAuthor :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;
```
- Explaination
SELECT from accounts, where id is equals to the 1st input argument
so we use the label: "one"
"LIMIT" 1, because we just want to select 1 single record

```
-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;
```
- Explaination
ListAccount will return multiple accounts records.
so we use the label: "many"
select from accounts, then order the record by id,
bacause there can be a lot of accounts in the database.
we should not select all of them at once.
Instead, we do pagination, 
so we use "LIMIT" to set the number of rows we want to get.
And use "OFFET" to tell postgres to skip this many rows
before starting to return the result.

root
secret

### Update operation
```
-- name: UpdateAccount :exec
UPDATE accounts 
SET balance = $2
WHERE id = $1;
```
- Explaination
The label "exec" won't return any data, 
bacause just updates 1 row in the database.
"UPDATE" (table name)
"SET" means only allow updating the account balance,
The account owner and currency should not be changed.
Use "WHERE" clause to specify the id of the account we want to update.
then run "make sqlc" to generate code.
if we want to return data when update table, the code in account.sql is:
```
-- name: UpdateAccount :ont
UPDATE accounts 
SET balance = $2
WHERE id = $1
RETURNING *;
```

### DELETE operation
```
-- name: DeleteAccount :exec
DELETE FROM accounts 
WHERE id = $1;
```
- Explaination

Run "make sqlc", to genrate code
we can see 
```
const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts 
WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteAccountStmt, deleteAccount, id)
	return err
}
```

### Unit Test

- 自動隨機產生test data  -> util package

### Transaction Example
Transfer 10 US from bank account1 to bank account2

1. Create a transfer record with amount = 10
2. Create an account entry for account1 with amount = -10
3. Create an account entry for account2 with amount = +10
4. Subtract 10 from the balance of account1
5. Add 10 to the balance of account2

### How to run SQL Tx
```
BEGIN;
...
COMMIT;
```

### How to debug a deadlock in DB transaction

1. make some log to see which transacition is calling which query and in which order.
2. we have to assign a name for each transaction and pass it into the TransferTx() function
   via the context argument.
3. 可以用 postgres wiki中的sql語法 找出dead lock的原因
[This link] https://wiki.postgresql.org/wiki/Lock_Monitoring

4. 將原本的 GetAccountForUpdate修改
```
-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;
```
FOR NO KEY UPDATE會告訴Postgres 我們沒有要更新account table中的Key或是ID