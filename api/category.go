package api

import (
	"net/http"
	"time"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryResponse struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	CreateAt time.Time `json:"created_at"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.CreateCategory(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	categoryResponse := CategoryResponse{
		ID:       category.ID,
		Name:     category.Name,
		CreateAt: category.CreatedAt,
	}

	ctx.JSON(http.StatusOK, categoryResponse)
}

type GetCategoryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCategory(ctx *gin.Context) {
	var req GetCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetCategory(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	categoryResponse := CategoryResponse{
		ID:       category.ID,
		Name:     category.Name,
		CreateAt: category.CreatedAt,
	}

	ctx.JSON(http.StatusOK, categoryResponse)
}

type ListCategoryRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCategory(ctx *gin.Context) {
	var req ListCategoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	categories, err := server.store.ListCategory(ctx, db.ListCategoryParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var categoryResponses = make([]CategoryResponse, 0)
	for _, category := range categories {
		categoryResponses = append(categoryResponses, CategoryResponse{
			ID:       category.ID,
			Name:     category.Name,
			CreateAt: category.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, categoryResponses)
}

type DeleteCategoryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteCategory(ctx *gin.Context) {
	var req DeleteCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetCategory(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = server.store.DeleteCategory(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}


type UpdateCategoryURI struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

type UpdateCategoryJSON struct {
    Name string `json:"name" binding:"required"`
}

func (server *Server) updateCategory(ctx *gin.Context) {
	var reqUriID UpdateCategoryURI
	if err := ctx.ShouldBindUri(&reqUriID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var reqJson UpdateCategoryJSON
	if err := ctx.ShouldBindJSON(&reqJson); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.GetCategory(ctx, reqUriID.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	category, err := server.store.UpdateCategory(ctx, db.UpdateCategoryParams{
		ID:   reqUriID.ID,
		Name: reqJson.Name,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	categoryResponse := CategoryResponse{
		ID:       category.ID,
		Name:     category.Name,
		CreateAt: category.CreatedAt,
	}

	ctx.JSON(http.StatusOK, categoryResponse)
}
