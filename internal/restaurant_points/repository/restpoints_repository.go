package repository

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurant_points"
)

type RestPointsRepository struct {
	db *sql.DB
}

func NewRestPosintsRepository(db *sql.DB) restaurant_points.Repository {
	return &RestPointsRepository{
		db: db,
	}
}

func (rr *RestPointsRepository) InsertInto(p *models.RestaurantPoint) error {
	if err := rr.db.QueryRow("INSERT INTO rest_points (address, latitude, longitude, restid, radius) "+
		"VALUES ($1, $2, $3, $4, $5) RETURNING id", p.Address, p.MapPoint.Latitude, p.MapPoint.Longitude,
		p.RestID, p.ServiceRadius).Scan(&p.ID); err != nil {
		return err
	}

	return nil
}

func (rr *RestPointsRepository) Delete(id uint64) error {
	if _, err := rr.db.Exec("DELETE FROM rest_points WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}

func (rr *RestPointsRepository) GetAll() ([]*models.RestaurantPoint, error) {
	rows, err := rr.db.Query("SELECT id, address, latitude, longitude, restid, radius " +
		"FROM rest_points")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	returnPoints := []*models.RestaurantPoint{}
	for rows.Next() {
		rp := &models.RestaurantPoint{MapPoint: &models.GeoPos{}}

		if err = rows.Scan(&rp.ID, &rp.Address, &rp.MapPoint.Latitude,
			&rp.MapPoint.Longitude, &rp.RestID, &rp.ServiceRadius); err != nil {
			return nil, err
		}

		returnPoints = append(returnPoints, rp)
	}

	return returnPoints, nil
}

func (rr *RestPointsRepository) GetPointsByRestID(restID,
	count, page uint64) ([]*models.RestaurantPoint, uint64, error) {
	rows, err := rr.db.Query("SELECT id, address, latitude, longitude, restid, radius FROM rest_points "+
		"WHERE restid = $1 ORDER BY address LIMIT $2 OFFSET $3", restID, count, count*(page-1))
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	returnPoints := []*models.RestaurantPoint{}

	for rows.Next() {
		p := &models.RestaurantPoint{MapPoint: &models.GeoPos{}}

		err = rows.Scan(&p.ID, &p.Address, &p.MapPoint.Latitude, &p.MapPoint.Longitude,
			&p.RestID, &p.ServiceRadius)
		if err != nil {
			return nil, 0, err
		}

		returnPoints = append(returnPoints, p)
	}

	total := uint64(0)
	if err := rr.db.QueryRow("SELECT count(*) FROM rest_points WHERE restid = $1", restID).
		Scan(&total); err != nil {
		return nil, 0, err
	}

	return returnPoints, total, nil
}

func (rr *RestPointsRepository) GetCloserPointByRestID(restID uint64,
	pos *models.GeoPos) (*models.RestaurantPoint, error) {
	returnPoint := &models.RestaurantPoint{MapPoint: &models.GeoPos{}}
	if err := rr.db.QueryRow("SELECT id, restid, latitude, longitude, address, radius "+
		"FROM rest_points "+
		"WHERE restid = $1 AND ST_Distance(ST_MakePoint(latitude, longitude)::geography, "+
		"ST_Makepoint($2, $3)::geography) <= radius * 1000 "+
		"ORDER BY ST_Distance(ST_MakePoint(latitude, longitude)::geography, "+
		"ST_MakePoint($2, $3)::geography) ASC LIMIT 1",
		restID, pos.Latitude, pos.Longitude).Scan(&returnPoint.ID, &returnPoint.RestID,
		&returnPoint.MapPoint.Latitude, &returnPoint.MapPoint.Longitude,
		&returnPoint.Address, &returnPoint.ServiceRadius); err != nil {
		return nil, err
	}

	return returnPoint, nil
}

func (rr *RestPointsRepository) GetPointByID(id uint64) (*models.RestaurantPoint, error) {
	returnPoint := &models.RestaurantPoint{MapPoint: &models.GeoPos{}}
	if err := rr.db.QueryRow("SELECT id, address, latitude, longitude, restid, radius FROM rest_points "+
		"WHERE id = $1", id).Scan(&returnPoint.ID, &returnPoint.Address, &returnPoint.MapPoint.Latitude,
		&returnPoint.MapPoint.Longitude, &returnPoint.RestID, &returnPoint.ServiceRadius); err != nil {
		return nil, err
	}

	return returnPoint, nil
}
