package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/models"
	"github.com/Spiria-Digital/expense-manager/server/service"
)

// ListCategories
// @Summary List all categories
// @Description List all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Category
// @Failure 500 {object} models.ErrorResponse
// @Router /categories [get]
func ListCategories(ctx *gin.Context) {
	db := ctx.MustGet("db").(*bun.DB)
	categories, err := service.GetCategories(ctx, db)
	if err != nil {
		log.Err(err).Msg("failed to get categories")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to get categories",
		})
		return
	}
	ctx.JSON(http.StatusOK, categories)
}

// CreateCategory
// @Summary Create a category
// @Description Create a category
// @Tags Categories
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param category body models.Category true "Category object"
// @Success 201 {object} models.Category
// @Router /categories [post]
func CreateCategory(ctx *gin.Context) {
	db := ctx.MustGet("db").(*bun.DB)
	var category models.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		log.Err(err).Msg("failed to bind category")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "failed to bind category",
			Message: err.Error(),
		})
		return
	}
	err := service.CreateCategory(ctx, db, &category)
	if err != nil {
		log.Err(err).Msg("failed to create category")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to create category",
		})
		return
	}

	ctx.JSON(http.StatusCreated, category)
}
