package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/handler"
)

func Router(handler *handler.UserHandler) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "Upgrade", "Connection"},
		AllowCredentials: true,
	}))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
	router.POST("/user/register", handler.RegisterUser())
	router.GET("/list/users", handler.ListAllUser())
	router.GET("/user/:id", handler.GetUserByID())

	router.POST("/user/login", handler.Login())

	return router
}
