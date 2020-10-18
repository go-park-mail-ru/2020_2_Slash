package session

import (
	"context"
	"database/sql"
	"errors"
	uuid "github.com/satori/go.uuid"
	"time"

	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
)

var (
	ErrNoAuth = errors.New("No session found")
)

type Session struct {
	Id        uint64
	Value     string
	UserID    uint64
	ExpiresAt time.Time
}

func NewSession(userID uint64) *Session {
	value := uuid.NewV4().String()
	return &Session{
		Value:     value,
		UserID:    userID,
		ExpiresAt: time.Now().Add(10 * time.Hour),
	}
}

type SessionManager struct {
	db *sql.DB
}

func NewSessionManager(newDb *sql.DB) *SessionManager {
	return &SessionManager{
		db: newDb,
	}
}

func (sm *SessionManager) Create(user *user.User) (*Session, error) {
	session := NewSession(user.ID)

	tx, err := sm.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(
		`INSERT INTO session(value, expires, profile_id)
		VALUES ($1, $2, $3) RETURNING id`,
		session.Value, session.ExpiresAt, session.UserID).
		Scan(&session.Id)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return session, nil
}

func (sm *SessionManager) Get(sessValue string) (*Session, bool) {
	session := &Session{}
	row := sm.db.QueryRow(
		`SELECT id, value, expires, profile_id
		FROM session WHERE value=$1`, sessValue)
	has := row.Scan(&session.Id, &session.Value,
		&session.ExpiresAt, &session.UserID)
	if has != nil {
		return nil, false
	}
	return session, true
}

func (sm *SessionManager) IsValid(session *Session) bool {
	hasExpired := session.ExpiresAt.Before(time.Now())
	if hasExpired {
		return false
	}
	return true
}

func (sm *SessionManager) GetUserSession(user *user.User) *Session {
	session := &Session{}
	row := sm.db.QueryRow(
		`SELECT id, value, expires, profile_id
		FROM session WHERE profile_id=$1`, user.ID)
	has := row.Scan(&session.Id, &session.Value,
		&session.ExpiresAt, &session.UserID)
	if has != nil {
		return nil
	}
	return session
}

func (sm *SessionManager) IsAuthorized(user *user.User) (bool, error) {
	session := sm.GetUserSession(user)
	if session == nil {
		return false, nil
	}
	if session.ExpiresAt.Before(time.Now()) {
		err := sm.Delete(session.Value)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func (sm *SessionManager) Delete(cookieValue string) error {
	_, has := sm.Get(cookieValue)
	if !has {
		return errors.New("There is no such user")
	}

	tx, err := sm.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = sm.db.Exec(
		`DELETE FROM session
		WHERE value=$1`, cookieValue)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
