package models

import "time"

type Order struct {
	OrderID    uint      `gorm:"primaryKey;autoIncrement" json:"order_id"`
	CustomerID uint      `gorm:"not null" json:"customer_id"`
	Amount     float64   `gorm:"not null" json:"amount"`
	Status     string    `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	OrderItem  OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"order_items"`
}
type OrderItem struct {
	OrderItemID uint    `gorm:"primaryKey;autoIncrement" json:"order_item_id"`
	OrderID     uint    `gorm:"not null" json:"order_id"`
	ProductID   uint    `gorm:"not null" json:"product_id"`
	Price       float64 `gorm:"not null" json:"price"`
}

type OrderEvent struct {
	EventType string `json:"event_type"`
	ProductID int    `json:"product_id"`
}
