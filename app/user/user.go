package user

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Avatar   string `json:"avatar"`
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(newDb *sql.DB) *UserRepo {
	return &UserRepo{
		db: newDb,
	}
}

func (ur *UserRepo) Get(userID uint64) (*User, bool) {
	user := &User{}
	row := ur.db.QueryRow(
		`SELECT id, nickname, email, password, avatar
		FROM profile
		WHERE id=$1`, userID)
	has := row.Scan(&user.ID, &user.Nickname,
		&user.Email, &user.Password, &user.Avatar)
	if has != nil {
		return nil, false
	}
	return user, true
}

func (ur *UserRepo) Delete(userID uint64) error {
	_, has := ur.Get(userID)
	if !has {
		return errors.New("There is no such user")
	}

	tx, err := ur.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = ur.db.Exec(
		`DELETE FROM profile
		WHERE id=$1`, userID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) UpdateEmail(userID uint64, email string) error {
	user, has := ur.Get(userID)
	if !has {
		return errors.New("There is no such user")
	}
	if user.Email == email {
		return nil
	}
	if !ur.IsUniqEmail(email) {
		return errors.New("User with this Email already exists")
	}

	tx, err := ur.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	if _, err = ur.db.Exec(
		`UPDATE profile
		SET email=$1
		WHERE id=$2`, email, userID); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) Exists(userID uint64) bool {
	_, has := ur.Get(userID)
	return has
}

func (ur *UserRepo) IsUniqEmail(email string) bool {
	_, has := ur.GetByEmail(email)
	return !has
}

func (ur *UserRepo) Register(user *User) error {
	if !ur.IsUniqEmail(user.Email) {
		return errors.New("User with this Email already exists")
	}

	tx, err := ur.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	row := ur.db.QueryRow(
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

func (ur *UserRepo) GetByEmail(userEmail string) (*User, bool) {
	user := &User{}
	row := ur.db.QueryRow(
		`SELECT id, nickname, email, password, avatar
		FROM profile
		WHERE email=$1`, userEmail)
	has := row.Scan(&user.ID, &user.Nickname,
		&user.Email, &user.Password, &user.Avatar)
	if has != nil {
		return nil, false
	}
	return user, true
}
