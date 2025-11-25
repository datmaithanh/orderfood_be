package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserResponse struct {
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	CreateAt time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:     req.Username,
		HashPassword: hashedPassword,
		FullName:     req.FullName,
		Email:        req.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	userResponse := UserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		CreateAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, userResponse)
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

type getUserResponse struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	CreateAt time.Time `json:"created_at"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	userResponse := getUserResponse{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		CreateAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, userResponse)
}

type getListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req getListUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	users, err := server.store.ListUser(ctx, db.ListUserParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	usersResponse := make([]getUserResponse, 0)
	for _, user := range users {
		userResponse := getUserResponse{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
			CreateAt: user.CreatedAt,
		}
		usersResponse = append(usersResponse, userResponse)
	}
	ctx.JSON(http.StatusOK, usersResponse)
}

type deleteUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	err = server.store.DeleteUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

type updateUserIdRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

type updateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var idReq updateUserIdRequest
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, idReq.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	updateParams := db.UpdateUserParams{
		Username: user.Username,
		FullName: sql.NullString{
			String: req.FullName,
			Valid:  true,
		},
		Role: sql.NullString{
			String: req.Role,
			Valid:  true,
		},
		Email: sql.NullString{
			String: req.Email,
			Valid:  true,
		},
	}
	updatedUser, err := server.store.UpdateUser(ctx, updateParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	userResponse := getUserResponse{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		FullName: updatedUser.FullName,
		Email:    updatedUser.Email,
		CreateAt: updatedUser.CreatedAt,
	}
	ctx.JSON(http.StatusOK, userResponse)
}

type updateUserWithPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

func (server *Server) updateUserWithPassword(ctx *gin.Context) {
	var idReq updateUserIdRequest
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.GetUser(ctx, idReq.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var req updateUserWithPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updateUserPassword, err := server.store.UpdateUserWithPassword(ctx, db.UpdateUserWithPasswordParams{
		Username:           idReq.Username,
		HashPassword: hashedPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	userResponse := getUserResponse{
		ID:       updateUserPassword.ID,
		Username: updateUserPassword.Username,
		FullName: updateUserPassword.FullName,
		Email:    updateUserPassword.Email,
		CreateAt: updateUserPassword.CreatedAt,
	}
	ctx.JSON(http.StatusOK, userResponse)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID            uuid.UUID    `json:"session_id"`
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiesAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken         string       `json:"refresh_token"`
	RefreshTokenExpiesAt time.Time    `json:"refresh_token_expires_at"`
	User                 UserResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, user.HashPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, user.Role, utils.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, user.Role, utils.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	UserResponse := UserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		CreateAt: user.CreatedAt,
	}
	loginUserResponse := loginUserResponse{
		SessionID:            session.ID,
		AccessToken:          accessToken,
		AccessTokenExpiesAt:  accessPayload.ExpiredAt,
		RefreshToken:         refreshToken,
		RefreshTokenExpiesAt: refreshPayload.ExpiredAt,
		User:                 UserResponse,
	}
	ctx.JSON(http.StatusOK, loginUserResponse)
}
