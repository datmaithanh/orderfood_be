package gapi

import (
	"context"
	"database/sql"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/pb"
	"github.com/datmaithanh/orderfood/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	err = utils.CheckPassword(req.Password, user.HashPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "incorrect password: %v", err)
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, user.Role, utils.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access tonken: %v", err)

	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, user.Role, utils.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh tonken : %v", err)

	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	loginUserResponse := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User: &pb.User{
			Username:  user.Username,
			FullName:  user.FullName,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}
	return loginUserResponse, nil
}
