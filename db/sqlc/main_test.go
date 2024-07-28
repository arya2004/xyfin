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


func TestMain(m *testing.M){
	conn,err := sql.Open(dbDriver,dbSource )
	if err != nil {
		log.Fatal("can't connect to Db", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}