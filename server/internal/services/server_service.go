package services

import (
	"gorm.io/gorm"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/repositories"
)

type ServerService struct {
	DB *gorm.DB
}

func NewServerService(db *gorm.DB) *ServerService {
	return &ServerService{DB: db}
}

func (service *ServerService) CreateServer(server *entity.Server) error {
	serverRepository := repositories.NewServerRepository(service.DB)
	return serverRepository.CreateServer(server)
}

func (service *ServerService) GetServer(serverId string) (*entity.Server, error) {
	serverRepository := repositories.NewServerRepository(service.DB)
	return serverRepository.GetServer(serverId)
}

func (service *ServerService) GetServers() ([]entity.Server, error) {
	serverRepository := repositories.NewServerRepository(service.DB)
	return serverRepository.GetServers()
}

func (service *ServerService) UpdateServer(server *entity.Server) error {
	serverRepository := repositories.NewServerRepository(service.DB)
	return serverRepository.UpdateServer(server)
}

func (service *ServerService) DeleteServer(server *entity.Server) error {
	serverRepository := repositories.NewServerRepository(service.DB)
	return serverRepository.DeleteServer(server)
}

func (service *ServerService) GetServerMembers(serverId string) ([]entity.User, error) {
	serverRepository := repositories.NewServerRepository(service.DB)
	return serverRepository.GetServerMembers(serverId)
}

func (service *ServerService) GetServerChannels(serverId string) ([]entity.Channel, error) {
	serverRepository := repositories.NewServerRepository(service.DB)
	return serverRepository.GetServerChannels(serverId)
}
