package gapi

import (
	"context"
	"database/sql"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/pb"
	"github.com/datmaithanh/orderfood/utils"
	"github.com/datmaithanh/orderfood/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) UpdatePasswordUser(ctx context.Context, req *pb.UpdatePasswordUserRequest) (*pb.UpdatePasswordUserResponse, error) {

	violations := validateUpdatePasswordUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to hash password: %v", err)
	}
	user, err := server.store.UpdateUserWithPassword(ctx, db.UpdateUserWithPasswordParams{
		ID:           req.GetId(),
		HashPassword: hashedPassword,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "faild to update user password", err)
	}

	userResponse := &pb.UpdatePasswordUserResponse{
		User: &pb.User{
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      user.Role,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}
	return userResponse, nil
}

func validateUpdatePasswordUserRequest(req *pb.UpdatePasswordUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
