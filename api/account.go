package api

import (
	"database/sql"
	"net/http"

	db "github.com/bank-demo/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// creatAccountRequest to store creat account request
// think what will happen when an account be create, initial balance should be 0.
// because when server recieve those parameters from HTTP request are JSON, should be json tag in struct
// also valid those input parameter by
type creatAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD TWD"`
}

// implement creatAccount API, gin's HandlFunc take gin.Context as input
func (server *Server) createAccount(ctx *gin.Context) {
	var request creatAccountRequest
	// take request to ShoulBindJson
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// if err := ctx.ShouldBindJSON(&request); err != nil {
	// 	// return a error with JSON format
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// 	return
	// }
	// if iput parameters are valided correct
	// create account into DB, use them into CreateAccountParams from account.sql.go
	arg := db.CreateAccountParams{
		Owner:    request.Owner,
		Currency: request.Currency,
		Balance:  0,
	}
	// pass arg to CreatAccount() from account.sql.go, that will return Account and error
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		// the error is happen at inserting arg to DB, so use http.StatusInternalServerError
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// if no error,send StatusOK and create account object to client
	ctx.JSON(http.StatusOK, account)

}

type getAccountRequest struct {
	// because the url is accounts/:id, have to tell Gin that id is a URI parameter.
	ID int64 `uri:"id" binding:"required,min=1"`
}

// implement getAccount API
func (server *Server) getAccount(ctx *gin.Context) {
	var request getAccountRequest
	// use ShouldBindUri to bind uri
	if err := ctx.ShouldBindUri(&request); err != nil {
		// return 400 code with JSON format to client
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// if no error, connect the db and get info of account by method server.store.GetAccount
	account, err := server.store.GetAccount(ctx, request.ID)
	if err != nil {
		// if the account not existed in db
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// other error between server and db operation
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// if no error, return account in JSON format to client
	ctx.JSON(http.StatusOK, account)

}

// the url is /accounts, didn't take parameter from uri
// On Postman, use Query param -- key-value as input
// Page_ID : index of page number when query
// page_Size: number of records on 1 page - min = 5, max = 10
// Page_ID, Page_Size have to use tag "form"
type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var request listAccountRequest
	//use ShouldBindQuery to tell gin to get data from query string
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// if no error after get parameter from request to connect to db ,
	// using ListAccount to query page of account records from db
	// think the SQL syntax in account.sql
	// Limit is PageSize, Offet is number of records that db should skip.
	arg := db.ListAccountsParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get list of accounts success, and return to client
	ctx.JSON(http.StatusOK, accounts)

}
