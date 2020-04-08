package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/sessions"
)

type UseCase struct {
	sessionRepo sessions.Repository
}

func NewSessionUseCase(sessionRepo sessions.Repository) *UseCase {
	return &UseCase{
		sessionRepo: sessionRepo,
	}
}

func (sUC *UseCase) StoreSession(session *models.Session) error {
	if err := sUC.sessionRepo.InsertInto(session); err != nil {
		return err
	}

	return nil
}

func (sUC *UseCase) GetSession(token string) (*models.Session, error) {
	currSession := &models.Session{
		Token: token,
	}

	if err := sUC.sessionRepo.Get(currSession); err != nil {
		return nil, err
	}

	return currSession, nil
}

func (sUC *UseCase) DeleteSession(sessionId uint64) error {
	currSession := &models.Session{
		ID: sessionId,
	}

	if err := sUC.sessionRepo.Delete(currSession); err != nil {
		return err
	}

	return nil
}
