package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/models"
	"github.com/Spiria-Digital/expense-manager/server/service"
)

type createExpenseRequest struct {
	Amount      float64 `json:"amount"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Merchant    string  `json:"merchant"`
	Date        string  `json:"date"`
	CategoryId  int     `json:"categoryId"`
}

// CreateExpense creates a new expense
// @Summary Create a new expense
// @Description Create a new expense
// @Accept  json
// @Produce  json
// @Tags expenses
// @Param Authorization header string true "Bearer token"
// @Param expense body createExpenseRequest true "Expense object"
// @Success 201 {object} models.Expense
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses [post]
func CreateExpense(ctx *gin.Context) {
	var req createExpenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Err(err).Msg("Error binding JSON")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expenseDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		log.Err(err).Msg("Error parsing date")
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"message": err.Error(), "error": "Invalid date format, expected YYYY-MM-DD"})
		return
	}

	currentUser := ctx.MustGet("user").(*models.User)
	entity := models.Expense{
		Amount:      req.Amount,
		Title:       req.Title,
		Description: req.Description,
		Merchant:    req.Merchant,
		Date:        expenseDate,
		CategoryID:  req.CategoryId,
		OwnerID:     currentUser.ID,
	}
	db := ctx.MustGet("db").(*bun.DB)
	if err := service.CreateExpense(ctx, db, &entity); err != nil {
		log.Err(err).Msg("Error creating expense")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating expense"})
		return
	}
	ctx.JSON(201, entity)
}

// ListExpenses returns a list of expenses
// @Summary Get a list of expenses
// @Description Get a list of expenses
// @Accept  json
// @Produce  json
// @Tags expenses
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Expense
// @Failure 500 {object} map[string]string
// @Router /expenses [get]
func ListExpenses(ctx *gin.Context) {
	db := ctx.MustGet("db").(*bun.DB)
	currentUser := ctx.MustGet("user").(*models.User)
	// TODO check if there is a category filter in the query string and use it to call ListExpensesByCategory
	expenses, err := service.ListExpenses(ctx, db, currentUser.ID)
	if err != nil {
		log.Err(err).Msg("Error getting expenses")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting expenses"})
		return
	}
	ctx.JSON(200, expenses)
}

// GetExpense returns a single expense
// @Summary Get a single expense
// @Description Get a single expense
// @Accept  json
// @Produce  json
// @Tags expenses
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Expense ID"
// @Success 200 {object} models.Expense
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/{id} [get]
func GetExpense(ctx *gin.Context) {
	expenseID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Err(err).Msg("Error parsing expense ID")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}
	db := ctx.MustGet("db").(*bun.DB)
	currentUser := ctx.MustGet("user").(*models.User)

	expense, err := service.GetExpense(ctx, db, expenseID, currentUser.ID)
	if err != nil {
		log.Err(err).Msg("Error getting expense")
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting expense"})
		return
	}
	ctx.JSON(200, expense)
}

// UpdateExpense updates an existing expense
// @Summary Update an existing expense
// @Description Update an existing expense
// @Accept  json
// @Produce  json
// @Tags expenses
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Expense ID"
// @Param expense body createExpenseRequest true "Expense object"
// @Success 200 {object} models.Expense
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/{id} [put]
func UpdateExpense(ctx *gin.Context) {
	var req createExpenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Err(err).Msg("Error binding JSON")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expenseDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		log.Err(err).Msg("Error parsing date")
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"message": err.Error(), "error": "Invalid date format, expected YYYY-MM-DD"})
		return
	}

	expenseID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Err(err).Msg("Error parsing expense ID")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}
	db := ctx.MustGet("db").(*bun.DB)
	currentUser := ctx.MustGet("user").(*models.User)

	expense, err := service.GetExpense(ctx, db, expenseID, currentUser.ID)
	if err != nil {
		log.Err(err).Msg("Error getting expense")
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting expense"})
		return
	}

	expense.Title = req.Title
	expense.Amount = req.Amount
	expense.Description = req.Description
	expense.Merchant = req.Merchant
	expense.Date = expenseDate

	if err := service.UpdateExpense(ctx, db, expense); err != nil {
		log.Err(err).Msg("Error updating expense")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error updating expense"})
		return
	}
	ctx.JSON(200, expense)
}

// DeleteExpense deletes an existing expense
// @Summary Delete an existing expense
// @Description Delete an existing expense
// @Accept  json
// @Produce  json
// @Tags expenses
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Expense ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/{id} [delete]
func DeleteExpense(ctx *gin.Context) {
	expenseID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Err(err).Msg("Error parsing expense ID")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}
	db := ctx.MustGet("db").(*bun.DB)
	currentUser := ctx.MustGet("user").(*models.User)

	if err := service.DeleteExpense(ctx, db, expenseID, currentUser.ID); err != nil {
		log.Err(err).Msg("Error deleting expense")
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error deleting expense"})
		return
	}
	ctx.Status(204)
}
