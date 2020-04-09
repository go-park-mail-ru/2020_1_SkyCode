package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_sessions "github.com/2020_1_Skycode/internal/sessions/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUseCase_DeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockSessionsRepo := mock_sessions.NewMockRepository(ctrl)

	testSession := &models.Session{
		ID: uint64(1),
	}

	mockSessionsRepo.EXPECT().Delete(testSession).Return(nil)
	sessionUCase := NewSessionUseCase(mockSessionsRepo)

	if err := sessionUCase.DeleteSession(testSession.ID); err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}

func TestUseCase_GetSession(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockSessionsRepo := mock_sessions.NewMockRepository(ctrl)

	testSession := &models.Session{
		Token: "test-token",
	}

	mockSessionsRepo.EXPECT().Get(testSession).Return(nil)
	sessionUCase := NewSessionUseCase(mockSessionsRepo)

	resultSess, err := sessionUCase.GetSession(testSession.Token)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, testSession, resultSess)
}

func TestUseCase_StoreSession(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockSessionsRepo := mock_sessions.NewMockRepository(ctrl)

	testSession := &models.Session{
		UserId: uint64(1),
		Token:  "test-token",
	}

	mockSessionsRepo.EXPECT().InsertInto(testSession).Return(nil)
	sessionUCase := NewSessionUseCase(mockSessionsRepo)

	err := sessionUCase.StoreSession(testSession)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}
