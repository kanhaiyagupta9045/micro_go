package repository

import (
	"errors"
	"time"

	"github.com/kanhaiyagupta9045/pratilipi/user/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(dburl string) (*UserRepository, error) {
	DB, err := gorm.Open(postgres.Open(dburl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return &UserRepository{
		db: DB,
	}, nil
}

func (r *UserRepository) RegisterUser(user *model.User) error {

	var existinguser model.User

	if err := r.db.Where("email = ? OR mobile_number = ?", user.Email, user.MobileNumber).First(&existinguser).Error; err == nil {
		if existinguser.Email == user.Email {
			return errors.New("email already exist")
		}
		if existinguser.MobileNumber == user.MobileNumber {
			return errors.New("phone number already exist")
		}
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	res := r.db.Create(&user)

	if res.RowsAffected == 0 {
		return errors.New("error while inserting user")
	}

	return nil
}

func (r *UserRepository) GetAllUser() ([]model.User, error) {
	var users []model.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user does n't exist")
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(user model.UpdateData) error {
	var existinguser model.User

	if err := r.db.Where("email = ?", user.Email).First(&existinguser).Error; err != nil {
		return errors.New("user doesn't exist")
	}

	if user.MobileNumber != "" {

		var userWithSameMobile model.User
		if err := r.db.Where("mobile_number = ?", user.MobileNumber).First(&userWithSameMobile).Error; err == nil {
			return errors.New("mobile number already exists")
		}
		existinguser.MobileNumber = user.MobileNumber
	}

	if user.FirstName != "" {
		existinguser.FirstName = user.FirstName
	}

	if user.LastName != "" {
		existinguser.LastName = user.LastName
	}

	res := r.db.Save(&existinguser)

	if res.Error != nil {
		return res.Error
	}

	return nil
}
