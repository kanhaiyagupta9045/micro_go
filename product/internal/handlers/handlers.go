package handlers

import (
	"log"
	"net/http"
	"strconv"

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

		err, ok := Authentication(c)
		if ok {
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

		res, err := p.service.CreateProduct(&product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, res)
	}
}

func Authentication(c *gin.Context) (error, bool) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide Authorization Token"})
		return nil, true
	}
	client := &http.Client{}
	url := "http://localhost:5000/validate/user"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return nil, true
	}

	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return nil, true
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Invalid token or validation failed"})
		return nil, true
	}
	return err, false
}

func (p *ProductHandler) ListAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := p.service.ListAllProducts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func (p *ProductHandler) GetProductById() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_id := c.Param("id")
		if product_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error:": "please proivde product_id"})
			return
		}
		id, err := strconv.ParseUint(product_id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
			return
		}

		product, err := p.service.GetProductById(int(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func (p *ProductHandler) UpdateInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := Authentication(c)
		if ok {
			return
		}
		product_id := c.Param("id")
		if product_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error:": "please proivde product_id"})
			return
		}

		userID, err := strconv.ParseUint(product_id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
			return
		}
		var new_stock models.Stock
		if err := c.BindJSON(&new_stock); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
			return
		}
		if validatorErr := validate.Struct(&new_stock); validatorErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error:": validatorErr.Error()})
			return
		}

		if err := p.service.UpdateInventory(int(userID), new_stock.Stock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success:": "Inventory Update Sucessfully"})
	}
}
