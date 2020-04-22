package usecase

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants"
	"github.com/2020_1_Skycode/internal/reviews"
	"github.com/2020_1_Skycode/internal/tools"
	"math"
)

type RestaurantUseCase struct {
	restaurantRepo restaurants.Repository
	reviewsRepo    reviews.Repository
}

func NewRestaurantsUseCase(rr restaurants.Repository, rvr reviews.Repository) *RestaurantUseCase {
	return &RestaurantUseCase{
		restaurantRepo: rr,
		reviewsRepo:    rvr,
	}
}

func (rUC *RestaurantUseCase) GetRestaurants(count uint64, page uint64) ([]*models.Restaurant, uint64, error) {
	restList, total, err := rUC.restaurantRepo.GetAll(count, page)
	if err != nil {
		return nil, total, err
	}

	return restList, total, nil
}

func (rUC *RestaurantUseCase) GetRestaurantByID(id uint64) (*models.Restaurant, error) {
	rest, err := rUC.restaurantRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return rest, nil
}

func (rUC *RestaurantUseCase) CreateRestaurant(rest *models.Restaurant) error {
	if err := rUC.restaurantRepo.InsertInto(rest); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) UpdateRestaurant(rest *models.Restaurant) error {
	updRestaurant := &models.Restaurant{
		ID:          rest.ID,
		Name:        rest.Name,
		Description: rest.Description,
	}

	if err := rUC.restaurantRepo.Update(updRestaurant); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) UpdateImage(restID uint64, filename string) error {
	updRestaurant := &models.Restaurant{
		ID:    restID,
		Image: filename,
	}

	if err := rUC.restaurantRepo.UpdateImage(updRestaurant); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) Delete(restID uint64) error {
	if err := rUC.restaurantRepo.Delete(restID); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) AddReview(review *models.Review) error {
	if _, err := rUC.restaurantRepo.GetByID(review.RestID); err != nil {
		if err == sql.ErrNoRows {
			return tools.RestaurantNotFoundError
		}
	}

	if curReview, err := rUC.reviewsRepo.GetRestaurantReviewByUser(
		review.RestID, review.Author); err != nil || curReview != nil {
		if curReview != nil {
			return tools.ReviewAlreadyExists
		}

		return err
	}

	review.Rate = math.Round(review.Rate*100) / 100
	if err := rUC.reviewsRepo.CreateReview(review); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) GetReviews(restID, userID, count, page uint64) (
	[]*models.Review, *models.Review, uint64, error) {
	if _, err := rUC.restaurantRepo.GetByID(restID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, 0, tools.RestaurantNotFoundError
		}
	}

	returnReviews, err := rUC.reviewsRepo.GetReviewsByRestID(restID, count, page)
	if err != nil {
		return nil, nil, 0, err
	}

	curReview, err := rUC.reviewsRepo.GetRestaurantReviewByUser(restID, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, 0, err
	}

	total, err := rUC.reviewsRepo.GetReviewsCountByRestID(restID)
	if err != nil {
		return nil, nil, 0, err
	}

	return returnReviews, curReview, total, nil
}
