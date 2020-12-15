package CSRFManager

import (
	// #nosec
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"strconv"
	"strings"
	"time"
)

func CreateToken(sess *models.Session) (string, *errors.Error) {
	hashed, err := hashSessData(sess)
	if err != nil {
		return "", err
	}

	token := fmt.Sprintf("%s:%d", hashed, sess.ExpiresAt.Unix())
	return token, nil
}

func ValidateToken(sess *models.Session, actualToken string) *errors.Error {
	expected, customErr := hashSessData(sess)
	if customErr != nil {
		return customErr
	}

	tokenParts := strings.Split(actualToken, ":")
	if len(tokenParts) != 2 {
		return errors.Get(CodeWrongCSRFToken)
	}

	actual := tokenParts[0]
	if actual != expected {
		return errors.Get(CodeWrongCSRFToken)
	}

	tokenExpires, err := strconv.ParseInt(tokenParts[1], 10, 64)
	if err != nil {
		return errors.New(CodeWrongCSRFToken, err)
	}
	if tokenExpires < time.Now().Unix() {
		return errors.New(CodeSessionExpired, err)
	}

	return nil
}

func hashSessData(sess *models.Session) (string, *errors.Error) {
	data := fmt.Sprintf("%d:%s", sess.UserID, sess.Value)
	// #nosec
	hasher := sha1.New()
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", errors.New(CodeInternalError, err)
	}

	hashed := hex.EncodeToString(hasher.Sum(nil))
	return hashed, nil
}
