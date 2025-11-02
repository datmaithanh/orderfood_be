package api

import (
	"net/http"
	"time"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createMenuRequest struct {
	Name       string `json:"name" binding:"required"`
	Price      string `json:"price" binding:"required,number"`
	CategoryID int64  `json:"category_id" binding:"required"`
}

type menuResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Price      string    `json:"price"`
	CategoryID int64     `json:"category_id"`
	Status     bool      `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (server *Server) createMenu(ctx *gin.Context) {
	var req createMenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	menu, err := server.store.CreateMenu(ctx, db.CreateMenuParams{
		Name:       req.Name,
		Price:      req.Price,
		CategoryID: req.CategoryID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	menuResponse := menuResponse{
		ID:         menu.ID,
		Name:       menu.Name,
		Price:      menu.Price,
		CategoryID: menu.CategoryID,
		Status:     menu.Status,
		CreatedAt:  menu.CreatedAt,
	}

	ctx.JSON(http.StatusOK, menuResponse)
}

type getMenuRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getMenu(ctx *gin.Context) {
	var req getMenuRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	menu, err := server.store.GetMenu(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	menuResponse := menuResponse{
		ID:         menu.ID,
		Name:       menu.Name,
		Price:      menu.Price,
		CategoryID: menu.CategoryID,
		Status:     menu.Status,
		CreatedAt:  menu.CreatedAt,
	}

	ctx.JSON(http.StatusOK, menuResponse)
}

type listMenuRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listMenu(ctx *gin.Context) {
	var req listMenuRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	menus, err := server.store.ListMenu(ctx, db.ListMenuParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	menusResponse := make([]menuResponse, 0)
	for _, menu := range menus {
		menuResp := menuResponse{
			ID:         menu.ID,
			Name:       menu.Name,
			Price:      menu.Price,
			CategoryID: menu.CategoryID,
			Status:     menu.Status,
			CreatedAt:  menu.CreatedAt,
		}
		menusResponse = append(menusResponse, menuResp)
	}

	ctx.JSON(http.StatusOK, menusResponse)
}

type updateMenuUriIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateMenuJSONRequest struct {
	Name       string `json:"name" binding:"required"`
	Price      string `json:"price" binding:"required,number"`
	CategoryID int64  `json:"category_id" binding:"required"`
}

func (server *Server) updateMenu(ctx *gin.Context) {
	var reqUriID updateMenuUriIDRequest
	if err := ctx.ShouldBindUri(&reqUriID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqJson updateMenuJSONRequest
	if err := ctx.ShouldBindJSON(&reqJson); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetMenu(ctx, reqUriID.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	menu, err := server.store.UpdateMenu(ctx, db.UpdateMenuParams{
		ID:         reqUriID.ID,
		Name:       reqJson.Name,
		Price:      reqJson.Price,
		CategoryID: reqJson.CategoryID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	menuResponse := menuResponse{
		ID:         menu.ID,
		Name:       menu.Name,
		Price:      menu.Price,
		CategoryID: menu.CategoryID,
		Status:     menu.Status,
		CreatedAt:  menu.CreatedAt,
	}

	ctx.JSON(http.StatusOK, menuResponse)
}

type deleteMenuRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteMenu(ctx *gin.Context) {
	var req deleteMenuRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetMenu(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = server.store.DeleteMenu(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "menu deleted successfully"})
}
