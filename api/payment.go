package api

import (
	"net/http"
	"time"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createPaymentRequest struct {
	OrderID       int64  `json:"order_id" binding:"required,min=1"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=Cash Card Online"`
}

type paymentResponse struct {
	ID            int64     `json:"id"`
	OrderID       int64     `json:"order_id"`
	PaymentMethod string    `json:"payment_method"`
	Amount        string    `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func (server *Server) createPayment(ctx *gin.Context) {
	var req createPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, req.OrderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payment, err := server.store.CreatePayment(ctx, db.CreatePaymentParams{
		OrderID:       req.OrderID,
		PaymentMethod: req.PaymentMethod,
		Amount:        order.TotalPrice,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	paymentResponse := paymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		PaymentMethod: payment.PaymentMethod,
		Amount:        payment.Amount,
		Status:        payment.Status,
		CreatedAt:     payment.CreatedAt,
	}

	ctx.JSON(http.StatusOK, paymentResponse)
}


type getPaymentRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getPayment(ctx *gin.Context) {
	var req getPaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment, err := server.store.GetPayment(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	paymentResponse := paymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		PaymentMethod: payment.PaymentMethod,
		Amount:        payment.Amount,
		Status:        payment.Status,
		CreatedAt:     payment.CreatedAt,
	}

	ctx.JSON(http.StatusOK, paymentResponse)
}


type listPaymentsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listPayments (ctx *gin.Context){
	var req listPaymentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payments, err := server.store.ListPayment(ctx, db.ListPaymentParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var paymentResponses = make([]paymentResponse, 0)
	for _, payment := range payments {
		paymentResponse := paymentResponse{
			ID:            payment.ID,
			OrderID:       payment.OrderID,
			PaymentMethod: payment.PaymentMethod,
			Amount:        payment.Amount,
			Status:        payment.Status,
			CreatedAt:     payment.CreatedAt,
		}
		paymentResponses = append(paymentResponses, paymentResponse)
	}

	ctx.JSON(http.StatusOK, paymentResponses)
}


type deletePaymentRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deletePayment(ctx *gin.Context) {
	var req deletePaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeletePayment(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "payment deleted successfully"})
}

type updatePaymentIDUriRequest struct {
	ID     int64  `uri:"id" binding:"required,min=1"`
}

type updatePaymentRequest struct {
	Status        string `json:"status" binding:"required,oneof=Pending Completed Failed"`
}

func (server *Server) updatePaymentStatus(ctx *gin.Context) {
	var uriReq updatePaymentIDUriRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment,  err := server.store.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
		ID:     uriReq.ID,
		Status: req.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	paymentResponse := paymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		PaymentMethod: payment.PaymentMethod,
		Amount:        payment.Amount,
		Status:        payment.Status,
		CreatedAt:     payment.CreatedAt,
	}

	ctx.JSON(http.StatusOK, paymentResponse)
}