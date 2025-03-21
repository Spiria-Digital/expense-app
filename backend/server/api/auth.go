package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/middleware"
	"github.com/Spiria-Digital/expense-manager/server/models"
	"github.com/Spiria-Digital/expense-manager/server/service"
	"github.com/Spiria-Digital/expense-manager/server/utils"
)

var InvalidEmailOrPassword = "invalid email or password"

type userRegistrationRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type userRegistrationResponse struct {
	Message string `json:"message"`
}

// UserRegistration
// @Summary User registration
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body userRegistrationRequest true "User registration request"
// @Success 201 {object} userRegistrationResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
func UserRegistration(ctx *gin.Context) {
	var user userRegistrationRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request body",
			Message: err.Error(),
		})
		return
	}

	// hash the password
	hashedPw, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Err(err).Msg("password hashing error")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "password hashing error",
		})
		return
	}

	user.Password = hashedPw

	// save the user
	db := ctx.MustGet("db").(*bun.DB)
	entity := &models.User{
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	if err := service.CreateUser(ctx, db, entity); err != nil {
		log.Err(err).Msg("user creation error")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "user creation error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, userRegistrationResponse{Message: "User created successfully"})
}

type userLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userLoginResponse struct {
	Token string `json:"token"`
}

// UserLogin
// @Summary User login
// @Description Login a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body userLoginRequest true "User login request"
// @Success 200 {object} userLoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func UserLogin(ctx *gin.Context) {
	var user userLoginRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request body",
			Message: err.Error(),
		})
		return
	}

	// get the user
	db := ctx.MustGet("db").(*bun.DB)
	entity, err := service.GetUserByEmail(ctx, db, user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: InvalidEmailOrPassword,
			})
			return
		}

		log.Err(err).Msg("user retrieval error")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "user retrieval error",
		})
		return
	}

	// check the password
	if !utils.CheckPasswordHash(user.Password, entity.Password) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: InvalidEmailOrPassword,
		})
		return
	}

	// generate the token
	token, err := middleware.GenerateToken(entity.ID)
	if err != nil {
		log.Err(err).Msg("token generation error")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "token generation error",
		})
		return
	}

	ctx.JSON(http.StatusOK, userLoginResponse{Token: token})
}
