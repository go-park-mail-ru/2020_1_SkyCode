package repository

import (
	"github.com/jackc/pgx"

	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants"
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
		restaraunt := &models.Restaurant{}
		err = rows.Scan(&restaraunt.ID, &restaraunt.Name, &restaraunt.Rating, &restaraunt.Image)
		if err != nil {
			return nil, err
		}
		restaurantsList = append(restaurantsList, restaraunt)
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
