package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/arya2004/xyfin/api"
	db "github.com/arya2004/xyfin/db/sqlc"
	"github.com/arya2004/xyfin/gapi"
	"github.com/arya2004/xyfin/pb"
	"github.com/arya2004/xyfin/utils"
	"github.com/arya2004/xyfin/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	//run migration
	runMigrations(config.MigrationURL,config.DbSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)
	go runGatewayServer(config, store,taskDistributor)
	runGrpcServer(config, store, taskDistributor)
	
}

func runMigrations(migrationUrl string,dbSource string ) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal("cannot create migration: ", err)
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("cannot run migration: ", err)
	}

	log.Println("migration completed")
}


func runTaskProcessor(redisopt asynq.RedisClientOpt, store db.Store) {
	processor := worker.NewRedisTaskProcessor(redisopt, store)
	err := processor.Start()
	if err != nil {
		log.Fatal("cannot run task processor", err)
	}
}


func runGrpcServer(config utils.Configuration, store db.Store, taskDist worker.TaskDistributor) {
	
	server, err := gapi.NewServer(config, store,taskDist)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)

	grpcServer := grpc.NewServer(grpcLogger)
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

func runGatewayServer(config utils.Configuration, store db.Store, taskDist worker.TaskDistributor) {
	
	server, err := gapi.NewServer(config, store, taskDist)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err  = pb.RegisterXyfinHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register service", err)
	}

	gatewayMux := http.NewServeMux()
	gatewayMux.Handle("/", grpcMux)


	listner, err := net.Listen("tcp", config.HTTPSServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	log.Printf("start gRPC server on %s", listner.Addr().String())

	err = http.Serve(listner, gatewayMux)
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