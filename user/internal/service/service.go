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

	err := s.kafka.ProduceMessage(kafka.USER_TOPIC, "USER Registered")
	if err != nil {
		log.Println(err.Error())
	}

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

func (s *UserService) UpdateProfile(updateddata model.UpdateData) error {

	if err := s.repo.UpdateUser(updateddata); err != nil {
		return err
	}
	if err := s.kafka.ProduceMessage(kafka.USER_TOPIC, "User Profile Updated"); err != nil {
		log.Println(err.Error())
	}
	return nil
}
