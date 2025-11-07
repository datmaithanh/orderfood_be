package gapi

import (
	"fmt"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/pb"
	"github.com/datmaithanh/orderfood/token"
	"github.com/datmaithanh/orderfood/utils"
)

type Server struct {
	pb.UnimplementedOrderFoodServiceServer
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(utils.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}
	return server, nil
}
