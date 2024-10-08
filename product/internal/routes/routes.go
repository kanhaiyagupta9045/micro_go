package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/product_service/internal/handlers"
)

func ProductRoutes(handler *handlers.ProductHandler) *gin.Engine {
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
	router.POST("/create/product", handler.CreateProduct())
	router.GET("/products/list_all_products", handler.ListAllProducts())
	router.GET("/product/:id", handler.GetProductById())
	router.PUT("/update/inventory/:id", handler.UpdateInventory())
	return router
}
