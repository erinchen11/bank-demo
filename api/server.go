package api

import (

	db "github.com/bank-demo/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server servers HTTP request for bank service
type Server struct {
	router *gin.Engine
	store  db.Store
}

// NewServer return a new HTTP server and setup router
// -------after mock DB------
// *db.Store change to db.Store, because interface
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	// add routes to router, first API is POST method
	// server.creatAccount is a hendler
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)

	server.router = router
	return server

}



// because router is private in api package, need to an public method for call from outside
// Start method is run HTTP Server on the input address and listen for API request
func (server *Server) Start(address string) error {
	return server.router.Run(address)

}
// errorResponse will return map[string]interface{} to client
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// then go to main.go to start sever
