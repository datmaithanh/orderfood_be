package gapi

import (
	"context"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/pb"
	"github.com/datmaithanh/orderfood/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to hash password: %v", err)
	}
	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:     req.GetUsername(),
		HashPassword: hashedPassword,
		FullName:     req.GetFullname(),
		Email:        req.GetEmail(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to create user: %v", err)

	}
	userResponse := &pb.CreateUserResponse{
		User: &pb.User{
			Username:  user.Username,
			Fullname:  user.FullName,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt) ,
		},
	}
	return userResponse, nil
}
