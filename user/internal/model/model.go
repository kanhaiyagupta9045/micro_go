package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName    string `gorm:"first_name;not null" json:"first_name" validate:"required"`
	LastName     string `gorm:"last_name;not null" json:"last_name" validate:"required"`
	MobileNumber string `gorm:"mobile_number;unique;not null" json:"mobile_number" validate:"required"`
	Email        string `gorm:"email;unique;not null" json:"email" validate:"required,email"`
	Password     string `gorm:"password;not null" validate:"required,min=8" json:"-"`
	UserType     string `gorm:"usertype;not null" json:"user_type" validate:"required"`
}

type LoginData struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateData struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MobileNumber string `json:"mobile_number"`
	Email        string `json:"email" validate:"required,email"`
}
