package usecase

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/reviews"
	"github.com/2020_1_Skycode/internal/tools"
	"math"
)

type ReviewsUseCase struct {
	reviewsRepo reviews.Repository
}

func NewReviewsUseCase(rr reviews.Repository) reviews.UseCase {
	return &ReviewsUseCase{
		reviewsRepo: rr,
	}
}

func (rUC *ReviewsUseCase) GetReview(id uint64) (*models.Review, error) {
	returnReview, err := rUC.reviewsRepo.GetReviewByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tools.ReviewNotFoundError
		}

		return nil, err
	}

	return returnReview, nil
}

func (rUC *ReviewsUseCase) GetUserReviews(userID, count, page uint64) ([]*models.Review, uint64, error) {
	returnReviews, err := rUC.reviewsRepo.GetReviewsByUserID(userID, count, page)
	if err != nil {
		return nil, 0, err
	}

	total, err := rUC.reviewsRepo.GetReviewsCountByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	return returnReviews, total, nil
}

func (rUC *ReviewsUseCase) UpdateReview(r *models.Review, u *models.User) error {
	returnReview, err := rUC.reviewsRepo.GetReviewByID(r.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return tools.ReviewNotFoundError
		}

		return err
	}

	if u.ID != returnReview.Author.ID {
		return tools.DontEnoughRights
	}

	r.Rate = math.Round(r.Rate*100) / 100
	err = rUC.reviewsRepo.UpdateReview(r)
	if err != nil {
		return err
	}

	return nil
}

func (rUC *ReviewsUseCase) DeleteReview(id uint64, u *models.User) error {
	returnReview, err := rUC.reviewsRepo.GetReviewByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return tools.ReviewNotFoundError
		}

		return err
	}

	if u.ID != returnReview.Author.ID {
		return tools.DontEnoughRights
	}

	err = rUC.reviewsRepo.DeleteReview(id)
	if err != nil {
		return err
	}

	return nil
}
