package sessions

import (
	"github.com/2020_1_Skycode/internal/models"
)

type Repository interface {
	InsertInto(session *models.Session) error
	Delete(session *models.Session) error
}
