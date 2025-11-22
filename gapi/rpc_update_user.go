package gapi

import (
	"context"
	"database/sql"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/pb"
	"github.com/datmaithanh/orderfood/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "you can only update your own account")
	}

	arg := db.UpdateUserParams{
		Username: req.Username,
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Role: sql.NullString{
			String: req.GetRole(),
			Valid:  req.Role != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	userResponse := &pb.UpdateUserResponse{
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

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("fullname", err))
	}
	if err := val.ValidateRole(req.GetRole()); err != nil {
		violations = append(violations, fieldViolation("role", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
