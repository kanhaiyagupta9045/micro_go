package service

import (
	"log"

	"github.com/kanhaiyagupta9045/pratilipi/order/internal/models"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/producer"
)

func (o *OrderService) PlaceOrder(customer_id int, product_id int, orderItem models.OrderItem) error {

	err := o.repo.PlaceOrder(customer_id, product_id, orderItem)
	event := models.OrderEvent{
		EventType: "Order Placed",
		ProductID: product_id,
	}
	if err == nil {
		go func() {
			if err := o.producer.ProduceMessage(producer.Order_Topic, event); err != nil {
				log.Println("Error producing Kafka message:", err)
			}
		}()
	}
	return err
}

func (o *OrderService) GetOrderByID(order_id int) (*models.Order, error) {
	return o.repo.GetOrderById(order_id)
}

func (o *OrderService) ShipOrder(order_id int) error {
	return o.repo.ShipOrder(order_id)
}
