package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
)

type UserPgRepository struct {
	dbConn *sql.DB
}

func NewUserPgRepository(conn *sql.DB) user.UserRepository {
	return &UserPgRepository{
		dbConn: conn,
	}
}

func (ur *UserPgRepository) Insert(user *models.User) error {
	tx, err := ur.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	row := ur.dbConn.QueryRow(
		`INSERT INTO profile(nickname, email, password, avatar)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		user.Nickname, user.Email, user.Password, user.Avatar)

	err = row.Scan(&user.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserPgRepository) SelectByEmail(email string) (*models.User, error) {
	user := &models.User{}

	row := ur.dbConn.QueryRow(
		`SELECT id, nickname, email, password, avatar
		FROM profile
		WHERE email=$1`, email)

	err := row.Scan(&user.ID, &user.Nickname, &user.Email, &user.Password, &user.Avatar)
	if err != nil {
		return nil, err
	}
	return user, nil
}
