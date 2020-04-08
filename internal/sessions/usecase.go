package sessions

import (
	"github.com/2020_1_Skycode/internal/models"
)

type UseCase interface {
	StoreSession(user *models.Session) error
	DeleteSession(sessionId uint64) error
	GetSession(token string) (*models.Session, error)
}
