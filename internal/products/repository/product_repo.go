package repository

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/products"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) products.Repository {
	return &ProductRepository{
		db: db,
	}
}

func (pr *ProductRepository) GetProductsByRestID(
	restID uint64) ([]*models.Product, error) {
	productList := []*models.Product{}

	rows, err := pr.db.Query("SELECT id, name, price, image, coalesce(tag, 0) "+
		"FROM products WHERE rest_id = $1", restID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		product := &models.Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Image, &product.Tag)
		if err != nil {
			return nil, err
		}
		productList = append(productList, product)
	}

	return productList, nil
}

func (pr *ProductRepository) GetProductByID(prodID uint64) (*models.Product, error) {
	product := &models.Product{}

	if err := pr.db.QueryRow("SELECT id, name, price, image, rest_id, coalesce(tag, 0) "+
		"FROM products WHERE id = $1",
		prodID).Scan(&product.ID, &product.Name, &product.Price,
		&product.Image, &product.RestId, &product.Tag); err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *ProductRepository) InsertInto(product *models.Product) error {
	if err := pr.db.QueryRow("INSERT INTO products (name, price, image, rest_id, tag) "+
		"VALUES ($1, $2, $3, $4, CASE WHEN $5 = 0 THEN NULL ELSE $5 END) RETURNING id",
		product.Name,
		product.Price,
		product.Image,
		product.RestId,
		product.Tag).Scan(&product.ID); err != nil {
		return err
	}

	return nil
}

func (pr *ProductRepository) Delete(prodID uint64) error {
	if _, err := pr.db.Exec("DELETE FROM products WHERE id = $1", prodID); err != nil {
		return err
	}

	return nil
}

func (pr *ProductRepository) Update(product *models.Product) error {
	if _, err := pr.db.Exec("UPDATE products SET name = $2, price = $3, "+
		"tag = (CASE WHEN $4 = 0 THEN NULL ELSE $4 END) WHERE id = $1",
		product.ID,
		product.Name,
		product.Price,
		product.Tag); err != nil {
		return err
	}

	return nil
}

func (pr *ProductRepository) UpdateImage(product *models.Product) error {
	if _, err := pr.db.Exec("UPDATE products SET image = $2 WHERE id = $1",
		product.ID,
		product.Image); err != nil {
		return err
	}

	return nil
}
