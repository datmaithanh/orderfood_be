package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/datmaithanh/orderfood/utils"
	"github.com/gin-gonic/gin"
)

type newAccessTokenUserRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type newAccessTokenUserResponse struct {
	AccessToken         string    `json:"access_token"`
	AccessTokenExpiesAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) reNewAccessToken(ctx *gin.Context) {
	var req newAccessTokenUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.ID != refreshPayload.ID {
		err := fmt.Errorf("mismatched session ID")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched refresh token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.ExpiresAt.Before(time.Now()) {
		err := fmt.Errorf("session has expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, refreshPayload.Role, utils.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	newAccessTokenResponse := newAccessTokenUserResponse{
		AccessToken:         accessToken,
		AccessTokenExpiesAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, newAccessTokenResponse)
}
