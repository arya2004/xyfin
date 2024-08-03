package main

import (
	"database/sql"
	"log"

	"github.com/arya2004/xyfin/api"
	db "github.com/arya2004/xyfin/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/xyfin?sslmode=disable"
	serverAddress = "0.0.0.0:8080"

)


func main() {
	// This is the main function
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
	
}