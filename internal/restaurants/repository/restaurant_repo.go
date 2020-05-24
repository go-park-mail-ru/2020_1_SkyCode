package repository

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants"
	"github.com/sirupsen/logrus"
)

type RestaurantRepository struct {
	db *sql.DB
}

func NewRestaurantRepository(db *sql.DB) restaurants.Repository {
	return &RestaurantRepository{
		db: db,
	}
}

func (rr *RestaurantRepository) GetAll(count uint64, page uint64, tagID uint64) ([]*models.Restaurant, uint64, error) {
	var restaurantsList []*models.Restaurant
	var rows *sql.Rows
	var err error
	if tagID == 0 {
		rows, err = rr.db.Query("SELECT id, name, rating, image FROM restaurants "+
			"LIMIT $1 OFFSET $2", count, (page-1)*count)
	} else {
		rows, err = rr.db.Query("SELECT r.id, r.name, r.rating, r.image FROM restaurants r "+
			"JOIN restaurants_and_tags rt ON (r.id = rt.rest_id) WHERE rt.resttag_id = $3 "+
			"LIMIT $1 OFFSET $2", count, (page-1)*count, tagID)
	}
	if err != nil {
		return nil, 0, err
	}

	var total uint64
	if tagID == 0 {
		err = rr.db.QueryRow("SELECT COUNT(*) FROM restaurants").Scan(&total)
	} else {
		err = rr.db.QueryRow("SELECT COUNT(*) FROM restaurants r "+
			"JOIN restaurants_and_tags rt ON (r.id = rt.rest_id) WHERE rt.resttag_id = $1", tagID).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	for rows.Next() {
		restaurant := &models.Restaurant{}
		err = rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Rating, &restaurant.Image)
		if err != nil {
			return nil, 0, err
		}
		restaurantsList = append(restaurantsList, restaurant)
	}

	return restaurantsList, total, nil
}

func (rr *RestaurantRepository) GetRecommendationsByOrder(userID uint64, count uint64) ([]*models.Restaurant, error) {
	rows, err := rr.db.Query("WITH ordered_rests AS (SELECT DISTINCT restid FROM orders WHERE userid = $1), "+
		"ordered_tags AS (SELECT DISTINCT resttag_id FROM restaurants_and_tags "+
		"WHERE rest_id IN (SELECT * FROM ordered_rests)) "+
		"SELECT id, name, description, rating, image FROM restaurants "+
		"WHERE id IN (SELECT r.id FROM restaurants r "+
		"JOIN restaurants_and_tags rat ON r.id = rat.rest_id WHERE rat.resttag_id IN (SELECT * FROM ordered_tags) "+
		"AND r.id NOT IN (SELECT * FROM ordered_rests) "+
		"GROUP BY r.id) "+

		"UNION "+

		"SELECT id, name, description, rating, image FROM restaurants "+
		"WHERE id IN (SELECT r.id FROM restaurants r "+
		"WHERE r.id NOT IN (SELECT * FROM ordered_rests) "+
		"GROUP BY r.id) ORDER BY rating DESC LIMIT $2", userID, count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rest []*models.Restaurant

	for rows.Next() {
		r := &models.Restaurant{}

		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.Rating, &r.Image); err != nil {
			return nil, err
		}

		rest = append(rest, r)
	}

	return rest, nil
}

func (rr *RestaurantRepository) GetByID(id uint64) (*models.Restaurant, error) {
	restaurant := &models.Restaurant{}

	err := rr.db.
		QueryRow("SELECT id, moderId, name, description, rating, image FROM restaurants WHERE id = $1", id).
		Scan(&restaurant.ID, &restaurant.ManagerID, &restaurant.Name, &restaurant.Description,
			&restaurant.Rating, &restaurant.Image)

	if err != nil {
		return nil, err
	}
	return restaurant, nil
}

func (rr *RestaurantRepository) GetAllInServiceRadius(
	pos *models.GeoPos, count, page, tagID uint64) ([]*models.Restaurant, uint64, error) {
	var rows *sql.Rows
	var err error
	if tagID == 0 {
		rows, err = rr.db.Query("SELECT r.id, r.name, r.description, r.rating, r.image, "+
			"min(st_distance("+
			"st_makepoint(rp.latitude, rp.longitude)::geography, "+
			"st_makepoint($1, $2)::geography)) as dst "+
			"FROM restaurants r "+
			"JOIN rest_points rp ON (r.id = rp.restid) "+
			"WHERE ST_DWithin("+
			"ST_MakePoint(rp.latitude, rp.longitude)::geography, "+
			"ST_MakePoint($1, $2)::geography, rp.radius * 1000) "+
			"GROUP BY r.id, r.rating "+
			"ORDER BY dst ASC, r.rating DESC "+
			"LIMIT $3 OFFSET $4", pos.Latitude, pos.Longitude, count, count*(page-1))
	} else {
		rows, err = rr.db.Query("SELECT r.id, r.name, r.description, r.rating, r.image, "+
			"min(st_distance("+
			"st_makepoint(rp.latitude, rp.longitude)::geography, "+
			"st_makepoint($1, $2)::geography)) as dst "+
			"FROM restaurants r "+
			"JOIN rest_points rp ON (r.id = rp.restid) "+
			"JOIN restaurants_and_tags rt ON (r.id = rt.rest_id) "+
			"WHERE ST_DWithin("+
			"ST_MakePoint(rp.latitude, rp.longitude)::geography, "+
			"ST_MakePoint($1, $2)::geography, rp.radius * 1000) AND "+
			"rt.resttag_id = $5 "+
			"GROUP BY r.id, r.rating "+
			"ORDER BY dst ASC, r.rating DESC "+
			"LIMIT $3 OFFSET $4", pos.Latitude, pos.Longitude, count, count*(page-1), tagID)
	}

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	returnRests := []*models.Restaurant{}
	for rows.Next() {
		rest := &models.Restaurant{}
		var dist float64
		if err := rows.Scan(&rest.ID, &rest.Name,
			&rest.Description, &rest.Rating, &rest.Image, &dist); err != nil {
			return nil, 0, err
		}

		returnRests = append(returnRests, rest)
	}

	total := uint64(0)
	if tagID == 0 {
		if err := rr.db.QueryRow("SELECT COUNT(r.id) "+
			"FROM restaurants r "+
			"JOIN rest_points rp ON (r.id = rp.restid) "+
			"WHERE ST_DWithin("+
			"ST_MakePoint(rp.latitude, rp.longitude)::geography, "+
			"ST_MakePoint($1, $2)::geography, rp.radius * 1000) "+
			"GROUP BY r.id", pos.Latitude, pos.Longitude).Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := rr.db.QueryRow("SELECT COUNT(r.id) "+
			"FROM restaurants r "+
			"JOIN rest_points rp ON (r.id = rp.restid) "+
			"JOIN restaurants_and_tags rt ON (r.id = rt.rest_id) "+
			"WHERE ST_DWithin("+
			"ST_MakePoint(rp.latitude, rp.longitude)::geography, "+
			"ST_MakePoint($1, $2)::geography, rp.radius * 1000) AND "+
			"rt.resttag_id = $3 "+
			"GROUP BY r.id", pos.Latitude, pos.Longitude, tagID).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	return returnRests, total, nil
}

func (rr *RestaurantRepository) GetRecommendationsInRadius(pos *models.GeoPos,
	userID uint64, count uint64) ([]*models.Restaurant, error) {
	rows, err := rr.db.Query("WITH ordered_rests AS (SELECT DISTINCT restid FROM orders WHERE userid = $1), "+
		"ordered_tags AS (SELECT DISTINCT resttag_id FROM restaurants_and_tags "+
		"WHERE rest_id IN (SELECT * FROM ordered_rests)) "+
		"SELECT r.id, r.name, r.description, r.rating, r.image, "+
		"min(st_distance("+
		"st_makepoint(rp.latitude, rp.longitude)::geography, "+
		"st_makepoint($2, $3)::geography)) as dst "+
		"FROM restaurants r "+
		"JOIN rest_points rp ON (r.id = rp.restid) "+
		"JOIN restaurants_and_tags rat ON r.id = rat.rest_id WHERE rat.resttag_id IN (SELECT * FROM ordered_tags) "+
		"AND ST_DWithin("+
		"ST_MakePoint(rp.latitude, rp.longitude)::geography, "+
		"ST_MakePoint($2, $3)::geography, rp.radius * 1000) AND "+
		"r.id NOT IN (SELECT * FROM ordered_rests) "+
		"GROUP BY r.id, r.rating "+

		"UNION "+

		"SELECT r.id, r.name, r.description, r.rating, r.image, "+
		"min(st_distance("+
		"st_makepoint(rp.latitude, rp.longitude)::geography, "+
		"st_makepoint($2, $3)::geography)) as dst "+
		"FROM restaurants r "+
		"JOIN rest_points rp ON (r.id = rp.restid) "+
		"WHERE ST_DWithin("+
		"ST_MakePoint(rp.latitude, rp.longitude)::geography, "+
		"ST_MakePoint($2, $3)::geography, rp.radius * 1000) AND "+
		"r.id NOT IN (SELECT * FROM ordered_rests) "+
		"GROUP BY r.id, r.rating "+
		"ORDER BY rating DESC "+
		"LIMIT $4", userID, pos.Latitude, pos.Longitude, count)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rest []*models.Restaurant

	for rows.Next() {
		r := &models.Restaurant{}
		var dist float64

		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.Rating, &r.Image, &dist); err != nil {
			return nil, err
		}

		rest = append(rest, r)
	}

	return rest, nil
}

func (rr *RestaurantRepository) InsertInto(rest *models.Restaurant) error {
	if err := rr.db.QueryRow("INSERT INTO restaurants (moderId, name, description, rating, image) "+
		"VALUES ($1, $2, $3, $4, $5) RETURNING id",
		rest.ManagerID,
		rest.Name,
		rest.Description,
		rest.Rating,
		rest.Image).Scan(&rest.ID); err != nil {
		return err
	}
	logrus.Info(rest.ID)

	return nil
}

func (rr *RestaurantRepository) Update(rest *models.Restaurant) error {
	if _, err := rr.db.Exec("UPDATE restaurants SET name = $2, description = $3 "+
		"WHERE id = $1", rest.ID, rest.Name, rest.Description); err != nil {
		return err
	}

	return nil
}

func (rr *RestaurantRepository) UpdateImage(rest *models.Restaurant) error {
	if _, err := rr.db.Exec("UPDATE restaurants SET image = $2 "+
		"WHERE id = $1", rest.ID, rest.Image); err != nil {
		return err
	}

	return nil
}

func (rr *RestaurantRepository) Delete(restID uint64) error {
	if _, err := rr.db.Exec("DELETE FROM restaurants WHERE id = $1", restID); err != nil {
		return err
	}

	return nil
}
