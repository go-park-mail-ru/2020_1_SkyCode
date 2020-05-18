package repository

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/product_tags"
)

type ProductTagsRepository struct {
	db *sql.DB
}

func NewProductTagsRepository(db *sql.DB) product_tags.Repository {
	return &ProductTagsRepository{
		db: db,
	}
}

func (ptr *ProductTagsRepository) InsertInto(tag *models.ProductTag) error {
	if err := ptr.db.QueryRow("INSERT INTO product_tags (name, rest_id) "+
		"VALUES ($1, $2) RETURNING id", tag.Name, tag.RestID).Scan(&tag.ID); err != nil {
		return err
	}

	return nil
}

func (ptr *ProductTagsRepository) GetByID(ID uint64) (*models.ProductTag, error) {
	tag := &models.ProductTag{}
	if err := ptr.db.QueryRow("SELECT id, name, rest_id FROM product_tags "+
		"WHERE id = $1", ID).Scan(&tag.ID, &tag.Name, &tag.RestID); err != nil {
		return nil, err
	}

	return tag, nil
}

func (ptr *ProductTagsRepository) GetByRestID(restID uint64) ([]*models.ProductTag, error) {
	rows, err := ptr.db.Query("SELECT id, name, rest_id FROM product_tags "+
		"WHERE rest_id = $1", restID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []*models.ProductTag
	for rows.Next() {
		tag := &models.ProductTag{}

		if err := rows.Scan(&tag.ID, &tag.Name, &tag.RestID); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (ptr *ProductTagsRepository) Delete(ID uint64) error {
	if _, err := ptr.db.Exec("DELETE FROM product_tags WHERE id = $1", ID); err != nil {
		return err
	}

	return nil
}
