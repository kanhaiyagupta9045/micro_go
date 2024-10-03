package repository

import (
	"errors"
	"fmt"

	"github.com/kanhaiyagupta9045/pratilipi/order/internal/models"
	"gorm.io/gorm"
)

func (o *OrderRepositry) CreateOrder(product *models.Product) error {
	if err := o.db.Create(&product).Error; err != nil {
		return err
	}
	return nil
}

func (o *OrderRepositry) UpdateInventory(productID uint, newStock int) error {
	err := o.db.Transaction(func(tx *gorm.DB) error {
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

func (p *OrderRepositry) GetProductById(product_id int) (*models.Product, error) {
	var product *models.Product
	if err := p.db.Preload("Inventory").Where("product_id = ?", product_id).First(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}
