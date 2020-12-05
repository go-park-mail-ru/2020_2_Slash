package grpc

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/config"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user/repository"
	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Loader
)

func TestMain(m *testing.M) {
	config, configErr := config.LoadConfig("../../../../config.json")
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
		testfixtures.Directory("../../test/fixture"),
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

func TestUserblockMicroservice_Create_DB_OK(t *testing.T) {
	prepareTestDatabase()
	userRep := repository.NewUserPgRepository(db)
	userblockMicroservice := NewUserblockMicroservice(userRep)
	userBuilder := NewUserBuilder()
	user := userBuilder.CreateRegularUser()

	createdUser, err := userblockMicroservice.Create(context.Background(), user)

	// fixture id logic
	assert.NoError(t, err)
	assert.Equal(t, uint64(10001), createdUser.ID)
}

func TestUserblockMicroservice_Create_DB_EmailConflict(t *testing.T) {
	prepareTestDatabase()
	userRep := repository.NewUserPgRepository(db)
	userblockMicroservice := NewUserblockMicroservice(userRep)
	userBuilder := NewUserBuilder()
	existedUser := userBuilder.CreateNinthUserFromDB()

	createdUser, err := userblockMicroservice.Create(context.Background(), existedUser)

	assert.Equal(t, status.Error(codes.Code(consts.CodeEmailAlreadyExists),
		""), err)
	assert.Nil(t, createdUser)
}

func TestUserblockMicroservice_Create_DB_EmptyNickname(t *testing.T) {
	prepareTestDatabase()
	userRep := repository.NewUserPgRepository(db)
	userblockMicroservice := NewUserblockMicroservice(userRep)
	userBuilder := NewUserBuilder()
	emptyNicknameUser := userBuilder.CreateUserWithEmptyNickname()

	createdUser, err := userblockMicroservice.Create(context.Background(), emptyNicknameUser)

	assert.NoError(t, err)
	assert.Equal(t, "nick", createdUser.Nickname)
}

func TestUserblockMicroservice_UpdateProfile_DB_OK(t *testing.T) {
	prepareTestDatabase()
	userRep := repository.NewUserPgRepository(db)
	userblockMicroservice := NewUserblockMicroservice(userRep)
	userBuilder := NewUserBuilder()
	userForUpdate := userBuilder.CreateNinthUserFromDB()
	newNickname, newEmail := "new nickname", "newEmail@mail.ru"
	userForUpdate.Nickname = newNickname
	userForUpdate.Email = newEmail

	updatedUser, err := userblockMicroservice.UpdateProfile(context.Background(), userForUpdate)

	assert.NoError(t, err)
	assert.Equal(t, newNickname, updatedUser.Nickname)
	assert.Equal(t, newEmail, updatedUser.Email)
}

func TestUserblockMicroservice_UpdateProfile_DB_EmailConflicts(t *testing.T) {
	prepareTestDatabase()
	userRep := repository.NewUserPgRepository(db)
	userblockMicroservice := NewUserblockMicroservice(userRep)
	userBuilder := NewUserBuilder()
	userForUpdate := userBuilder.CreateUserWithConflictEmail()

	updatedUser, err := userblockMicroservice.UpdateProfile(context.Background(), userForUpdate)

	assert.Equal(t, status.Error(codes.Code(consts.CodeEmailAlreadyExists),
		""), err)
	assert.Nil(t, updatedUser)
}

func TestUserblockMicroservice_UpdateProfile_DB_IdDoesNotExist(t *testing.T) {
	prepareTestDatabase()
	userRep := repository.NewUserPgRepository(db)
	userblockMicroservice := NewUserblockMicroservice(userRep)
	userBuilder := NewUserBuilder()
	userForUpdate := userBuilder.CreateUserWithNotExistedID()

	updatedUser, err := userblockMicroservice.UpdateProfile(context.Background(), userForUpdate)

	assert.Equal(t, status.Error(codes.Code(consts.CodeUserDoesNotExist),
		""), err)
	assert.Nil(t, updatedUser)
}

func TestUserblockMicroservice_GetByID_DB_OK(t *testing.T) {
	prepareTestDatabase()
	userRep := repository.NewUserPgRepository(db)
	userblockMicroservice := NewUserblockMicroservice(userRep)
	userBuilder := NewUserBuilder()
	ninthUserFromDB := userBuilder.CreateNinthUserFromDB()

	userFromDB, err := userblockMicroservice.GetByID(context.Background(), &ID{ID: ninthUserFromDB.ID})

	assert.NoError(t, err)
	assert.Equal(t, ninthUserFromDB, userFromDB)
}

func TestUserblockMicroservice_GetByID_DB_NoUserWithThisID(t *testing.T) {
	prepareTestDatabase()
	userRep := repository.NewUserPgRepository(db)
	userblockMicroservice := NewUserblockMicroservice(userRep)
	userBuilder := NewUserBuilder()
	userWithNotExistedID := userBuilder.CreateUserWithNotExistedID()

	userFromDB, err := userblockMicroservice.GetByID(context.Background(), &ID{ID: userWithNotExistedID.ID})

	assert.Equal(t, status.Error(codes.Code(consts.CodeUserDoesNotExist),
		""), err)
	assert.Nil(t, userFromDB)
}

// Can't test with real DB because method removes file
func TestUserUseCase_UpdateAvatar_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRep := mocks.NewMockUserRepository(ctrl)
	userblockMicroservice := NewUserblockMicroservice(userRep)
	userBuilder := NewUserBuilder()
	regularUser := userBuilder.CreateRegularUser()
	newAvatar := "/avatar"

	idAvatar := &IdAvatar{
		Id:     &ID{ID: regularUser.ID},
		Avatar: &Avatar{Avatar: newAvatar},
	}

	newUserModel := regularUser
	newUserModel.Avatar = newAvatar

	userRep.
		EXPECT().
		SelectByID(gomock.Eq(regularUser.ID)).
		Return(GrpcUserToModel(regularUser), nil)

	userRep.
		EXPECT().
		Update(gomock.Eq(GrpcUserToModel(regularUser))).
		Return(nil)

	_, err := userblockMicroservice.UpdateAvatar(context.Background(), idAvatar)
	assert.NotEqual(t, err, (error)(nil))
}

func TestUserUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRep := mocks.NewMockUserRepository(ctrl)
	userblockMicroservice := NewUserblockMicroservice(userRep)
	userBuilder := NewUserBuilder()
	regularUser := userBuilder.CreateRegularUser()
	modelUser := GrpcUserToModel(regularUser)

	userRep.
		EXPECT().
		SelectByEmail(modelUser.Email).
		Return(modelUser, nil)

	_, err := userblockMicroservice.Create(context.Background(), regularUser)
	assert.Equal(t, err, status.Error(codes.Code(consts.CodeEmailAlreadyExists), ""))
}

func TestUserUseCase_Create_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRep := mocks.NewMockUserRepository(ctrl)
	userblockMicroservice := NewUserblockMicroservice(userRep)

	userBuilder := NewUserBuilder()
	regularUser := userBuilder.CreateRegularUser()
	modelUser := GrpcUserToModel(regularUser)

	userRep.
		EXPECT().
		SelectByEmail(gomock.Eq(modelUser.Email)).
		Return(modelUser, nil)

	_, err := userblockMicroservice.Create(context.Background(), regularUser)
	assert.Equal(t, err, status.Error(codes.Code(consts.CodeEmailAlreadyExists), ""))
}

func TestUserUseCase_Update_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRep := mocks.NewMockUserRepository(ctrl)
	userblockMicroservice := NewUserblockMicroservice(userRep)

	userBuilder := NewUserBuilder()
	regularUser := userBuilder.CreateRegularUser()
	userInDatabaseModel := GrpcUserToModel(regularUser)
	userInDatabaseModel.Email = "oldEmail@mail.ru"
	existedUser := userBuilder.CreateNinthUserFromDB()
	existedUserModel := GrpcUserToModel(existedUser)

	userRep.
		EXPECT().
		SelectByID(regularUser.ID).
		Return(userInDatabaseModel, nil)

	userRep.
		EXPECT().
		SelectByEmail(regularUser.Email).
		Return(existedUserModel, nil)

	_, err := userblockMicroservice.UpdateProfile(context.Background(), regularUser)
	assert.Equal(t, err, status.Error(codes.Code(consts.CodeEmailAlreadyExists), ""))
}

func TestUserUseCase_GetByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRep := mocks.NewMockUserRepository(ctrl)
	userblockMicroservice := NewUserblockMicroservice(userRep)

	userBuilder := NewUserBuilder()
	regularUser := userBuilder.CreateRegularUser()
	modelUser := GrpcUserToModel(regularUser)

	userRep.
		EXPECT().
		SelectByID(gomock.Eq(regularUser.ID)).
		Return(modelUser, nil)

	dbUser, err := userblockMicroservice.GetByID(context.Background(), &ID{ID: regularUser.ID})
	assert.Equal(t, err, (error)(nil))
	assert.Equal(t, regularUser, dbUser)
}

func TestUserUseCase_GetByEmail_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRep := mocks.NewMockUserRepository(ctrl)
	userblockMicroservice := NewUserblockMicroservice(userRep)

	userBuilder := NewUserBuilder()
	regularUser := userBuilder.CreateRegularUser()
	modelUser := GrpcUserToModel(regularUser)

	userRep.
		EXPECT().
		SelectByEmail(gomock.Eq(regularUser.Email)).
		Return(modelUser, nil)

	dbUser, err := userblockMicroservice.GetByEmail(context.Background(), &Email{Email: regularUser.Email})
	assert.Equal(t, err, (error)(nil))
	assert.Equal(t, dbUser, regularUser)
}
