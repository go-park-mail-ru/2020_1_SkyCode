package users

import (
	"github.com/2020_1_Skycode/internal/models"
)

type Repository interface {
	InsertInto(user *models.User) error
	Update(user *models.User) error
	Get(user *models.User) error
}
