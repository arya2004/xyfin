package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/arya2004/xyfin/utils"
	_ "github.com/lib/pq"
)

// contains DBTX
var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M){
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	testDB,err = sql.Open(config.DbDriver,config.DbSource )
	if err != nil {
		log.Fatal("can't connect to Db", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}