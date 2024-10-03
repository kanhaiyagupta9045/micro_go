package repository

import (
	"errors"
	"fmt"

	"github.com/kanhaiyagupta9045/product_service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ProductRepositry struct {
	db *gorm.DB
}

func NewProductRepository(connectionstring string) (*ProductRepositry, error) {
	DB, err := gorm.Open(postgres.Open(connectionstring), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	if err := DB.AutoMigrate(&models.Product{}, &models.Inventory{}); err != nil {
		return nil, err
	}

	return &ProductRepositry{
		db: DB,
	}, nil
}

func (p *ProductRepositry) InsertProduct(product *models.Product) (*models.Product, error) {
	fmt.Print(product)
	err := p.db.Create(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductRepositry) ListAllProducts() ([]models.Product, error) {
	var products []models.Product

	if err := p.db.Preload("Inventory").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil

}

func (p *ProductRepositry) GetProductById(product_id int) (*models.Product, error) {
	var product *models.Product
	if err := p.db.Preload("Inventory").Where("product_id = ?", product_id).First(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductRepositry) UpdateInventory(productID uint, newStock int) error {

	err := p.db.Transaction(func(tx *gorm.DB) error {
		var inventory models.Inventory

		if err := tx.Where("product_id = ?", productID).First(&inventory).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("inventory not found for product ID: %d", productID)
			}
			return err
		}

		inventory.StockLevel = newStock

		if err := tx.Save(&inventory).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (p *ProductRepositry) UpdateInventoryEvent(product *models.Product) error {

	err := p.db.Transaction(func(tx *gorm.DB) error {
		product.Inventory.StockLevel -= 1

		if err := tx.Save(&product).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
