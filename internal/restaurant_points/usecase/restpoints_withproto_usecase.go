package usecase

import (
	"context"
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurant_points"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/tools/protobuf/adminwork"
	"google.golang.org/grpc"
)

type RestPointsWithProtoUseCase struct {
	RestPointsRepo restaurant_points.Repository
	adminManager   adminwork.RestaurantAdminWorkerClient
}

func NewRestPointsWithProtoUseCase(rpr restaurant_points.Repository, conn *grpc.ClientConn) restaurant_points.UseCase {
	return &RestPointsWithProtoUseCase{
		RestPointsRepo: rpr,
		adminManager:   adminwork.NewRestaurantAdminWorkerClient(conn),
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
		&adminwork.ProtoID{ID: id})

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestPointNotFound
		}

		return err
	}

	return nil
}
