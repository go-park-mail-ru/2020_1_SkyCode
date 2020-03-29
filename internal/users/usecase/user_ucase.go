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

func (uUC *UserUseCase) UpdateBio(user *models.User) error {
	if err := uUC.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}

func (uUC *UserUseCase) GetUser(user *models.User) error {
	if err := uUC.userRepo.Get(user); err != nil {
		return err
	}

	return nil
}

