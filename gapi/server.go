package gapi

import (
	"fmt"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/pb"
	"github.com/datmaithanh/orderfood/token"
	"github.com/datmaithanh/orderfood/utils"
	"github.com/datmaithanh/orderfood/worker"
)

type Server struct {
	pb.UnimplementedOrderFoodServiceServer
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(utils.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}
	return server, nil
}
