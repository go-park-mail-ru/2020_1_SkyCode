package repository

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/reviews"
	"math"
)

type ReviewsRepository struct {
	db *sql.DB
}

func NewReviewsRepository(db *sql.DB) reviews.Repository {
	return &ReviewsRepository{
		db: db,
	}
}

func (rr *ReviewsRepository) GetRatingByRestID(restID uint64) (float64, error) {
	rating := float64(0)

	if err := rr.db.QueryRow("SELECT avg(rating) FROM reviews WHERE rest_id = $1", restID).
		Scan(&rating); err != nil {
		return 0, err
	}

	return math.Round(rating*100) / 100, nil
}

func (rr *ReviewsRepository) GetReviewsByRestID(restID, count, page uint64) ([]*models.Review, error) {
	rows, err := rr.db.Query(
		"SELECT r.id, r.restId, u.id, u.firstname, u.lastname, r.message, r.creationDate, r.rate "+
			"FROM reviews AS r "+
			"JOIN users AS u ON (r.userid = u.id) "+
			"WHERE r.restID = $1 AND r.message <> '' ORDER BY r.rate DESC, r.creationDate DESC LIMIT $2 OFFSET $3",
		restID, count, (page-1)*count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	returnReviews := []*models.Review{}
	for rows.Next() {
		r := &models.Review{Author: &models.User{}}
		if err := rows.Scan(&r.ID, &r.RestID, &r.Author.ID, &r.Author.FirstName,
			&r.Author.LastName, &r.Text, &r.CreationDate, &r.Rate); err != nil {
			return nil, err
		}

		returnReviews = append(returnReviews, r)
	}

	return returnReviews, nil
}

func (rr *ReviewsRepository) GetReviewsCountByRestID(restID uint64) (uint64, error) {
	var count uint64
	if err := rr.db.QueryRow("SELECT COUNT(*) FROM reviews WHERE restId = $1", restID).
		Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (rr *ReviewsRepository) GetReviewsByUserID(userID, count, page uint64) ([]*models.Review, error) {
	rows, err := rr.db.Query(
		"SELECT r.id, r.restId, u.id, u.firstname, u.lastname, r.message, r.creationDate, r.rate "+
			"FROM reviews AS r "+
			"JOIN users AS u ON (r.userid = u.id) "+
			"WHERE r.userId = $1 ORDER BY r.creationDate DESC LIMIT $2 OFFSET $3", userID, count, (page-1)*count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	returnReviews := []*models.Review{}
	for rows.Next() {
		r := &models.Review{Author: &models.User{}}

		if err := rows.Scan(&r.ID, &r.RestID, &r.Author.ID, &r.Author.FirstName,
			&r.Author.LastName, &r.Text, &r.CreationDate, &r.Rate); err != nil {
			return nil, err
		}

		returnReviews = append(returnReviews, r)
	}

	return returnReviews, nil
}

func (rr *ReviewsRepository) GetReviewsCountByUserID(userID uint64) (uint64, error) {
	var count uint64
	if err := rr.db.QueryRow("SELECT COUNT(*) FROM reviews WHERE userId = $1", userID).
		Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (rr *ReviewsRepository) GetRestaurantReviewByUser(restID, userID uint64) (*models.Review, error) {
	r := &models.Review{Author: &models.User{}}
	if err := rr.db.QueryRow("SELECT r.id, r.restId, u.id, u.firstname, u.lastname, r.message, r.creationDate, r.rate "+
		"FROM reviews AS r "+
		"JOIN users AS u ON (r.userid = u.id) "+
		"WHERE r.restId = $1 AND r.userId = $2",
		restID, userID).Scan(&r.ID, &r.RestID, &r.Author.ID, &r.Author.FirstName,
		&r.Author.LastName, &r.Text, &r.CreationDate, &r.Rate); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r, nil
}

func (rr *ReviewsRepository) GetReviewByID(id uint64) (*models.Review, error) {
	r := &models.Review{Author: &models.User{}}

	if err := rr.db.QueryRow(
		"SELECT r.id, r.restId, u.id, u.firstname, u.lastname, r.message, r.creationDate, r.rate "+
			"FROM reviews AS r "+
			"JOIN users AS u ON (r.userid = u.id) "+
			"WHERE r.id = $1", id).Scan(&r.ID, &r.RestID, &r.Author.ID, &r.Author.FirstName,
		&r.Author.LastName, &r.Text, &r.CreationDate, &r.Rate); err != nil {
		return nil, err
	}

	return r, nil
}

func (rr *ReviewsRepository) CreateReview(r *models.Review) error {
	if err := rr.db.QueryRow("INSERT INTO reviews (restId, userId, message, creationDate, rate) "+
		"VALUES ($1, $2, $3, $4, $5) RETURNING id", r.RestID, r.Author.ID, r.Text, r.CreationDate, r.Rate).
		Scan(&r.ID); err != nil {
		return err
	}

	return nil
}

func (rr *ReviewsRepository) UpdateReview(r *models.Review) error {
	if _, err := rr.db.Exec("UPDATE reviews SET message = $2, rate = $3 WHERE id = $1",
		r.ID, r.Text, r.Rate); err != nil {
		return err
	}

	return nil
}

func (rr *ReviewsRepository) DeleteReview(id uint64) error {
	if _, err := rr.db.Exec("DELETE FROM reviews "+
		"WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
