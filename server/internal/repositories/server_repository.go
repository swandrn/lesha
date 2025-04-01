// repositories/server_repository.go
package repositories

import (
	"gorm.io/gorm"
	"lesha.com/server/internal/entity"
)

type ServerRepository struct {
	DB *gorm.DB
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{DB: db}
}

func (repo *ServerRepository) CreateServer(server *entity.Server) error {
	return repo.DB.Create(server).Error
}
func (repo *ServerRepository) GetServer(serverId string) (*entity.Server, error) {
	var server entity.Server
	err := repo.DB.Where("id = ?", serverId).First(&server).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}
func (repo *ServerRepository) GetServers() ([]entity.Server, error) {
	var servers []entity.Server
	err := repo.DB.Find(&servers).Error
	if err != nil {
		return nil, err
	}
	return servers, nil
}
func (repo *ServerRepository) GetUserServers(userID uint) ([]entity.Server, error) {
	var servers []entity.Server
	err := repo.DB.Joins("JOIN user_servers ON servers.id = user_servers.server_id").
		Where("user_servers.user_id = ?", userID).
		Find(&servers).Error
	if err != nil {
		return nil, err
	}
	return servers, nil
}
func (repo *ServerRepository) UpdateServer(server *entity.Server) error {
	return repo.DB.Save(server).Error
}
func (repo *ServerRepository) DeleteServer(server *entity.Server) error {
	return repo.DB.Delete(server).Error
}
func (repo *ServerRepository) GetServerMembers(serverId string) ([]entity.User, error) {
	var members []entity.User
	err := repo.DB.Where("server_id = ?", serverId).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
func (repo *ServerRepository) GetServerChannels(serverId string) ([]entity.Channel, error) {
	var channels []entity.Channel
	err := repo.DB.Where("server_id = ?", serverId).Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}
