package repository

import (
	"github.com/jackc/pgx"

	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/products"
)

type ProductRepository struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) products.Repository {
	return &ProductRepository{
		db: db,
	}
}

func (pr *ProductRepository) GetRestaurantProducts(restID uint64) ([]*models.Product, error) {
	productList := []*models.Product{}

	rows, err := pr.db.Query("SELECT id, name, price, image FROM products WHERE rest_id = $1", restID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		product := &models.Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Image)
		if err != nil {
			return nil, err
		}
		productList = append(productList, product)
	}

	return productList, nil
}
