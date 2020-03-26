package users

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	CreateUser(user *models.User) error
}