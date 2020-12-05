package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/config"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user/delivery/grpc"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user/mocks"
	"github.com/stretchr/testify/assert"
)

var userInst = &models.User{
	Nickname: "Jhon",
	Email:    "jhon@gmail.com",
	Password: "hardpassword",
}

var (
	db       *sql.DB
	fixtures *testfixtures.Loader
)

func TestMain(m *testing.M) {
	config, configErr := config.LoadConfig("../../../config.json")
	if configErr != nil {
		log.Fatal(configErr)
	}

	var err error
	db, err = sql.Open("postgres", config.GetTestDbConnString())
	if err != nil {
		log.Fatal(err)
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("../test/fixture"),
	)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func TestUserPgRepository_Insert_DB_OK(t *testing.T) {
	prepareTestDatabase()
	userRep := NewUserPgRepository(db)
	userBuilder := grpc.NewUserBuilder()
	regularUser := userBuilder.CreateRegularUserModel()

	err := userRep.Insert(regularUser)

	assert.NoError(t, err)
	assert.Equal(t, uint64(10001), regularUser.ID)
}

func TestUserPgRepository_Insert_DB_EmailConflict(t *testing.T) {
	prepareTestDatabase()
	userRep := NewUserPgRepository(db)
	userBuilder := grpc.NewUserBuilder()
	userWithConflictEmail := userBuilder.CreateUserWithConflictEmailModel()

	err := userRep.Insert(userWithConflictEmail)

	assert.Error(t, err)
}

func TestUserPgRepository_SelectByEmail_DB_OK(t *testing.T) {
	prepareTestDatabase()
	userRep := NewUserPgRepository(db)
	userBuilder := grpc.NewUserBuilder()
	existedUser := userBuilder.CreateNinthUserFromDBModel()

	userFromDB, err := userRep.SelectByEmail(existedUser.Email)

	assert.NoError(t, err)
	assert.Equal(t, existedUser, userFromDB)
}

func TestUserPgRepository_SelectByEmail_DB_NoUser(t *testing.T) {
	prepareTestDatabase()
	userRep := NewUserPgRepository(db)
	userBuilder := grpc.NewUserBuilder()
	regularUser := userBuilder.CreateRegularUserModel()

	userFromDB, err := userRep.SelectByEmail(regularUser.Email)

	assert.Equal(t, sql.ErrNoRows, err)
	assert.Nil(t, userFromDB)
}

func TestUserPgRepository_SelectByID_DB_OK(t *testing.T) {
	prepareTestDatabase()
	userRep := NewUserPgRepository(db)
	userBuilder := grpc.NewUserBuilder()
	existedUser := userBuilder.CreateNinthUserFromDBModel()

	userFromDB, err := userRep.SelectByID(existedUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, existedUser, userFromDB)
}

func TestUserPgRepository_SelectByID_DB_NoUser(t *testing.T) {
	prepareTestDatabase()
	userRep := NewUserPgRepository(db)
	userBuilder := grpc.NewUserBuilder()
	regularUser := userBuilder.CreateRegularUserModel()

	userFromDB, err := userRep.SelectByID(regularUser.ID)

	assert.Equal(t, sql.ErrNoRows, err)
	assert.Nil(t, userFromDB)
}

func TestUserPgRepository_Update_DB_OK(t *testing.T) {
	prepareTestDatabase()
	userRep := NewUserPgRepository(db)
	userBuilder := grpc.NewUserBuilder()
	existedUser := userBuilder.CreateNinthUserFromDBModel()

	err := userRep.Update(existedUser)

	assert.NoError(t, err)
}

func TestUserPgRepository_Update_DB_EmailConflicts(t *testing.T) {
	prepareTestDatabase()
	userRep := NewUserPgRepository(db)
	userBuilder := grpc.NewUserBuilder()
	userWithConflictEmail := userBuilder.CreateUserWithConflictEmailModel()

	err := userRep.Update(userWithConflictEmail)

	assert.Error(t, err)
}

func TestUserPgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	userPgRep := NewUserPgRepository(db)
	mocks.MockUserRepoInsertReturnRows(mock, userInst)

	err = userPgRep.Insert(userInst)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPgRepository_Insert_UserAlreadyExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	userPgRep := NewUserPgRepository(db)

	mocks.MockUserRepoInsertReturnRows(mock, userInst)
	err = userPgRep.Insert(userInst)
	assert.NoError(t, err)

	mocks.MockUserRepoInsertReturnErrNoUniq(mock, userInst)
	err = userPgRep.Insert(userInst)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPgRepository_Update_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	userPgRep := NewUserPgRepository(db)
	mocks.MockUserRepoUpdateReturnResultOk(mock, userInst)

	err = userPgRep.Update(userInst)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPgRepository_SelectByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	userPgRep := NewUserPgRepository(db)
	mocks.MockUserRepoSelectByIDReturnRows(mock, userInst)

	dbUser, err := userPgRep.SelectByID(userInst.ID)

	assert.Equal(t, userInst, dbUser)
	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPgRepository_SelectByEmail_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	userPgRep := NewUserPgRepository(db)
	mocks.MockUserRepoSelectByEmailReturnRows(mock, userInst)

	dbUser, err := userPgRep.SelectByEmail(userInst.Email)

	assert.Equal(t, userInst, dbUser)
	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
