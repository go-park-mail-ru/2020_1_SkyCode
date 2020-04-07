package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_users "github.com/2020_1_Skycode/internal/users/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserUseCase_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)

	testUser := &models.User{
		Email:     "testmail@m.ru",
		Password:  "1234",
		FirstName: "A",
		LastName:  "B",
		Phone:     "89765433221",
		Avatar:    "./default.jpg",
	}

	mockUserRepo.EXPECT().InsertInto(testUser).Return(nil)
	userUCase := NewUserUseCase(mockUserRepo)

	if err := userUCase.CreateUser(testUser); err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, "User", testUser.Role)
}

func TestUserUseCase_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)

	testUser := &models.User{
		ID: uint64(1),
	}

	mockUserRepo.EXPECT().GetById(testUser).Return(nil)
	userUCase := NewUserUseCase(mockUserRepo)

	resultUser, err := userUCase.GetUserById(testUser.ID)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, testUser, resultUser)
}

func TestUserUseCase_GetUserByPhone(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)

	testUser := &models.User{
		Phone: "89765433221",
	}

	mockUserRepo.EXPECT().GetByPhone(testUser).Return(nil)
	userUCase := NewUserUseCase(mockUserRepo)

	resultUser, err := userUCase.GetUserByPhone(testUser.Phone)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, testUser, resultUser)
}

func TestUserUseCase_UpdateAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)

	testUser := &models.User{
		ID:     uint64(1),
		Avatar: "./default.jpg",
	}

	mockUserRepo.EXPECT().UpdateAvatar(testUser).Return(nil)
	userUCase := NewUserUseCase(mockUserRepo)

	if err := userUCase.UpdateAvatar(testUser.ID, testUser.Avatar); err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}

func TestUserUseCase_UpdateBio(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)

	testUser := &models.User{
		ID:        uint64(1),
		Email:     "testmail@m.ru",
		FirstName: "A",
		LastName:  "B",
	}

	mockUserRepo.EXPECT().Update(testUser).Return(nil)
	userUCase := NewUserUseCase(mockUserRepo)

	if err := userUCase.UpdateBio(testUser); err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}

func TestUserUseCase_UpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)

	testUser := &models.User{
		ID:       uint64(1),
		Password: "1234",
	}

	mockUserRepo.EXPECT().UpdatePassword(testUser).Return(nil)
	userUCase := NewUserUseCase(mockUserRepo)

	if err := userUCase.UpdatePassword(testUser.ID, testUser.Password); err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}

func TestUserUseCase_UpdatePhoneNumber(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)

	testUser := &models.User{
		ID:    uint64(1),
		Phone: "89765433221",
	}

	mockUserRepo.EXPECT().UpdatePhone(testUser).Return(nil)
	userUCase := NewUserUseCase(mockUserRepo)

	if err := userUCase.UpdatePhoneNumber(testUser.ID, testUser.Phone); err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}
