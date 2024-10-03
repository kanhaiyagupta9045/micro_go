package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName    string  `gorm:"type:varchar(255);not null" json:"first_name" validate:"required"`
	LastName     string  `gorm:"type:varchar(255);not null" json:"last_name" validate:"required"`
	MobileNumber string  `gorm:"type:varchar(255);unique;not null" json:"mobile_number" validate:"required"`
	Email        string  `gorm:"type:varchar(255);unique;not null" json:"email" validate:"required,email"`
	Password     string  `gorm:"type:varchar(255);not null" json:"password" validate:"required,min=8"`
	UserType     string  `gorm:"type:varchar(50);not null" json:"user_type" validate:"required"`
	Address      Address `gorm:"embedded" json:"address"`
}

type Address struct {
	Village  string `gorm:"type:varchar(255);not null" json:"village" validate:"required"`
	City     string `gorm:"type:varchar(255);not null" json:"city" validate:"required"`
	District string `gorm:"type:varchar(255);not null" json:"district" validate:"required"`
	State    string `gorm:"type:varchar(255);not null" json:"state" validate:"required"`
}

type LoginData struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Data struct {
	ID           uint            `json:"id"`
	FirstName    string          `json:"first_name"`
	LastName     string          `json:"last_name"`
	MobileNumber string          `json:"mobile_number"`
	Email        string          `json:"email"`
	Address      ModifiedAddress `json:"address"`
}
type ModifiedAddress struct {
	Village  string `json:"village"`
	City     string `json:"city"`
	District string `json:"district"`
	State    string `json:"state"`
}

type UserEvent struct {
	EventType string `json:"event_type"`
	Data      Data   `json:"data"`
}
