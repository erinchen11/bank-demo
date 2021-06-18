package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/bank-demo?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// connect to db
	// conn, err := sql.Open(dbDriver, dbSource)
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db : ", err)
	}
	//the New function is defined in db.go that sqlc has generated for us
	// testQueries = New(conn)
	testQueries = New(testDB)
	os.Exit(m.Run())
}
