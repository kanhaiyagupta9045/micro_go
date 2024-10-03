package repository

import (
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type OrderRepositry struct {
	db *gorm.DB
}

func NewOrderRepositry(connectionstring string) (*OrderRepositry, error) {
	DB, err := gorm.Open(postgres.Open(connectionstring), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	if err := DB.AutoMigrate(&models.UserCatalog{}, &models.ProductCatalog{}, &models.InventoryCatalog{}, &models.Order{}, &models.OrderItem{}); err != nil {
		return nil, err
	}

	return &OrderRepositry{
		db: DB,
	}, nil
}

func (o *OrderRepositry) CreateUser(user *models.UserCatalog) error {
	if err := o.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (o *OrderRepositry) UpdateUser(id int, eventData models.Data) error {
	var existinguser models.UserCatalog
	if err := o.db.Where("id = ?", id).First(&existinguser).Error; err != nil {
		return err
	}

	if eventData.Address.Village != "" {
		existinguser.Address.Village = eventData.Address.Village
	}
	if eventData.Address.City != "" {
		existinguser.Address.City = eventData.Address.City
	}
	if eventData.Address.District != "" {
		existinguser.Address.District = eventData.Address.District
	}
	if eventData.Address.State != "" {
		existinguser.Address.State = eventData.Address.State
	}
	if eventData.Email != "" {
		existinguser.Email = eventData.Email
	}
	if eventData.FirstName != "" {
		existinguser.FirstName = eventData.FirstName
	}
	if eventData.LastName != "" {
		existinguser.LastName = eventData.LastName
	}
	if eventData.MobileNumber != "" {
		existinguser.MobileNumber = eventData.MobileNumber
	}

	result := o.db.Save(&existinguser)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

func (o *OrderRepositry) GetUserById(id int) (*models.UserCatalog, error) {
	var user models.UserCatalog
	if err := o.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
