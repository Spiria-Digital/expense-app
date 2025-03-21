package server

import (
	"github.com/Spiria-Digital/expense-manager/server/api"
	"github.com/Spiria-Digital/expense-manager/server/docs"
	"github.com/Spiria-Digital/expense-manager/server/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

var Router = gin.Default()

func init() {
	apiGroup := Router.Group("/api")
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}
	apiGroup.Use(cors.New(corsConfig))
	docs.SwaggerInfo.BasePath = "/api"
	apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup.GET("/status", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// inject database
	apiGroup.Use(func(context *gin.Context) {
		context.Set("db", BunDB)
		context.Next()
	})

	auth := apiGroup.Group("/auth")
	{
		auth.POST("/login", api.UserLogin)
		auth.POST("/register", api.UserRegistration)
	}

	user := apiGroup.Group("/users")
	{
		user.GET("/", api.ListUsers)
		user.Use(middleware.JWTMiddleware())
	}

	expense := apiGroup.Group("/expenses")
	{
		expense.Use(middleware.JWTMiddleware())
		expense.POST("/", api.CreateExpense)
		expense.GET("/", api.ListExpenses)
		expense.GET("/:id", api.GetExpense)
		expense.PUT("/:id", api.UpdateExpense)
		expense.DELETE("/:id", api.DeleteExpense)
	}

	categories := apiGroup.Group("/categories")
	{
		categories.Use(middleware.JWTMiddleware())
		categories.POST("/", api.CreateCategory)
		categories.GET("/", api.ListCategories)
	}
}
