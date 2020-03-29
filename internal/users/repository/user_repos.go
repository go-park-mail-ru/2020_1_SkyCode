package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/users"
	"github.com/jackc/pgx"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) users.Repository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) InsertInto(user *models.User) error {
	if err := ur.db.QueryRow("INSERT INTO users(firstName, lastName, email, phone, password, avatar) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.Password,
		user.Avatar).Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Update(user *models.User) error {
	if err := ur.db.QueryRow("UPDATE users SET first_name = $2, last_name = $3, email = $4 WHERE id = $1" +
		" RETURNING firstName, lastName, email",
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email).Scan(&user); err != nil {
			return err
		}

	return nil
}

func (ur *UserRepository) Get(user *models.User) error {
	if err := ur.db.QueryRow("SELECT id, firstName, lastName, email, password, avatar FROM users WHERE phone = $1",
		user.Phone).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Avatar); err != nil {
		return err
	}

	return nil
}
