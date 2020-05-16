package usecase

import (
	"context"
	"database/sql"
	"github.com/2020_1_Skycode/internal/geodata"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/product_tags"
	"github.com/2020_1_Skycode/internal/restaurant_points"
	"github.com/2020_1_Skycode/internal/restaurants"
	protobuf_admin_rest "github.com/2020_1_Skycode/internal/restaurants/delivery/protobuf"
	"github.com/2020_1_Skycode/internal/restaurants_tags"
	"github.com/2020_1_Skycode/internal/reviews"
	"github.com/2020_1_Skycode/internal/tools"
	"google.golang.org/grpc"
	"math"
)

type RestaurantWithProtoUseCase struct {
	restaurantRepo  restaurants.Repository
	reviewsRepo     reviews.Repository
	geoDataRepo     geodata.Repository
	restPointsRepo  restaurant_points.Repository
	restTagsRepo    restaurants_tags.Repository
	productTagsRepo product_tags.Repository
	adminManager    protobuf_admin_rest.RestaurantAdminWorkerClient
}

func NewRestaurantsWithProtoUseCase(rr restaurants.Repository, rpr restaurant_points.Repository,
	rvr reviews.Repository, gdr geodata.Repository, rtr restaurants_tags.Repository,
	ptr product_tags.Repository, conn *grpc.ClientConn) restaurants.UseCase {
	return &RestaurantWithProtoUseCase{
		restaurantRepo:  rr,
		reviewsRepo:     rvr,
		geoDataRepo:     gdr,
		restPointsRepo:  rpr,
		restTagsRepo:    rtr,
		productTagsRepo: ptr,
		adminManager:    protobuf_admin_rest.NewRestaurantAdminWorkerClient(conn),
	}
}

func (rUC *RestaurantWithProtoUseCase) GetRestaurants(
	count uint64, page uint64, tagID uint64) ([]*models.Restaurant, uint64, error) {
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

func (rUC *RestaurantWithProtoUseCase) GetRestaurantByID(id uint64) (*models.Restaurant, error) {
	rest, err := rUC.restaurantRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return rest, nil
}

func (rUC *RestaurantWithProtoUseCase) CreateRestaurant(rest *models.Restaurant) error {
	_, err := rUC.adminManager.CreateRestaurant(
		context.Background(),
		&protobuf_admin_rest.ProtoRestaurant{
			ManagerID:   rest.ManagerID,
			Name:        rest.Name,
			Description: rest.Description,
			ImagePath:   rest.Image,
		})

	if err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) UpdateRestaurant(rest *models.Restaurant) error {
	answ, err := rUC.adminManager.UpdateRestaurant(
		context.Background(),
		&protobuf_admin_rest.ProtoRestaurant{
			ID:          rest.ID,
			Name:        rest.Name,
			Description: rest.Description,
		})

	if err != nil {
		return err
	}

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestaurantNotFoundError
		}
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) UpdateImage(restID uint64, filename string) error {
	answ, err := rUC.adminManager.UpdateRestaurantImage(
		context.Background(),
		&protobuf_admin_rest.ProtoImage{
			ID:        restID,
			ImagePath: filename,
		})

	if err != nil {
		return err
	}

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestaurantNotFoundError
		}
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) Delete(restID uint64) error {
	answ, err := rUC.adminManager.DeleteRestaurant(
		context.Background(),
		&protobuf_admin_rest.ProtoID{
			ID: restID,
		})

	if err != nil {
		return err
	}

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestaurantNotFoundError
		}
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) AddPoint(p *models.RestaurantPoint) error {
	answ, err := rUC.adminManager.CreatePoint(
		context.Background(),
		&protobuf_admin_rest.ProtoPoint{
			Address: p.Address,
			RestID:  p.RestID,
			Radius:  float32(p.ServiceRadius),
		})

	if err != nil {
		return err
	}

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestaurantNotFoundError
		}
		if answ.ID == tools.AddressNotHouse {
			return tools.ApiNotHouseAnswerError
		}
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) GetRestaurantsInServiceRadius(
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

func (rUC *RestaurantWithProtoUseCase) GetPoints(restID, count, page uint64) ([]*models.RestaurantPoint, uint64, error) {
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

func (rUC *RestaurantWithProtoUseCase) AddReview(review *models.Review) error {
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

func (rUC *RestaurantWithProtoUseCase) GetReviews(restID, userID, count, page uint64) (
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

func (rUC *RestaurantWithProtoUseCase) AddTag(restID, tagID uint64) error {
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

func (rUC *RestaurantWithProtoUseCase) DeleteTag(restID, tagID uint64) error {
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

func (rUC *RestaurantWithProtoUseCase) GetProductTagsByID(restID uint64) ([]*models.ProductTag, error) {
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

func (rUC *RestaurantWithProtoUseCase) AddProductTag(tag *models.ProductTag) error {
	if err := rUC.productTagsRepo.InsertInto(tag); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) DeleteProductTag(ID uint64) error {
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
