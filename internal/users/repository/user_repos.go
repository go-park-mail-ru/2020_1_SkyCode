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
	if err := ur.db.QueryRow("UPDATE users SET firstName = $2, lastName = $3, email = $4 WHERE id = $1" +
		" RETURNING firstName, lastName, email",
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email).Scan(&user.FirstName, &user.LastName, &user.Email); err != nil {
			return err
		}

	return nil
}

func (ur *UserRepository) UpdatePhone(user *models.User) error {
	if err := ur.db.QueryRow("UPDATE users SET Phone = $2 WHERE id = $1" +
		" RETURNING firstName, lastName, email, avatar",
		user.ID,
		user.Phone).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Avatar); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateAvatar(user *models.User) error {
	if err := ur.db.QueryRow("UPDATE users SET Avatar = $2 WHERE id = $1" +
		" RETURNING firstName, lastName, email, phone",
		user.ID,
		user.Avatar).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Phone); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdatePassword(user *models.User) error {
	if err := ur.db.QueryRow("UPDATE users SET Password = $2 WHERE id = $1" +
		" RETURNING firstName, lastName, email, avatar",
		user.ID,
		user.Password).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Avatar); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetById(user *models.User) error {
	if err := ur.db.QueryRow("SELECT firstName, lastName, email, password, avatar, phone FROM users WHERE id = $1",
		user.ID).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Avatar, &user.Phone); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetByPhone(user *models.User) error {
	if err := ur.db.QueryRow("SELECT id, firstName, lastName, email, password, avatar FROM users WHERE phone = $1",
		user.Phone).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Avatar); err != nil {
		return err
	}

	return nil
}
