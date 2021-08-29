package main

import (
	"database/sql"
	"log"

	"github.com/bank-demo/api"
	db "github.com/bank-demo/db/sqlc"
	"github.com/bank-demo/util"
	_ "github.com/lib/pq"
)

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgresql://root:secret@localhost:5432/bank-demo?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// )
// ----- after write config.go don
func main() {
	// use the config.go to get configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config file: ", err)
	}
	// for create a Server, need to connect to DB and creat a Store.
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to Postgres DB: ", err)
	}
	// if connection to DB success, use the conn as db.NewStore()'s input
	// and use store as server's input, so the server can handle request about DB
	store := db.NewStore(conn)
	server := api.NewServer(store)

	// start the HTTP server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start HTTP server: ", err)
	}
}
