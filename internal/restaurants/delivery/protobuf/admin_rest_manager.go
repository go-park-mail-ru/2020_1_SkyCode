package protobuf_admin_rest

import (
	"context"
	"database/sql"
	"github.com/2020_1_Skycode/internal/geodata"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/products"
	"github.com/2020_1_Skycode/internal/restaurant_points"
	"github.com/2020_1_Skycode/internal/restaurants"
	"github.com/2020_1_Skycode/internal/tools"
)

type AdminRestaurantManager struct {
	RestaurantRepo restaurants.Repository
	ProductRepo    products.Repository
	RestPointsRepo restaurant_points.Repository
	GeoDataRepo    geodata.Repository
}

func NewAdminRestaurantManager(rr restaurants.Repository, pr products.Repository,
	rpr restaurant_points.Repository, gdr geodata.Repository) *AdminRestaurantManager {
	return &AdminRestaurantManager{
		RestaurantRepo: rr,
		ProductRepo:    pr,
		RestPointsRepo: rpr,
		GeoDataRepo:    gdr,
	}
}

func (am *AdminRestaurantManager) CreateRestaurant(ctx context.Context, r *ProtoRestaurant) (*ErrorCode, error) {
	rest := &models.Restaurant{
		ManagerID:   r.ManagerID,
		Name:        r.Name,
		Description: r.Description,
		Image:       r.ImagePath,
	}

	if err := am.RestaurantRepo.InsertInto(rest); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}

func (am *AdminRestaurantManager) UpdateRestaurant(ctx context.Context, r *ProtoRestaurant) (*ErrorCode, error) {
	_, err := am.RestaurantRepo.GetByID(r.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ErrorCode{ID: tools.DoesntExist}, nil
		}
		return &ErrorCode{ID: tools.InternalError}, err
	}

	rest := &models.Restaurant{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	if err := am.RestaurantRepo.Update(rest); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}

func (am *AdminRestaurantManager) UpdateRestaurantImage(ctx context.Context, im *ProtoImage) (*ErrorCode, error) {
	_, err := am.RestaurantRepo.GetByID(im.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ErrorCode{ID: tools.DoesntExist}, nil
		}
		return &ErrorCode{ID: tools.InternalError}, err
	}

	rest := &models.Restaurant{
		ID:    im.ID,
		Image: im.ImagePath,
	}

	if err := am.RestaurantRepo.UpdateImage(rest); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}

func (am *AdminRestaurantManager) DeleteRestaurant(ctx context.Context, id *ProtoID) (*ErrorCode, error) {
	_, err := am.RestaurantRepo.GetByID(id.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ErrorCode{ID: tools.DoesntExist}, nil
		}
		return &ErrorCode{ID: tools.InternalError}, err
	}

	if err := am.RestaurantRepo.Delete(id.ID); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}

func (am *AdminRestaurantManager) CreateProduct(ctx context.Context, p *ProtoProduct) (*ErrorCode, error) {
	_, err := am.RestaurantRepo.GetByID(p.RestID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ErrorCode{ID: tools.DoesntExist}, nil
		}
		return &ErrorCode{ID: tools.InternalError}, err
	}

	prod := &models.Product{
		Name:   p.Name,
		Price:  p.Price,
		Image:  p.ImagePath,
		RestId: p.RestID,
	}

	if err := am.ProductRepo.InsertInto(prod); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}

func (am *AdminRestaurantManager) UpdateProduct(ctx context.Context, p *ProtoProduct) (*ErrorCode, error) {
	_, err := am.ProductRepo.GetProductByID(p.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ErrorCode{ID: tools.DoesntExist}, nil
		}
		return &ErrorCode{ID: tools.InternalError}, err
	}

	prod := &models.Product{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
	}

	if err := am.ProductRepo.Update(prod); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}

func (am *AdminRestaurantManager) UpdateProductImage(ctx context.Context, im *ProtoImage) (*ErrorCode, error) {
	_, err := am.ProductRepo.GetProductByID(im.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ErrorCode{ID: tools.DoesntExist}, nil
		}
		return &ErrorCode{ID: tools.InternalError}, err
	}

	prod := &models.Product{
		ID:    im.ID,
		Image: im.ImagePath,
	}

	if err := am.ProductRepo.UpdateImage(prod); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}

func (am *AdminRestaurantManager) DeleteProduct(ctx context.Context, id *ProtoID) (*ErrorCode, error) {
	_, err := am.ProductRepo.GetProductByID(id.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ErrorCode{ID: tools.DoesntExist}, nil
		}
		return &ErrorCode{ID: tools.InternalError}, err
	}

	if err := am.ProductRepo.Delete(id.ID); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}

func (am *AdminRestaurantManager) CreatePoint(ctx context.Context, pnt *ProtoPoint) (*ErrorCode, error) {
	_, err := am.RestaurantRepo.GetByID(pnt.RestID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ErrorCode{ID: tools.DoesntExist}, nil
		}
		return &ErrorCode{ID: tools.InternalError}, err
	}

	gp, err := am.GeoDataRepo.GetGeoPosByAddress(pnt.Address)

	if err != nil {
		if err == tools.ApiNotHouseAnswerError {
			return &ErrorCode{ID: tools.AddressNotHouse}, nil
		}

		return &ErrorCode{ID: tools.InternalError}, err
	}

	point := &models.RestaurantPoint{
		Address:       pnt.Address,
		MapPoint:      gp,
		ServiceRadius: float64(pnt.Radius),
		RestID:        pnt.RestID,
	}

	if err := am.RestPointsRepo.InsertInto(point); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}

func (am *AdminRestaurantManager) DeletePoint(ctx context.Context, id *ProtoID) (*ErrorCode, error) {
	_, err := am.RestPointsRepo.GetPointByID(id.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ErrorCode{ID: tools.DoesntExist}, nil
		}
		return &ErrorCode{ID: tools.InternalError}, err
	}

	if err := am.RestPointsRepo.Delete(id.ID); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}
