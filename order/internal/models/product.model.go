package models

import "time"

type Product struct {
	ProductID   uint      `gorm:"primaryKey;autoIncrement" json:"product_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Inventory   Inventory `gorm:"foreignKey:ProductID;references:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"inventory"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Inventory struct {
	InventoryID uint `gorm:"primaryKey;autoIncrement" json:"inventory_id"`
	ProductID   uint `gorm:"not null" json:"product_id"`
	StockLevel  int  `gorm:"not null" json:"stock_level"`
}

type ProductCreatedEvent struct {
	EventType string  `json:"event_type"`
	Data      Product `json:"data"`
}

type InventoryUpdate struct {
	EventType  string `json:"event_type"`
	ProductID  int    `gorm:"type:uuid;not null" json:"product_id"`
	StockLevel int    `gorm:"not null" json:"stock_level"`
}
