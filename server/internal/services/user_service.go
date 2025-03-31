package services

import (
	"gorm.io/gorm"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/repositories"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (service *UserService) CreateUser(user *entity.User) error {
	userRepository := repositories.NewUserRepository(service.DB)
	return userRepository.CreateUser(user)
}

func (service *UserService) GetUser(userId string) (*entity.User, error) {
	userRepository := repositories.NewUserRepository(service.DB)
	return userRepository.GetUserById(userId)
}

func (service *UserService) GetUsers() ([]entity.User, error) {
	userRepository := repositories.NewUserRepository(service.DB)
	return userRepository.GetAllUsers()
}

func (service *UserService) UpdateUser(user *entity.User) error {
	userRepository := repositories.NewUserRepository(service.DB)
	return userRepository.UpdateUser(user)
}

func (service *UserService) DeleteUser(user *entity.User) error {
	userRepository := repositories.NewUserRepository(service.DB)
	return userRepository.DeleteUser(user)
}
