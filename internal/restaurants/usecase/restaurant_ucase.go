package usecase

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/geodata"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/product_tags"
	"github.com/2020_1_Skycode/internal/restaurant_points"
	"github.com/2020_1_Skycode/internal/restaurants"
	"github.com/2020_1_Skycode/internal/restaurants_tags"
	"github.com/2020_1_Skycode/internal/reviews"
	"github.com/2020_1_Skycode/internal/tools"
	"math"
)

type RestaurantUseCase struct {
	restaurantRepo  restaurants.Repository
	reviewsRepo     reviews.Repository
	geoDataRepo     geodata.Repository
	restPointsRepo  restaurant_points.Repository
	restTagsRepo    restaurants_tags.Repository
	productTagsRepo product_tags.Repository
}

func NewRestaurantsUseCase(rr restaurants.Repository, rpr restaurant_points.Repository,
	rvr reviews.Repository, gdr geodata.Repository, rtr restaurants_tags.Repository,
	ptr product_tags.Repository) *RestaurantUseCase {
	return &RestaurantUseCase{
		restaurantRepo:  rr,
		reviewsRepo:     rvr,
		geoDataRepo:     gdr,
		restPointsRepo:  rpr,
		restTagsRepo:    rtr,
		productTagsRepo: ptr,
	}
}

func (rUC *RestaurantUseCase) GetRestaurants(count uint64, page uint64, tagID uint64) ([]*models.Restaurant, uint64, error) {
	if tagID != 0 {
		if _, err := rUC.restTagsRepo.GetByID(tagID); err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, tools.RestTagNotFound
			}
		}
	}

	restList, total, err := rUC.restaurantRepo.GetAll(count, page, tagID)
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

func (rUC *RestaurantUseCase) AddPoint(p *models.RestaurantPoint) error {
	if _, err := rUC.restaurantRepo.GetByID(p.RestID); err != nil {
		if err == sql.ErrNoRows {
			return tools.RestaurantNotFoundError
		}
	}

	geoPoint, err := rUC.geoDataRepo.GetGeoPosByAddress(p.Address)
	if err != nil {
		return err
	}

	p.MapPoint = geoPoint

	err = rUC.restPointsRepo.InsertInto(p)
	if err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) GetRestaurantsInServiceRadius(
	address string, count, page, tagID uint64) ([]*models.Restaurant, uint64, error) {
	pos, err := rUC.geoDataRepo.GetGeoPosByAddress(address)
	if err != nil {
		return nil, 0, err
	}

	if tagID != 0 {
		if _, err := rUC.restTagsRepo.GetByID(tagID); err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, tools.RestTagNotFound
			}
		}
	}

	returnRests, total, err := rUC.restaurantRepo.GetAllInServiceRadius(pos, count, page, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Restaurant{}, 0, nil
		}
		return nil, 0, err
	}
	for _, rest := range returnRests {
		rest.Points = []*models.RestaurantPoint{}
		closerPoint, err := rUC.restPointsRepo.GetCloserPointByRestID(rest.ID, pos)
		if err != nil {
			return nil, 0, err
		}

		rest.Points = append(rest.Points, closerPoint)
	}

	return returnRests, total, nil
}

func (rUC *RestaurantUseCase) GetPoints(restID, count, page uint64) ([]*models.RestaurantPoint, uint64, error) {
	if _, err := rUC.restaurantRepo.GetByID(restID); err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, tools.RestaurantNotFoundError
		}
	}

	returnPoints, total, err := rUC.restPointsRepo.GetPointsByRestID(restID, count, page)
	if err != nil {
		return nil, 0, err
	}

	return returnPoints, total, nil
}

func (rUC *RestaurantUseCase) AddReview(review *models.Review) error {
	if _, err := rUC.restaurantRepo.GetByID(review.RestID); err != nil {
		if err == sql.ErrNoRows {
			return tools.RestaurantNotFoundError
		}
	}

	if curReview, err := rUC.reviewsRepo.GetRestaurantReviewByUser(
		review.RestID, review.Author.ID); err != nil || curReview != nil {
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

func (rUC *RestaurantUseCase) AddTag(restID, tagID uint64) error {
	if _, err := rUC.restTagsRepo.GetByID(tagID); err != nil {
		if err == sql.ErrNoRows {
			return tools.RestTagNotFound
		}

		return err
	}

	if _, err := rUC.restaurantRepo.GetByID(restID); err != nil {
		if err == sql.ErrNoRows {
			return tools.RestaurantNotFoundError
		}
	}

	check, err := rUC.restTagsRepo.CheckTagRestRelation(restID, tagID)
	if err != nil {
		return err
	}

	if check {
		return tools.TagRestComboAlreadyExist
	}

	if err := rUC.restTagsRepo.CreateTagRestRelation(restID, tagID); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) DeleteTag(restID, tagID uint64) error {
	if _, err := rUC.restTagsRepo.GetByID(tagID); err != nil {
		if err == sql.ErrNoRows {
			return tools.RestTagNotFound
		}

		return err
	}

	if _, err := rUC.restaurantRepo.GetByID(restID); err != nil {
		if err == sql.ErrNoRows {
			return tools.RestaurantNotFoundError
		}
	}

	check, err := rUC.restTagsRepo.CheckTagRestRelation(restID, tagID)
	if err != nil {
		return err
	}

	if !check {
		return tools.TagRestComboDoesntExist
	}

	if err := rUC.restTagsRepo.DeleteTagRestRelation(restID, tagID); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) GetProductTagsByID(restID uint64) ([]*models.ProductTag, error) {
	if _, err := rUC.restaurantRepo.GetByID(restID); err != nil {
		if err == sql.ErrNoRows {
			return nil, tools.RestaurantNotFoundError
		}

		return nil, err
	}

	tags, err := rUC.productTagsRepo.GetByRestID(restID)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (rUC *RestaurantUseCase) AddProductTag(tag *models.ProductTag) error {
	if err := rUC.productTagsRepo.InsertInto(tag); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) DeleteProductTag(ID uint64) error {
	if _, err := rUC.productTagsRepo.GetByID(ID); err != nil {
		if err == sql.ErrNoRows {
			return tools.ProductTagNotFound
		}

		return err
	}

	if err := rUC.productTagsRepo.Delete(ID); err != nil {
		return err
	}

	return nil
}
