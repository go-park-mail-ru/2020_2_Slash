package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
)

type SessionPgRepository struct {
	dbConn *sql.DB
}

func NewSessionPgRepository(conn *sql.DB) session.SessionRepository {
	return &SessionPgRepository{
		dbConn: conn,
	}
}

func (sr *SessionPgRepository) Insert(session *models.Session) error {
	tx, err := sr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = tx.QueryRow(
		`INSERT INTO session(value, expires, profile_id)
		VALUES ($1, $2, $3) RETURNING id`,
		session.Value, session.ExpiresAt, session.UserID).Scan(&session.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (sr *SessionPgRepository) SelectByID(sessValue string) (*models.Session, error) {
	sess := &models.Session{}

	row := sr.dbConn.QueryRow(
		`SELECT id, value, expires, profile_id
		FROM session WHERE value=$1`, sessValue)

	err := row.Scan(&sess.ID, &sess.Value, &sess.ExpiresAt, &sess.UserID)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
