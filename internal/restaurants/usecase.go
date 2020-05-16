package restaurants

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	GetRestaurants(count uint64, page uint64, tagID uint64) ([]*models.Restaurant, uint64, error)
	GetRestaurantByID(id uint64) (*models.Restaurant, error)
	CreateRestaurant(rest *models.Restaurant) error
	UpdateRestaurant(rest *models.Restaurant) error
	UpdateImage(restID uint64, filename string) error
	Delete(restID uint64) error

	AddPoint(p *models.RestaurantPoint) error
	GetPoints(restID, count, page uint64) ([]*models.RestaurantPoint, uint64, error)
	GetRestaurantsInServiceRadius(address string, count, page, tagID uint64) ([]*models.Restaurant, uint64, error)

	AddReview(review *models.Review) error
	GetReviews(restID, userID, count, page uint64) ([]*models.Review, *models.Review, uint64, error)

	AddTag(restID, tagID uint64) error
	DeleteTag(restID, tagID uint64) error

	GetProductTagsByID(restID uint64) ([]*models.ProductTag, error)
	AddProductTag(tag *models.ProductTag) error
	DeleteProductTag(ID uint64) error
}
