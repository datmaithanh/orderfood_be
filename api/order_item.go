package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createOrderItemRequest struct {
	OrderID  int64  `json:"order_id" binding:"required,min=1"`
	MenuID   int64  `json:"menu_id" binding:"required,min=1"`
	Quantity int32  `json:"quantity" binding:"required,gt=0"`
}

type orderItemResponse struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	MenuID    int64     `json:"menu_id"`
	Quantity  int32     `json:"quantity"`
	Price     string    `json:"price"`
	NoteItem  string    `json:"note_item"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createOrderItem(ctx *gin.Context) {
	var req createOrderItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	menu, err := server.store.GetMenu(ctx, req.MenuID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderItem, err := server.store.CreateOrderItem(ctx, db.CreateOrderItemParams{
		OrderID:  req.OrderID,
		MenuID:   req.MenuID,
		Quantity: req.Quantity,
		Price:    menu.Price,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, req.OrderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	oldOrderPrice, err := strconv.ParseFloat(order.TotalPrice, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	priceNewItem, err := strconv.ParseFloat(orderItem.Price, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	newTotalOrderPrice := fmt.Sprintf("%.2f", oldOrderPrice+priceNewItem*float64(orderItem.Quantity))

	_, err = server.store.UpdateOrderTotalPrice(ctx, db.UpdateOrderTotalPriceParams{
		ID:         req.OrderID,
		TotalPrice: newTotalOrderPrice,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderItemResponse := orderItemResponse{
		ID:        orderItem.ID,
		OrderID:   orderItem.OrderID,
		MenuID:    orderItem.MenuID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		NoteItem:  orderItem.NoteItem,
		Status:    orderItem.Status,
		CreatedAt: orderItem.CreatedAt,
	}

	ctx.JSON(http.StatusOK, orderItemResponse)
}


type getOrderItemRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getOrderItem(ctx *gin.Context) {
	var req getOrderItemRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orderItem, err := server.store.GetOrderItem(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderItemResponse := orderItemResponse{
		ID:        orderItem.ID,
		OrderID:   orderItem.OrderID,
		MenuID:    orderItem.MenuID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		NoteItem:  orderItem.NoteItem,
		Status:    orderItem.Status,
		CreatedAt: orderItem.CreatedAt,
	}

	ctx.JSON(http.StatusOK, orderItemResponse)
}



type listOrderItemsRequest struct {
	PageID  int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listOrderItems(ctx *gin.Context) {
	var req listOrderItemsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListOrderItemParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	items, err := server.store.ListOrderItem(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var itemsResponse = make([]orderItemResponse, 0)
	for _, item := range items {
		itemResp := orderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			MenuID:    item.MenuID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			NoteItem:  item.NoteItem,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
		}
		itemsResponse = append(itemsResponse, itemResp)
	}

	ctx.JSON(http.StatusOK, itemsResponse)
}


type deleteOrderItemRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteOrderItem(ctx *gin.Context) {
	var req deleteOrderItemRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteOrderItem(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "order item deleted successfully"})
}


type updateOrderItemIDUriRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type updateOrderItemRequest struct {
	OrderID  int64  `json:"order_id" binding:"required,min=1"`
	MenuID   int64  `json:"menu_id" binding:"required,min=1"`
	Quantity int32  `json:"quantity" binding:"required,gt=0"`
	NoteItem string `json:"note_item"`
	Status   string `json:"status"`
}

func (server *Server) updateOrderItem(ctx *gin.Context) {
	var reqUri updateOrderItemIDUriRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var reqJson updateOrderItemRequest
	if err := ctx.ShouldBindJSON(&reqJson); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orderItem,  err := server.store.UpdateOrderItem(ctx, db.UpdateOrderItemParams{
		ID:        reqUri.ID,
		OrderID:   reqJson.OrderID,
		MenuID:    reqJson.MenuID,
		Quantity:  reqJson.Quantity,
		NoteItem:  reqJson.NoteItem,
		Status:    reqJson.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderItemResponse := orderItemResponse{
		ID:        orderItem.ID,
		OrderID:   orderItem.OrderID,
		MenuID:    orderItem.MenuID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		NoteItem:  orderItem.NoteItem,
		Status:    orderItem.Status,
		CreatedAt: orderItem.CreatedAt,
	}

	ctx.JSON(http.StatusOK, orderItemResponse)
}