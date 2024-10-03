package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/helpers"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/models"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/service"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(srv *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: *srv,
	}
}

var validate = validator.New()

func (o *OrderHandler) PlaceOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide Authorization token"})
			return
		}

		claims, err := helpers.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userID, err := strconv.ParseUint(claims.Id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		if _, err := o.service.GetUserById(int(userID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var req models.OrderItem
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(req); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		product, err := o.service.ProductById(int(req.ProductID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product does not exist"})
			return
		}

		if req.Price < product.Price {
			msg := fmt.Sprintf("The actual price of the product is : %.2f", product.Price)
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if err := o.service.PlaceOrder(int(userID), int(req.ProductID), req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Order placed successfully"})
	}
}

func (o *OrderHandler) OrderbyID() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		order_id, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		order, err := o.service.GetOrderByID(int(order_id))
		if err != nil {
			c.JSON(http.StatusBadRequest, "Order Does not exist")
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func (o *OrderHandler) ShipOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		order_id, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		if err := o.service.ShipOrder(int(order_id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": "Congratulations your Order Shipped"})

	}
}
