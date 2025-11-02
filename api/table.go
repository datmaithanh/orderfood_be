package api

import (
	"fmt"
	"net/http"

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
