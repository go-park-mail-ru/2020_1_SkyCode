package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants"
	"github.com/jackc/pgx"
)

type RestaurantRepository struct {
	db *pgx.Conn
}

func NewRestaurantRepository(db *pgx.Conn) restaurants.Repository {
	return &RestaurantRepository{
		db: db,
	}
}

func (rr *RestaurantRepository) GetAll() ([]*models.Restaurant, error) {
	restaurantsList := []*models.Restaurant{}

	rows, err := rr.db.Query("SELECT id, name, rating, image FROM restaurants")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		restaurant := &models.Restaurant{}
		err = rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Rating, &restaurant.Image)
		if err != nil {
			return nil, err
		}
		restaurantsList = append(restaurantsList, restaurant)
	}

	return restaurantsList, nil
}

func (rr *RestaurantRepository) GetByID(id uint64) (*models.Restaurant, error) {
	restaurant := &models.Restaurant{}

	err := rr.db.
		QueryRow("SELECT id, name, description, rating, image FROM restaurants WHERE id = $1", id).
		Scan(&restaurant.ID, &restaurant.Name, &restaurant.Description,
			&restaurant.Rating, &restaurant.Image)

	if err != nil {
		return nil, err
	}
	return restaurant, nil
}

func (rr *RestaurantRepository) InsertInto(rest *models.Restaurant) error {
	if err := rr.db.QueryRow("INSERT INTO restaurants (name, description, rating, image) "+
		"VALUES ($1, $2, $3, $4) RETURNING id",
		rest.Name,
		rest.Description,
		rest.Rating,
		rest.Image).Scan(&rest.ID); err != nil {
		return err
	}

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
