package helpers

import (
	uuid "github.com/satori/go.uuid"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\" +
	"/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?" +
	"(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	if !emailRegex.MatchString(email) {
		return false
	}

	parts := strings.Split(email, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}

var allowedImagesContentType = map[string]string{
	"image/png":  "png",
	"image/jpg":  "jpg",
	"image/jpeg": "jpeg",
}

func IsAllowedImageContentType(image []byte) (string, bool) {
	contentType := http.DetectContentType(image)
	extension, allowed := allowedImagesContentType[contentType]
	return extension, allowed
}

func GetUniqFileName(user *user.User, fileExtension string) string {
	randString := uuid.NewV4().String()
	return "userid_" + strconv.Itoa(int(user.ID)) + "_" + randString + "." + fileExtension
}

func InitAvatarStorage() {
	path := "./avatars" // TODO: config
	mode := int(0777)
	os.Mkdir(path, os.FileMode(mode))
}
