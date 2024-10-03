package service

import (
	"log"

	"github.com/kanhaiyagupta9045/pratilipi/order/internal/models"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/producer"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/repository"
)

type OrderService struct {
	repo     *repository.OrderRepositry
	producer *producer.Producer
}

func NewOrderService(repo *repository.OrderRepositry, producer *producer.Producer) *OrderService {
	return &OrderService{
		repo:     repo,
		producer: producer,
	}
}

func (o *OrderService) CreateUser(eventData models.Data) error {

	user := models.User{
		ID:           eventData.ID,
		FirstName:    eventData.FirstName,
		LastName:     eventData.LastName,
		MobileNumber: eventData.MobileNumber,
		Email:        eventData.Email,
		Address: models.Address{
			Village:  eventData.Address.Village,
			City:     eventData.Address.City,
			District: eventData.Address.District,
			State:    eventData.Address.State,
		},
	}

	if err := o.repo.CreateUser(&user); err != nil {
		log.Println("Error creating user:", err)
		return err
	}
	return nil
}

func (o *OrderService) UpdateUser(eventData models.Data) error {

	if err := o.repo.UpdateUser(int(eventData.ID), eventData); err != nil {
		return err
	}

	return nil
}

func (o *OrderService) GetUserById(id int) (*models.User, error) {

	return o.repo.GetUserById(id)
}
