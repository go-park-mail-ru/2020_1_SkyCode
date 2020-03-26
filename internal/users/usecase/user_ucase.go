package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/users"
)

type UserUseCase struct {
	userRepo users.Repository
}

func NewUserUseCase(ur users.Repository) users.UseCase {
	return &UserUseCase{
		userRepo: ur,
	}
}

func (uUC *UserUseCase) CreateUser(user *models.User) error {
	if err := uUC.userRepo.InsertInto(user); err != nil {
		return err
	}

	return nil
}