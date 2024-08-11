package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/arya2004/xyfin/api"
	db "github.com/arya2004/xyfin/db/sqlc"
	"github.com/arya2004/xyfin/gapi"
	"github.com/arya2004/xyfin/pb"
	"github.com/arya2004/xyfin/utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/xyfin?sslmode=disable"
	serverAddress = "0.0.0.0:8080"

)


func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	// This is the main function
	conn, err := sql.Open(config.DbDriver, config.DbSource)

	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
	
}

func runGrpcServer(config utils.Configuration, store db.Store) {
	grpcServer := grpc.NewServer()
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server", err)
	}
	pb.RegisterXyfinServer(grpcServer, server)
	reflection.Register(grpcServer)

	listner, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	log.Printf("start gRPC server on %s", listner.Addr().String())

	err = grpcServer.Serve(listner)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}


func runGinServer(config utils.Configuration, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.HTTPSServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}