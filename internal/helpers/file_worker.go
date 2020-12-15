package helpers

import (
	"errors"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	cstm_errors "github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	uuid "github.com/satori/go.uuid"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
)

var allowedImagesContentType = map[string]string{
	"image/png":  "png",
	"image/jpg":  "jpg",
	"image/jpeg": "jpeg",
}

var allowedVideoContentType = map[string]string{
	"video/mp4": "mp4",
}

func StoreFile(fileHeader *multipart.FileHeader, absFilePath string) *cstm_errors.Error {
	file, err := fileHeader.Open()
	if err != nil {
		return cstm_errors.New(CodeBadRequest, err)
	}
	defer file.Close()

	// Save file to storage
	fileMode := int(0777)
	newFile, err := os.OpenFile(filepath.Clean(absFilePath), os.O_WRONLY|os.O_CREATE, os.FileMode(fileMode))
	if err != nil {
		return cstm_errors.New(CodeInternalError, err)
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, file); err != nil {
		removeErr := os.Remove(absFilePath)
		if removeErr != nil {
			logger.Error(removeErr)
		}
		return cstm_errors.New(CodeInternalError, err)
	}
	return nil
}

func CheckImageContentType(image *multipart.FileHeader) *cstm_errors.Error {
	return checkFileContentType(image, allowedImagesContentType)
}

func CheckVideoContentType(video *multipart.FileHeader) *cstm_errors.Error {
	return checkFileContentType(video, allowedVideoContentType)
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

func GetUniqFileName(userID uint64, fileExtension string) string {
	randString := uuid.NewV4().String()
	return "userid_" + strconv.Itoa(int(userID)) + "_" + randString + "." + fileExtension
}

func InitStorage(path string) {
	mode := int(0777)
	if err := os.Mkdir(path, os.FileMode(mode)); err != nil {
		logger.Error(err)
	}
}

func InitTree(path string) {
	mode := int(0777)
	if err := os.MkdirAll(path, os.FileMode(mode)); err != nil {
		logger.Error(err)
	}
}

func checkFileContentType(file *multipart.FileHeader, allowedContentTypes map[string]string) *cstm_errors.Error {
	// Check content type from header
	if !isAllowedFileHeader(file, allowedContentTypes) {
		return cstm_errors.Get(CodeWrongImgExtension)
	}

	// Check real content type
	imageFile, err := file.Open()
	if err != nil {
		return cstm_errors.New(CodeBadRequest, err)
	}
	defer imageFile.Close()

	fileHeader := make([]byte, 512)
	if _, err := imageFile.Read(fileHeader); err != nil {
		return cstm_errors.New(CodeBadRequest, err)
	}

	if !isAllowedFileContentType(fileHeader, allowedContentTypes) {
		return cstm_errors.Get(CodeWrongImgExtension)
	}
	return nil
}

func determineFileContentType(fileHeader textproto.MIMEHeader) (string, error) {
	contentTypes := fileHeader["Content-Type"]
	if len(contentTypes) < 1 {
		return "", errors.New("Wrong file header")
	}
	return contentTypes[0], nil
}

func isAllowedFileHeader(file *multipart.FileHeader, allowedContentTypes map[string]string) bool {
	contentType, err := determineFileContentType(file.Header)
	if err != nil {
		return false
	}
	_, allowed := allowedContentTypes[contentType]
	return allowed
}

func getFileContentType(file []byte, allowedContentTypes map[string]string) (string, bool) {
	contentType := http.DetectContentType(file)
	extension, allowed := allowedContentTypes[contentType]
	return extension, allowed
}

func isAllowedFileContentType(file []byte, allowedContentTypes map[string]string) bool {
	_, allowed := getFileContentType(file, allowedContentTypes)
	return allowed
}
