package helpers

import (
	"errors"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	cstm_errors "github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	uuid "github.com/satori/go.uuid"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
)

var allowedImagesContentType = map[string]string{
	"image/png":  "png",
	"image/jpg":  "jpg",
	"image/jpeg": "jpeg",
}

func determineFileContentType(fileHeader textproto.MIMEHeader) (string, error) {
	contentTypes := fileHeader["Content-Type"]
	if len(contentTypes) < 1 {
		return "", errors.New("Wrong file header")
	}
	return contentTypes[0], nil
}

func CheckImageContentType(image *multipart.FileHeader) *cstm_errors.Error {
	// Check content type from header
	if !IsAllowedImageHeader(image) {
		return cstm_errors.Get(CodeWrongImgExtension)
	}

	// Check real content type
	imageFile, err := image.Open()
	if err != nil {
		return cstm_errors.New(CodeBadRequest, err)
	}
	defer imageFile.Close()

	fileHeader := make([]byte, 512)
	if _, err := imageFile.Read(fileHeader); err != nil {
		return cstm_errors.New(CodeBadRequest, err)
	}

	if !IsAllowedImageContentType(fileHeader) {
		return cstm_errors.Get(CodeWrongImgExtension)
	}
	return nil
}

func IsAllowedImageHeader(image *multipart.FileHeader) bool {
	contentType, err := determineFileContentType(image.Header)
	if err != nil {
		return false
	}
	_, allowed := allowedImagesContentType[contentType]
	return allowed
}

func GetImageExtension(image *multipart.FileHeader) (string, error) {
	contentType, err := determineFileContentType(image.Header)
	if err != nil {
		return "", err
	}

	extension, has := allowedImagesContentType[contentType]
	if !has {
		return "", errors.New("prohibited image extension")
	}
	return extension, nil
}

func GetImageContentType(image []byte) (string, bool) {
	contentType := http.DetectContentType(image)
	extension, allowed := allowedImagesContentType[contentType]
	return extension, allowed
}

func IsAllowedImageContentType(image []byte) bool {
	_, allowed := GetImageContentType(image)
	return allowed
}

func GetUniqFileName(userID uint64, fileExtension string) string {
	randString := uuid.NewV4().String()
	return "userid_" + strconv.Itoa(int(userID)) + "_" + randString + "." + fileExtension
}

func InitAvatarStorage(path string) {
	mode := int(0777)
	os.Mkdir(path, os.FileMode(mode))
}
