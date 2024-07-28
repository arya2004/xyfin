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
	dbSource = "postgresql://root:secret@localhost:5432/xyfin?sslmode=disable"
)

// contains DBTX
var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M){
	var err error
	testDB,err = sql.Open(dbDriver,dbSource )
	if err != nil {
		log.Fatal("can't connect to Db", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}