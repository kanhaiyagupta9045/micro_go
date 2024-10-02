package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kanhaiyagupta9045/product_service/internal/models"
	"github.com/kanhaiyagupta9045/product_service/internal/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(srv *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: srv,
	}
}

var validate = validator.New()

func (p *ProductHandler) CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide Authorization Token"})
			return
		}
		client := &http.Client{}
		url := "http://localhost:5000/validate/user"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Invalid token or validation failed"})
			return
		}

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validationErr := validate.Struct(product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		result, err := p.service.CreateProduct(&product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"Product Created Succesfully:": result})
	}
}
