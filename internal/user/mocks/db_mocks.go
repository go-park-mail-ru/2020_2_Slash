package mocks

import (
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func MockUserRepoInsertReturnRows(mock sqlmock.Sqlmock, user *models.User) {
	mock.ExpectBegin()
	insertAnswer := sqlmock.NewRows([]string{"id"}).AddRow(user.ID)
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(user.Nickname, user.Email, user.Password, user.Avatar, user.Role).
		WillReturnRows(insertAnswer)
	mock.ExpectCommit()
}

func MockUserRepoInsertReturnErrNoUniq(mock sqlmock.Sqlmock, user *models.User) {
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(user.Nickname, user.Email, user.Password, user.Avatar, user.Role).
		WillReturnError(errors.New("No UNIQUE"))
	mock.ExpectRollback()
}

func MockUserRepoUpdateReturnResultOk(mock sqlmock.Sqlmock, user *models.User) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE users`).
		WithArgs(user.ID, user.Nickname, user.Email, user.Password, user.Avatar, user.Role).
		WillReturnResult(sqlmock.NewResult(int64(user.ID), 1))
	mock.ExpectCommit()
}

func MockUserRepoSelectByIDReturnRows(mock sqlmock.Sqlmock, user *models.User) {
	rows := sqlmock.NewRows([]string{"id", "nickname", "email", "password",
		"avatar", "role"})
	rows.AddRow(user.ID, user.Nickname, user.Email, user.Password,
		user.Avatar, user.Role)
	mock.ExpectQuery(`SELECT`).WithArgs(user.ID).WillReturnRows(rows)
}

func MockUserRepoSelectByEmailReturnRows(mock sqlmock.Sqlmock, user *models.User) {
	rows := sqlmock.NewRows([]string{"id", "nickname", "email", "password",
		"avatar", "role"})
	rows.AddRow(user.ID, user.Nickname, user.Email, user.Password,
		user.Avatar, user.Role)
	mock.ExpectQuery(`SELECT`).WithArgs(user.Email).WillReturnRows(rows)
}
