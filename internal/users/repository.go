package users

import (
	"github.com/2020_1_Skycode/internal/models"
)

type Repository interface {
	InsertInto(user *models.User) error
	Update(user *models.User) error
	UpdatePhone(user *models.User) error
	UpdatePassword(user *models.User) error
	UpdateAvatar(user *models.User) error
	GetById(user *models.User) error
	GetByPhone(user *models.User) error
}
