package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/service"
)

var (
	singingKey = []byte("QUsIiwWh&8E8Qflbo^V1CoKqWn#mEndELkP")
)

const appName = "expense-manager-app"

func validateAuthHeader(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing authorization header")
	}
	if len(header) < 7 || header[:7] != "Bearer " {
		return "", errors.New("invalid authorization header")
	}

	token := strings.Split(header, " ")[1]
	if token == "" {
		return "", errors.New("malformed authorization header")
	}
	return token, nil
}

func generateJWT(ownerId int, exp time.Time) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(ownerId),
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    appName,
		Audience:  jwt.ClaimStrings{"expense-manager-api"},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(singingKey)
}

func GenerateToken(owner int) (string, error) {
	expiry := time.Now().Add(time.Hour * 24)
	token, err := generateJWT(owner, expiry)
	if err != nil {
		return "", err
	}
	return token, nil
}

func validateJWT(signedToken string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return singingKey, nil
		})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey, err := validateAuthHeader(c.GetHeader("Authorization"))
		if err != nil {
			log.Err(err).Msg("failed to validate auth header")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := validateJWT(authKey)
		if err != nil {
			log.Err(err).Msg("failed to validate jwt")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ownerId, err := strconv.Atoi(claims.Subject)
		if err != nil {
			log.Err(err).Msg("failed to parse owner id")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		db := c.MustGet("db").(*bun.DB)
		user, err := service.GetUserById(c, db, ownerId)
		if err != nil {
			log.Err(err).Msg("failed to get user from db")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
