package repository

import (
	"fmt"

	"github.com/kanhaiyagupta9045/pratilipi/order/internal/models"
	"gorm.io/gorm"
)

func (o *OrderRepositry) PlaceOrder(customerID, productID int, orderItem models.OrderItem) error {
	order := models.Order{
		CustomerID: uint(customerID),
		Amount:     orderItem.Price,
		Status:     "PENDING",
	}
	return o.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		orderItem.OrderID = order.OrderID
		if err := tx.Create(&orderItem).Error; err != nil {
			return err
		}

		return nil
	})
}

func (o *OrderRepositry) GetOrderById(orderID int) (*models.Order, error) {
	var order models.Order
	if err := o.db.Preload("OrderItem").First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *OrderRepositry) ShipOrder(order_id int) error {
	order, err := o.GetOrderById(order_id)
	if err != nil {
		return err
	}
	if order.Status == "Shipped" || order.Status == "Canceled" {
		return fmt.Errorf("order cannot be shipped. Current status: %s", order.Status)
	}
	order.Status = "Shipped"
	return o.db.Save(&order).Error

}
