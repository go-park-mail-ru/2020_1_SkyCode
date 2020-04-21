package reviews

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	GetReview(uint64) (*models.Review, error)
	GetUserReviews(uint64, uint64, uint64) ([]*models.Review, uint64, error)
	UpdateReview(*models.Review, *models.User) error
	DeleteReview(uint64, *models.User) error
}
