package api

import (
	"net/http"
	"time"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createCustomerRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
}

type customerResponse struct {
	ID          int64     `json:"id"`
	FullName    string    `json:"name"`
	PhoneNumber string    `json:"phone"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
}

func (server *Server) createCustomer(ctx *gin.Context) {
	var req createCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	customer, err := server.store.CreateCustomer(ctx, db.CreateCustomerParams{

		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	customerResponse := customerResponse{
		ID:          customer.ID,
		FullName:    customer.FullName,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
		CreatedAt:   customer.CreatedAt,
	}
	ctx.JSON(http.StatusOK, customerResponse)
}


type getCustomerRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCustomer(ctx *gin.Context) {
	var req getCustomerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	customer, err := server.store.GetCustomer(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	customerResponse := customerResponse{
		ID:          customer.ID,
		FullName:    customer.FullName,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
		CreatedAt:   customer.CreatedAt,
	}
	ctx.JSON(http.StatusOK, customerResponse)
}

type listCustomerRequest struct {
	PageID  int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCustomer(ctx *gin.Context) {
	var req listCustomerRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	customers, err := server.store.ListCustomer(ctx, db.ListCustomerParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	var customerResponses = make([]customerResponse,0)
	for _, customer := range customers {
		customerResponses = append(customerResponses, customerResponse{
			ID:          customer.ID,
			FullName:    customer.FullName,
			PhoneNumber: customer.PhoneNumber,
			Email:       customer.Email,
			CreatedAt:   customer.CreatedAt,
		})
	}
	
	ctx.JSON(http.StatusOK, customerResponses)
}


type deleteCustomerRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteCustomer(ctx *gin.Context) {
	var req deleteCustomerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetCustomer(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	err = server.store.DeleteCustomer(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "customer deleted successfully"})
}
