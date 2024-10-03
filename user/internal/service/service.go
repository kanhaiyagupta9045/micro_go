package service

import (
	"log"

	"github.com/kanhaiyagupta9045/pratilipi/user/internal/helpers"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/kafka"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/model"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/repository"
)

type UserService struct {
	repo  *repository.UserRepository
	kafka *kafka.Producer
}

func NewUserService(repo *repository.UserRepository, kafka *kafka.Producer) *UserService {
	return &UserService{repo: repo, kafka: kafka}
}

func (s *UserService) CreateUser(user *model.User) error {

	user.Password = helpers.HashPassPassword(user.Password)

	if err := s.repo.RegisterUser(user); err != nil {
		return err
	}

	modifedAddress := model.ModifiedAddress{
		Village:  user.Address.Village,
		City:     user.Address.City,
		District: user.Address.District,
		State:    user.Address.State,
	}

	data := model.Data{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		MobileNumber: user.MobileNumber,
		Email:        user.Email,
		Address:      modifedAddress,
	}
	event := model.UserEvent{
		EventType: "User Registered",
		Data:      data,
	}

	log.Println("Emitting event:", event)
	go func() {
		if err := s.kafka.ProduceMessage(kafka.USER_TOPIC, event); err != nil {
			log.Println("Error producing Kafka message:", err)
		}

	}()

	return nil
}

func (s *UserService) GetAllUser() ([]model.User, error) {

	users, err := s.repo.GetAllUser()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.repo.GetUserById(id)

	if err != nil {
		return nil, err
	}

	return user, nil

}
func (s *UserService) LoginUser(logindata model.LoginData) (*model.User, error) {

	user, err := s.repo.GetUserByEmail(logindata.Email)
	if err != nil {
		return nil, err
	}

	ok, err := helpers.VerifyPassword(logindata.Password, user.Password)

	if !ok {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateProfile(id int, updateddata *model.Data) error {
	updateddata.ID = uint(id)
	if err := s.repo.UpdateUser(id, updateddata); err != nil {
		return err
	}

	event := model.UserEvent{
		EventType: "User Profile Updated",
		Data:      *updateddata,
	}

	log.Println("Emitting Update User  event:", event)
	go func() {
		if err := s.kafka.ProduceMessage(kafka.USER_TOPIC, event); err != nil {
			log.Println("Error producing Kafka message:", err)
		}

	}()
	return nil
}
