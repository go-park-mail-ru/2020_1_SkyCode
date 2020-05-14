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

func (rtr *RestTagRepository) CreateTagRestRelation(restID, tagID uint64) error {
	if _, err := rtr.db.Exec("INSERT INTO restaurants_and_tags (rest_id, resttag_id) "+
		"VALUES ($1, $2) RETURNING id", restID, tagID); err != nil {
		return err
	}

	return nil
}

func (rtr *RestTagRepository) CheckTagRestRelation(restID, tagID uint64) (bool, error) {
	var id uint64
	if err := rtr.db.QueryRow("SELECT id FROM restaurants_and_tags "+
		"WHERE rest_id = $1 AND resttag_id = $2 LIMIT 1", restID, tagID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (rtr *RestTagRepository) DeleteTagRestRelation(restID, tagID uint64) error {
	if _, err := rtr.db.Exec("DELETE FROM restaurants_and_tags "+
		"WHERE rest_id = $1 AND resttag_id = $2", restID, tagID); err != nil {
		return err
	}

	return nil
}
