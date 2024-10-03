package service

import (
	"log"

	"github.com/kanhaiyagupta9045/product_service/internal/models"
	"github.com/kanhaiyagupta9045/product_service/internal/producer"
	"github.com/kanhaiyagupta9045/product_service/internal/repository"
)

type ProductService struct {
	repo  *repository.ProductRepositry
	kafka *producer.Producer
}

func NewProductService(repo *repository.ProductRepositry, kafka *producer.Producer) *ProductService {
	return &ProductService{
		repo:  repo,
		kafka: kafka,
	}
}

func (p *ProductService) CreateProduct(product *models.Product) (*models.Product, error) {
	product, err := p.repo.InsertProduct(product)
	if err != nil {
		return nil, err
	}
	event := models.ProductCreatedEvent{
		EventType: "Product Created",
		Data:      *product,
	}

	log.Println(event)
	go func() {
		if err := p.kafka.ProduceMessage(producer.PRODUCT_TOPIC, event); err != nil {
			log.Println("Error producing Kafka message:", err)
		}

	}()

	return product, nil
}

func (p *ProductService) ListAllProducts() ([]models.Product, error) {
	products, err := p.repo.ListAllProducts()

	if err != nil {
		return nil, err
	}

	return products, err
}

func (p *ProductService) GetProductById(product_id int) (*models.Product, error) {
	product, err := p.repo.GetProductById(product_id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductService) UpdateInventory(product_id int, newStock int) error {

	err := p.repo.UpdateInventory(uint(product_id), newStock)
	if err != nil {
		return err
	}

	event := models.InventoryUpdate{
		EventType:  "Inventory Updated",
		ProductID:  (product_id),
		StockLevel: newStock,
	}

	log.Println("Producing Kafka event for inventory update:", event)
	go func() {
		if err := p.kafka.ProduceMessage(producer.PRODUCT_TOPIC, event); err != nil {
			log.Println("Error producing Kafka message:", err)
		}

	}()

	return nil
}

func (p *ProductService) UpdateInventoryEvent(product_id int) error {
	product, err := p.repo.GetProductById(product_id)
	if err != nil {
		return err
	}

	return p.repo.UpdateInventoryEvent(product)

}
