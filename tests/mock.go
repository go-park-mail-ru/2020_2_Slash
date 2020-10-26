package tests

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"time"
)

func mockUserRepoGetRows(mock sqlmock.Sqlmock, id uint64, nickname, email, password, avatar string) {
	expRows := sqlmock.NewRows([]string{"id", "nickname", "email", "password",
		"avatar"}).AddRow(id, nickname, email, password, avatar)
	mock.ExpectQuery("SELECT id, nickname, email, password, avatar FROM profile").WithArgs(id).WillReturnRows(expRows)
}

func mockUserRepoGetErrNoRows(mock sqlmock.Sqlmock, id uint64, nickname, email, password, avatar string) {
	mock.ExpectQuery("SELECT id, nickname, email, password, avatar FROM profile").WithArgs(id).WillReturnError(sql.ErrNoRows)
}

func mockUserRepoDelete(mock sqlmock.Sqlmock, id uint64, lastInsertId int64) {
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM profile").WithArgs(id).WillReturnResult(sqlmock.NewResult(lastInsertId, 1))
	mock.ExpectCommit()
}

func mockUserRepoUpdateEmail(mock sqlmock.Sqlmock, id uint64, newEmail string, lastInsertId int64) {
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE profile"+
		" SET email").WithArgs(newEmail, id).WillReturnResult(sqlmock.NewResult(lastInsertId, 1))
	mock.ExpectCommit()
}

func mockUserRepoSelectByEmailReturnRows(mock sqlmock.Sqlmock, id uint64, nickname, email, password, avatar string) {
	rows := sqlmock.NewRows([]string{"id", "nickname", "email", "password", "avatar"})
	rows.AddRow(id, nickname, email, password, avatar)
	mock.ExpectQuery("SELECT id, nickname, email, password, avatar " +
		"FROM profile").WithArgs(email).WillReturnRows(rows)
}

func mockUserRepoSelectByEmailReturnErrNoRows(mock sqlmock.Sqlmock, id uint64, nickname, email, password, avatar string) {
	mock.ExpectQuery("SELECT id, nickname, email, password, avatar " +
		"FROM profile").WithArgs(email).WillReturnError(sql.ErrNoRows)
}

func mockUserRepoInsertReturnRows(mock sqlmock.Sqlmock, id uint64, nickname, email, password, avatar string) {
	mock.ExpectBegin()
	insertAnswer := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery("INSERT INTO profile").WithArgs(nickname, email,
		password, avatar).WillReturnRows(insertAnswer)
	mock.ExpectCommit()
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnySessionValue struct{}

func (a AnySessionValue) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

func mockInsertSessionReturnRows(mock sqlmock.Sqlmock, sessId uint64, uid int) {
	mock.ExpectBegin()
	idRows := sqlmock.NewRows([]string{"id"})
	idRows.AddRow(sessId)
	mock.ExpectQuery("INSERT INTO session").WithArgs(AnySessionValue{},
		AnyTime{}, uid).WillReturnRows(idRows)
	mock.ExpectCommit()
}

func mockGetUserSessionReturnErrNoRows(mock sqlmock.Sqlmock, id uint64, sessId int) {
	mock.ExpectQuery("SELECT id, value, expires, profile_id " +
		"FROM session").WithArgs(id).WillReturnError(sql.ErrNoRows)
}

func mockGetUserSessionReturnRows(mock sqlmock.Sqlmock, sess string, sessId int, sessExpires time.Time, userId uint64) {
	answerRows := sqlmock.NewRows([]string{"id", "value", "expires", "profile_id"})
	answerRows.AddRow(sessId, sess, sessExpires, userId)
	mock.ExpectQuery("SELECT id, value, expires, profile_id " +
		"FROM session").WithArgs(sess).WillReturnRows(answerRows)
}

func mockCheckAndInsertUser(mock sqlmock.Sqlmock, id uint64, nickname, email, password, avatar string) {
	mockUserRepoSelectByEmailReturnErrNoRows(mock, id, nickname, email, password, avatar)
	mockUserRepoInsertReturnRows(mock, id, nickname, email, password, avatar)
}

func mockCheckAndInsertSession(mock sqlmock.Sqlmock, id uint64, sessId int) {
	mockGetUserSessionReturnErrNoRows(mock, id, sessId)
	mockInsertSessionReturnRows(mock, uint64(sessId), int(id))
}

func mockAlreadyExistTest(mock sqlmock.Sqlmock, id int, nickname, email, password, avatar string) {
	mockUserRepoSelectByEmailReturnRows(mock, uint64(id), nickname, email, password, avatar)
}

func mockGetProfileSuccess(mock sqlmock.Sqlmock, sess string, sessId int, sessExpires time.Time, userId uint64, nickname, email, password, avatar string) {
	mockGetUserSessionReturnRows(mock, sess, sessId, sessExpires, userId)
	// check
	mockUserRepoGetRows(mock, userId, nickname, email, password, avatar)
	// get
	mockUserRepoGetRows(mock, userId, nickname, email, password, avatar)
}
