package api

import (
	"net/http"
	"time"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createOrderRequest struct {
	CustomerID int64  `json:"customer_id" binding:"required,min=1"`
	UserID     int64  `json:"user_id" binding:"required,min=1"`
	TableID    int64  `json:"table_id" binding:"required,min=1"`
	TotalPrice string `json:"total_price" binding:"required,number"`
}

type orderResponse struct {
	ID         int64     `json:"id"`
	CustomerID int64     `json:"customer_id"`
	UserID     int64     `json:"user_id"`
	TableID    int64     `json:"table_id"`
	TotalPrice string    `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.CreateOrder(ctx, db.CreateOrderParams{
		CustomerID: req.CustomerID,
		UserID:     req.UserID,
		TableID:    req.TableID,
		TotalPrice: req.TotalPrice,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderResponse := orderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		UserID:     order.UserID,
		TableID:    order.TableID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
	}

	ctx.JSON(http.StatusOK, orderResponse)
}

type getOrderRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getOrder(ctx *gin.Context) {
	var req getOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderResponse := orderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		UserID:     order.UserID,
		TableID:    order.TableID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
	}

	ctx.JSON(http.StatusOK, orderResponse)
}

type listOrdersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listOrders(ctx *gin.Context) {
	var req listOrdersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orders, err := server.store.ListOrder(ctx, db.ListOrderParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var ordersResponse = make([]orderResponse, 0)
	for _, order := range orders {
		orderResp := orderResponse{
			ID:         order.ID,
			CustomerID: order.CustomerID,
			UserID:     order.UserID,
			TableID:    order.TableID,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt,
		}
		ordersResponse = append(ordersResponse, orderResp)
	}

	ctx.JSON(http.StatusOK, ordersResponse)
}

type deleteOrderRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteOrder(ctx *gin.Context) {
	var req deleteOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteOrder(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "order deleted successfully"})
}

type orderIDUriRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type orderUpdateRequest struct {
	UserID     int64  `json:"user_id" binding:"required,min=1"`
	TableID    int64  `json:"table_id" binding:"required,min=1"`
	customerID int64  `json:"customer_id" binding:"required,min=1"`
	TotalPrice string `json:"total_price" binding:"required,number"`
	Status     string `json:"status" binding:"required"`
}

func (server *Server) updateOrder(ctx *gin.Context) {
	var reqUri orderIDUriRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqJson orderUpdateRequest
	if err := ctx.ShouldBindJSON(&reqJson); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.UpdateOrder(ctx, db.UpdateOrderParams{
		ID:         reqUri.ID,
		UserID:     reqJson.UserID,
		CustomerID: reqJson.customerID,
		TableID:    reqJson.TableID,
		TotalPrice: reqJson.TotalPrice,
		Status:     reqJson.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderResponse := orderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		UserID:     order.UserID,
		TableID:    order.TableID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
	}

	ctx.JSON(http.StatusOK, orderResponse)
}

type updateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (server *Server) updateOrderStatus(ctx *gin.Context) {
	var reqUri orderIDUriRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqJson updateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&reqJson); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     reqUri.ID,
		Status: reqJson.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderResponse := orderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		UserID:     order.UserID,
		TableID:    order.TableID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
	}

	ctx.JSON(http.StatusOK, orderResponse)
}
