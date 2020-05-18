package usecase

import (
	"context"
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurant_points"
	protobuf_admin_rest "github.com/2020_1_Skycode/internal/restaurants/delivery/protobuf"
	"github.com/2020_1_Skycode/internal/tools"
	"google.golang.org/grpc"
)

type RestPointsWithProtoUseCase struct {
	RestPointsRepo restaurant_points.Repository
	adminManager   protobuf_admin_rest.RestaurantAdminWorkerClient
}

func NewRestPointsWithProtoUseCase(rpr restaurant_points.Repository, conn *grpc.ClientConn) restaurant_points.UseCase {
	return &RestPointsWithProtoUseCase{
		RestPointsRepo: rpr,
		adminManager:   protobuf_admin_rest.NewRestaurantAdminWorkerClient(conn),
	}
}

func (rpUC *RestPointsWithProtoUseCase) GetPoint(id uint64) (*models.RestaurantPoint, error) {
	returnPoint, err := rpUC.RestPointsRepo.GetPointByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tools.RestPointNotFound
		}

		return nil, err
	}

	return returnPoint, nil
}

func (rpUC *RestPointsWithProtoUseCase) GetAllPoints() ([]*models.RestaurantPoint, error) {
	returnPoints, err := rpUC.RestPointsRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return returnPoints, err
}

func (rpUC *RestPointsWithProtoUseCase) Delete(id uint64) error {
	answ, err := rpUC.adminManager.DeletePoint(
		context.Background(),
		&protobuf_admin_rest.ProtoID{ID: id})

	if err != nil {
		return err
	}

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestPointNotFound
		}
	}

	return nil
}
