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
	if err := rr.db.QueryRow("INSERT INTO rest_points (address, latitude, longitude, restid) "+
		"VALUES ($1, $2, $3, $4) RETURNING id", p.Address, p.MapPoint.Latitude, p.MapPoint.Longitude, p.RestID).
		Scan(&p.ID); err != nil {
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

func (rr *RestPointsRepository) GetPointsByRestID(restID, count, page uint64) ([]*models.RestaurantPoint, uint64, error) {
	rows, err := rr.db.Query("SELECT id, address, latitude, longitude, restid FROM rest_points "+
		"WHERE restid = $1 ORDER BY address LIMIT $2 OFFSET $3", restID, count, count*(page-1))
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	returnPoints := []*models.RestaurantPoint{}

	for rows.Next() {
		p := &models.RestaurantPoint{MapPoint: &models.GeoPos{}}

		err = rows.Scan(&p.ID, &p.Address, &p.MapPoint.Latitude, &p.MapPoint.Longitude, &p.RestID)
		if err != nil {
			return nil, 0, err
		}

		returnPoints = append(returnPoints, p)
	}

	total := uint64(0)
	if err := rr.db.QueryRow("SELECT count(*) FROM rest_points WHERE restid = $1").
		Scan(&total); err != nil {
		return nil, 0, err
	}

	return returnPoints, total, nil
}

func (rr *RestPointsRepository) GetPointByID(id uint64) (*models.RestaurantPoint, error) {
	returnPoint := &models.RestaurantPoint{MapPoint: &models.GeoPos{}}
	if err := rr.db.QueryRow("SELECT id, address, latitude, longitude, restid FROM rest_points "+
		"WHERE id = $1", id).Scan(&returnPoint.ID, &returnPoint.Address, &returnPoint.MapPoint.Latitude,
		&returnPoint.MapPoint.Longitude, &returnPoint.RestID); err != nil {
		return nil, err
	}

	return returnPoint, nil
}
