package usecase

import (
	"context"
	"database/sql"
	"github.com/2020_1_Skycode/internal/geodata"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurant_points"
	"github.com/2020_1_Skycode/internal/restaurants"
	"github.com/2020_1_Skycode/internal/reviews"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/tools/protobuf/adminwork"
	"google.golang.org/grpc"
	"math"
)

type RestaurantWithProtoUseCase struct {
	restaurantRepo restaurants.Repository
	reviewsRepo    reviews.Repository
	geoDataRepo    geodata.Repository
	restPointsRepo restaurant_points.Repository
	adminManager   adminwork.RestaurantAdminWorkerClient
}

func NewRestaurantsWithProtoUseCase(rr restaurants.Repository, rpr restaurant_points.Repository,
	rvr reviews.Repository, gdr geodata.Repository, conn *grpc.ClientConn) restaurants.UseCase {
	return &RestaurantWithProtoUseCase{
		restaurantRepo: rr,
		reviewsRepo:    rvr,
		geoDataRepo:    gdr,
		restPointsRepo: rpr,
		adminManager:   adminwork.NewRestaurantAdminWorkerClient(conn),
	}
}

func (rUC *RestaurantWithProtoUseCase) GetRestaurants(count uint64, page uint64) ([]*models.Restaurant, uint64, error) {
	restList, total, err := rUC.restaurantRepo.GetAll(count, page)
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
	answ, err := rUC.adminManager.CreateRestaurant(
		context.Background(),
		&adminwork.ProtoRestaurant{
			ManagerID:   rest.ManagerID,
			Name:        rest.Name,
			Description: rest.Description,
			ImagePath:   rest.Image,
		})

	if answ.ID != tools.OK {
		if err != nil {
			return err
		}
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) UpdateRestaurant(rest *models.Restaurant) error {
	answ, err := rUC.adminManager.UpdateRestaurant(
		context.Background(),
		&adminwork.ProtoRestaurant{
			ID:          rest.ID,
			Name:        rest.Name,
			Description: rest.Description,
		})

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestaurantNotFoundError
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) UpdateImage(restID uint64, filename string) error {
	answ, err := rUC.adminManager.UpdateRestaurantImage(
		context.Background(),
		&adminwork.ProtoImage{
			ID:        restID,
			ImagePath: filename,
		})

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestaurantNotFoundError
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) Delete(restID uint64) error {
	answ, err := rUC.adminManager.DeleteRestaurant(
		context.Background(),
		&adminwork.ProtoID{
			ID: restID,
		})

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestaurantNotFoundError
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) AddPoint(p *models.RestaurantPoint) error {
	answ, err := rUC.adminManager.CreatePoint(
		context.Background(),
		&adminwork.ProtoPoint{
			Address: p.Address,
			RestID:  p.RestID,
			Radius:  float32(p.ServiceRadius),
		})

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestaurantNotFoundError
		}
		if answ.ID == tools.AddressNotHouse {
			return tools.ApiNotHouseAnswerError
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (rUC *RestaurantWithProtoUseCase) GetRestaurantsInServiceRadius(
	address string, count, page uint64) ([]*models.Restaurant, uint64, error) {
	pos, err := rUC.geoDataRepo.GetGeoPosByAddress(address)
	if err != nil {
		return nil, 0, err
	}

	returnRests, total, err := rUC.restaurantRepo.GetAllInServiceRadius(pos, count, page)
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
