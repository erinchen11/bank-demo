package api

import (
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
	Owner    string `json: "owner" binding:"required"`
	Currency string `json: "currency" binding:"required,oneof=USD TWD"`
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
