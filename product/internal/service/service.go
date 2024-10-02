package service

import (
	"log"

	"github.com/kanhaiyagupta9045/product_service/internal/kafka"
	"github.com/kanhaiyagupta9045/product_service/internal/models"
	"github.com/kanhaiyagupta9045/product_service/internal/repository"
)

type ProductService struct {
	repo  *repository.ProductRepositry
	kafka *kafka.Producer
}

func NewProductService(repo *repository.ProductRepositry, kafka *kafka.Producer) *ProductService {
	return &ProductService{
		repo:  repo,
		kafka: kafka,
	}
}

func (p *ProductService) CreateProduct(product *models.Product) (interface{}, error) {

	result, err := p.repo.InsertProduct(product)
	if err != nil {
		return "", err
	}

	if err := p.kafka.ProduceMessage(kafka.PRODUCT_TOPIC, "Product Created Successfully"); err != nil {
		log.Println(err.Error())
	}
	return result, nil
}
