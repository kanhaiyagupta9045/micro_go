package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/handlers"
)

func OrderRoutes(handler *handlers.OrderHandler) *gin.Engine {
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

	router.POST("/create/order", handler.PlaceOrder())
	router.POST("/order/:id", handler.OrderbyID())
	router.POST("/ship/order/:id", handler.ShipOrder())

	return router
}
