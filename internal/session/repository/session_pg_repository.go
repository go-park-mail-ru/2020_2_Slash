package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
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
		`INSERT INTO sessions(value, expires, user_id)
		VALUES ($1, $2, $3) RETURNING id`,
		session.Value, session.ExpiresAt, session.UserID).Scan(&session.ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (sr *SessionPgRepository) SelectByValue(sessValue string) (*models.Session, error) {
	sess := &models.Session{}

	row := sr.dbConn.QueryRow(
		`SELECT id, value, expires, user_id
		FROM sessions WHERE value=$1`, sessValue)

	err := row.Scan(&sess.ID, &sess.Value, &sess.ExpiresAt, &sess.UserID)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (sr *SessionPgRepository) DeleteByValue(sessionValue string) error {
	tx, err := sr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM sessions
		WHERE value=$1`, sessionValue)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
