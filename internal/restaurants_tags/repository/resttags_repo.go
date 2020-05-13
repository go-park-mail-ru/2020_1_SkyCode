package repository

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants_tags"
)

type RestTagRepository struct {
	db *sql.DB
}

func NewRestTagRepository(db *sql.DB) restaurants_tags.Repository {
	return &RestTagRepository{
		db: db,
	}
}

func (rtr *RestTagRepository) InsertInto(tag *models.RestTag) error {
	if err := rtr.db.QueryRow("INSERT INTO rest_tags (name, image) "+
		"VALUES ($1, $2) RETURNING id", tag.Name, tag.Image).Scan(&tag.ID); err != nil {
		return err
	}

	return nil
}

func (rtr *RestTagRepository) GetAll() ([]*models.RestTag, error) {
	rows, err := rtr.db.Query("SELECT id, name, image FROM rest_tags")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []*models.RestTag
	for rows.Next() {
		tag := &models.RestTag{}
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Image); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (rtr *RestTagRepository) GetByID(id uint64) (*models.RestTag, error) {
	tag := &models.RestTag{}
	if err := rtr.db.QueryRow("SELECT id, name, image FROM rest_tags WHERE id = $1", id).
		Scan(&tag.ID, &tag.Name, &tag.Image); err != nil {
		return nil, err
	}

	return tag, nil
}

func (rtr *RestTagRepository) Update(tag *models.RestTag) error {
	if _, err := rtr.db.Exec("UPDATE rest_tags SET name = $2, image = $3 WHERE id = $1",
		tag.ID, tag.Name, tag.Image); err != nil {
		return err
	}

	return nil
}

func (rtr *RestTagRepository) Delete(id uint64) error {
	if _, err := rtr.db.Exec("DELETE FROM rest_tags WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
