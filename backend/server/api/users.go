package api

import (
	"github.com/Spiria-Digital/expense-manager/server/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/service"
)

type userResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ListUsers
// @Summary List users
// @Description List users
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.OutgoingUser
// @Failure 500 {object} models.ErrorResponse
// @Router /users [get]
func ListUsers(ctx *gin.Context) {
	db := ctx.MustGet("db").(*bun.DB)

	users, err := service.ListUsers(ctx, db)
	if err != nil {
		log.Err(err).Msg("failed to list users")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to list users",
		})
		return
	}

	ctx.JSON(200, users)
}
