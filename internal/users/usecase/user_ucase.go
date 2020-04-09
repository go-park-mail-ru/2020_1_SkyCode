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
	user.SetUser()
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

func (uUC *UserUseCase) UpdatePhoneNumber(id uint64, phone string) error {
	user := &models.User{
		ID: id,
		Phone: phone,
	}
	if err := uUC.userRepo.UpdatePhone(user); err != nil {
		return err
	}

	return nil
}

func (uUC *UserUseCase) UpdateAvatar(id uint64, path string) error {
	user := &models.User{
		ID: id,
		Avatar: path,
	}
	if err := uUC.userRepo.UpdateAvatar(user); err != nil {
		return err
	}

	return nil
}

func (uUC *UserUseCase) UpdatePassword(id uint64, password string) error {
	user := &models.User{
		ID: id,
		Password: password,
	}
	if err := uUC.userRepo.UpdatePassword(user); err != nil {
		return err
	}

	return nil
}

func (uUC *UserUseCase) GetUserById(userId uint64) (*models.User, error) {
	user := &models.User{
		ID: userId,
	}
	if err := uUC.userRepo.GetById(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uUC *UserUseCase) GetUserByPhone(phone string) (*models.User, error) {
	user := &models.User{
		Phone: phone,
	}
	if err := uUC.userRepo.GetByPhone(user); err != nil {
		return nil, err
	}

	return user, nil
}

