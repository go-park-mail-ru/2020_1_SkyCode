package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/jackc/pgx"
)

type Repository struct {
	db *pgx.Conn
}

func NewSessionRepository(db *pgx.Conn) *Repository {
	return &Repository{
		db:db,
	}
}

func (sR *Repository) InsertInto(session *models.Session) error {
	if err := sR.db.QueryRow("INSERT INTO sessions(userId, token) VALUES ($1, $2) RETURNING id",
		session.UserId,
		session.Token).Scan(&session.ID); err != nil {
		return err
	}

	return nil
}

func (sR *Repository) Delete(session *models.Session) error {
	if err := sR.db.QueryRow("DELETE FROM sessions WHERE id = $1 RETURNING token, userId",
		session.ID).Scan(&session.Token, &session.UserId); err != nil {
		return err
	}

	return nil
}

func (sR *Repository) Get(session *models.Session) error {
	if err := sR.db.QueryRow("SELECT id, userId, token FROM sessions WHERE token = $1",
		session.Token).Scan(&session.ID, &session.UserId, &session.Token); err != nil {
		return err
	}

	return nil
}
