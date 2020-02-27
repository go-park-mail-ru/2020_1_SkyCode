package users

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	CreateUser(user *models.User) error
	UpdateBio(user *models.User) error
	GetUserByPhone(phone string) (*models.User, error)
	GetUserById(userId uint64) (*models.User, error)
	UpdatePhoneNumber(id uint64, phone string) error
	UpdatePassword(id uint64, password string) error
	UpdateAvatar(id uint64, path string) error
}
