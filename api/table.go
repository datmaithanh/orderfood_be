package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/utils"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

func (server *Server) createTable(ctx *gin.Context) {
	maxTableID, err := server.store.GetMaxTableID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tableName := fmt.Sprintf("Table-%d", maxTableID.(int64)+1)
	qrText := fmt.Sprintf("%s/table/%d", utils.UrlToWebsiteOrderFood, maxTableID.(int64)+1)

	filePath := fmt.Sprintf("./qrcodes/%s.png", tableName)

	err = qrcode.WriteFile(qrText, qrcode.Medium, 256, filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cld, err := cloudinary.NewFromURL(utils.CLOUDINARY_URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot connect to Cloudinary"})
		return
	}

	uploadResult, err := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{
		Folder:   "orderfood_qrcode",
		PublicID: tableName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot upload QR to Cloudinary"})
		return
	}

	qrImageURL := uploadResult.SecureURL

	table, err := server.store.CreateTable(ctx, db.CreateTableParams{
		Name:       tableName,
		QrText:     qrText,
		QrImageUrl: qrImageURL,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, table)
}

type getTableRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type tableResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	QrText     string    `json:"qr_text"`
	QrImageUrl string    `json:"qr_image_url"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (server *Server) getTable(ctx *gin.Context) {
	var req getTableRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	table, err := server.store.GetTable(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tableResponse := tableResponse{
		ID:         table.ID,
		Name:       table.Name,
		QrText:     table.QrText,
		QrImageUrl: table.QrImageUrl,
		Status:     table.Status,
		CreatedAt:  table.CreatedAt,
	}

	ctx.JSON(http.StatusOK, tableResponse)
}

type listTablesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listTables(ctx *gin.Context) {
	var req listTablesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListTableParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	tables, err := server.store.ListTable(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var tablesResponse = make([]tableResponse, 0)
	for _, table := range tables {
		tableResp := tableResponse{
			ID:         table.ID,
			Name:       table.Name,
			QrText:     table.QrText,
			QrImageUrl: table.QrImageUrl,
			Status:     table.Status,
			CreatedAt:  table.CreatedAt,
		}
		tablesResponse = append(tablesResponse, tableResp)
	}

	ctx.JSON(http.StatusOK, tablesResponse)
}

type updateUriIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateTableStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=available occupied reserved"`
}

func (server *Server) updateTableStatus(ctx *gin.Context) {
	var reqUriID updateUriIDRequest
	if err := ctx.ShouldBindUri(&reqUriID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqJson updateTableStatusRequest
	if err := ctx.ShouldBindJSON(&reqJson); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	table, err := server.store.UpdateTable(ctx, db.UpdateTableParams{
		ID:     reqUriID.ID,
		Status: reqJson.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	tableResponse := tableResponse{
		ID:         table.ID,
		Name:       table.Name,
		QrText:     table.QrText,
		QrImageUrl: table.QrImageUrl,
		Status:     table.Status,
		CreatedAt:  table.CreatedAt,
	}

	ctx.JSON(http.StatusOK, tableResponse)
}


type deleteTableRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteTable(ctx *gin.Context) {
	var req deleteTableRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetTable(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = server.store.DeleteTable(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "table deleted successfully"})
}
