package gapi

import (
	"fmt"

	db "github.com/arya2004/xyfin/db/sqlc"
	"github.com/arya2004/xyfin/pb"
	"github.com/arya2004/xyfin/token"
	"github.com/arya2004/xyfin/utils"
	"github.com/arya2004/xyfin/worker"
)


type Server struct {
	pb.UnimplementedXyfinServer
	config utils.Configuration
	store db.Store
	tokenMaker token.Maker
	taskDistributor worker.TaskDistributor

}

func NewServer(config utils.Configuration, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}
	
	


	return server, nil

}