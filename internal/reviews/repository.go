package reviews

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	GetRatingByRestID(uint64) (float64, error)
	GetReviewsByRestID(uint64, uint64, uint64) ([]*models.Review, error)
	GetReviewsCountByRestID(uint64) (uint64, error)
	GetReviewsByUserID(uint64, uint64, uint64) ([]*models.Review, error)
	GetReviewsCountByUserID(uint64) (uint64, error)
	CheckRestaurantReviewByUser(uint64, uint64) (bool, error)
	GetReviewByID(uint64) (*models.Review, error)
	CreateReview(*models.Review) error
	UpdateReview(*models.Review) error
	DeleteReview(uint64) error
}
