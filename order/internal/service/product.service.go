package service

import (
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/models"
)

func (o *OrderService) CreateOrder(product models.Product) error {
	return o.repo.CreateOrder(&product)
}

func (o *OrderService) UpdateInventory(product_id int, newStock int) error {
	return o.repo.UpdateInventory(uint(product_id), newStock)
}

func (o *OrderService) ProductById(product_id int) (*models.Product, error) {
	return o.repo.GetProductById(product_id)
}
