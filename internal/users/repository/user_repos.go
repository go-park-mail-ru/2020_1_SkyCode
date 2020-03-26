package repository

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/users"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) users.Repository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) InsertInto(user *models.User) error {
	if err := ur.db.QueryRow("INSERT INTO users(first_name, last_name, email, phone, password, avatar) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.FirstName, user.LastName, user.Email, user.Phone, user.Password, user.ProfilePhoto).Scan(&user.ID); err != nil {
		return err
	}

	return nil
}
