package grpc

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/sanitizer"
	"golang.org/x/crypto/bcrypt"
	_ "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"strings"
)

type UserblockMicroservice struct {
	userRepo user.UserRepository
}

func NewUserblockMicroservice(userRepo user.UserRepository) *UserblockMicroservice {
	return &UserblockMicroservice{userRepo: userRepo}
}

func (uu *UserblockMicroservice) Create(ctx context.Context, newUser *user.User) (*user.User, error) {
	sanitizer.Sanitize(newUser)
	if err := uu.checkByEmail(newUser.Email); err == nil {
		return nil, status.Error(codes.Code(consts.CodeEmailAlreadyExists), "")
	}

	// If nickname wasn't sent, set nickname as email before @
	if newUser.Nickname == "" {
		newUser.Nickname = strings.Split(newUser.Email, "@")[0]
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	newUser.Password = string(hashedPassword)

	modelUser := user.GrpcUserToModel(newUser)
	if err := uu.userRepo.Insert(modelUser); err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	newUser.ID = modelUser.ID

	return newUser, nil
}

func (uu *UserblockMicroservice) GetByEmail(ctx context.Context, email *user.Email) (*user.User, error) {
	dbUser, err := uu.userRepo.SelectByEmail(email.GetEmail())
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeUserDoesNotExist), "")
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), "")
	}
	return user.ModelUserToGrpc(dbUser), nil
}

func (uu *UserblockMicroservice) GetByID(ctx context.Context, id *user.ID) (*user.User, error) {
	dbUser, err := uu.userRepo.SelectByID(id.GetID())
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeUserDoesNotExist), "")
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), "")
	}
	return user.ModelUserToGrpc(dbUser), nil
}

func (uu *UserblockMicroservice) UpdateProfile(ctx context.Context, newUserData *user.User) (*user.User, error) {
	sanitizer.Sanitize(newUserData)
	dbUser, err := uu.GetByID(context.Background(), &user.ID{ID: newUserData.GetID()})
	if err != nil {
		return nil, err
	}

	// Update email
	if newUserData.Email != "" && dbUser.Email != newUserData.Email {
		if err := uu.checkByEmail(newUserData.Email); err == nil {
			return nil, status.Error(codes.Code(consts.CodeEmailAlreadyExists), "")
		}
		dbUser.Email = newUserData.Email
	}

	// Update nickname
	if newUserData.Nickname != "" {
		dbUser.Nickname = newUserData.Nickname
	}

	// Update password
	if newUserData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUserData.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, status.Error(codes.Code(consts.CodeInternalError), "")
		}
		dbUser.Password = string(hashedPassword)
	}

	if err := uu.userRepo.Update(user.GrpcUserToModel(dbUser)); err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), "")
	}
	return dbUser, nil
}

func (uu *UserblockMicroservice) UpdateAvatar(ctx context.Context, idAvatar *user.IdAvatar) (*user.User, error) {
	dbUser, err := uu.GetByID(context.Background(), &user.ID{ID: idAvatar.GetId().GetID()})
	if err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	// Update user avatar
	prevAvatar := dbUser.Avatar
	dbUser.Avatar = idAvatar.GetAvatar().GetAvatar()
	if err := uu.userRepo.Update(user.GrpcUserToModel(dbUser)); err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	// Delete prev avatar image
	if prevAvatar != "" {
		if err := os.Remove("." + prevAvatar); err != nil {
			return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
		}
	}
	return dbUser, nil
}

func (uu *UserblockMicroservice) checkByEmail(email string) error {
	_, err := uu.GetByEmail(context.Background(), &user.Email{Email: email})
	return err
}
