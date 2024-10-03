package models

type Address struct {
	Village  string `gorm:"type:varchar(255);not null" json:"village"`
	City     string `gorm:"type:varchar(255);not null" json:"city"`
	District string `gorm:"type:varchar(255);not null" json:"district"`
	State    string `gorm:"type:varchar(255);not null" json:"state"`
}

type UserCatalog struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName    string  `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName     string  `gorm:"type:varchar(255);not null" json:"last_name"`
	MobileNumber string  `gorm:"type:varchar(255);unique;not null" json:"mobile_number"`
	Email        string  `gorm:"type:varchar(255);unique;not null" json:"email"`
	Address      Address `gorm:"embedded" json:"address"`
}
type UserAddress struct {
	Village  string `gorm:"type:varchar(255);not null" json:"village"`
	City     string `gorm:"type:varchar(255);not null" json:"city"`
	District string `gorm:"type:varchar(255);not null" json:"district"`
	State    string `gorm:"type:varchar(255);not null" json:"state"`
}

type Data struct {
	ID           uint        `json:"id"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	MobileNumber string      `json:"mobile_number"`
	Email        string      `json:"email"`
	Address      UserAddress `json:"address"`
}

type UserEvent struct {
	EventType string `json:"event_type"`
	Data      Data   `json:"data"`
}
