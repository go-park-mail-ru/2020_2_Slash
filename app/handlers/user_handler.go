package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2020_2_Slash/app/helpers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/session"
	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
)

type UserHandler struct {
	UserRepo       *user.UserRepo
	SessionManager *session.SessionManager
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		UserRepo:       user.NewUserRepo(db),
		SessionManager: session.NewSessionManager(db),
	}
}

type Error struct {
	Message string `json:"error"`
}

type Result struct {
	Message string `json:"result"`
}

func CreateCookie(session *session.Session) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    session.Value,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	}
}

func (uh *UserHandler) GetValidSession(cookieVal string) (*session.Session, bool) {
	session, has := uh.SessionManager.Get(cookieVal)
	if !has || !uh.SessionManager.IsValid(session) || !uh.UserRepo.Exists(session.UserID) {
		return nil, false
	}
	return session, true
}

func (uh *UserHandler) GetUserProfile(cntx echo.Context) error {
	cookie, err := cntx.Cookie("session_id")
	if err != nil {
		data := Error{Message: "user isn't authorized"}
		return cntx.JSON(http.StatusUnauthorized, data)
	}
	session, valid := uh.GetValidSession(cookie.Value)
	if !valid {
		data := Error{Message: "session is invalid"}
		return cntx.JSON(http.StatusUnauthorized, data)
	}

	curUser, _ := uh.UserRepo.Get(session.UserID)
	userProfile := curUser.GetProfile()
	return cntx.JSON(http.StatusOK, userProfile)
}

func (uh *UserHandler) ChangeUserProfile(cntx echo.Context) error {
	cookie, err := cntx.Cookie("session_id")
	if err != nil {
		data := Error{Message: "user isn't authorized"}
		return cntx.JSON(http.StatusUnauthorized, data)
	}
	session, valid := uh.GetValidSession(cookie.Value)
	if !valid {
		data := Error{Message: "session is invalid"}
		return cntx.JSON(http.StatusUnauthorized, data)
	}

	curUser, _ := uh.UserRepo.Get(session.UserID)

	if nickname := cntx.FormValue("nickname"); nickname != "" {
		curUser.Nickname = nickname
	}
	if email := cntx.FormValue("email"); email != "" && helpers.IsValidEmail(email) {
		err := uh.UserRepo.UpdateEmail(curUser.ID, email)
		if err != nil {
			data := Error{Message: "email already exists"}
			return cntx.JSON(http.StatusBadRequest, data)
		}
	}

	data := Result{Message: "ok"}
	return cntx.JSON(http.StatusOK, data)
}

func (h *UserHandler) Login(cntx echo.Context) error {
	newUser := &user.User{}
	if err := cntx.Bind(newUser); err != nil {
		log.Println("Error in decoding user data: ", err)
		data := Error{Message: err.Error()}
		return cntx.JSON(http.StatusBadRequest, data)
	}

	if err, ok := isUserDataValid(newUser); !ok {
		return cntx.JSON(http.StatusBadRequest, err)
	}

	dbUser, ok := h.UserRepo.GetByEmail(newUser.Email)
	if !ok {
		data := Error{Message: WrongEmailMsg}
		return cntx.JSON(http.StatusBadRequest, data)
	}
	if err, ok := isPasswordRight(dbUser, newUser); !ok {
		return cntx.JSON(http.StatusBadRequest, err)
	}

	// save session to db
	session, err := h.SessionManager.Create(dbUser)
	if err != nil {
		log.Println(err)
		data := Error{Message: err.Error()}
		return cntx.JSON(http.StatusBadRequest, data)
	}
	// set cookie in browser
	cookie := CreateCookie(session)
	cntx.SetCookie(cookie)

	data := NewLoginResponse(dbUser.ID, dbUser.Nickname, dbUser.Avatar)
	return cntx.JSON(http.StatusOK, data)
}

type LoginResponse struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func NewLoginResponse(id uint64, nickname string, avatar string) *LoginResponse {
	return &LoginResponse{
		ID:       id,
		Nickname: nickname,
		Avatar:   avatar,
	}
}

func getUserFromRequest(r *http.Request) (*user.User, error) {
	decoder := json.NewDecoder(r.Body)
	newUser := &user.User{}
	err := decoder.Decode(newUser)
	return newUser, err
}

func isPasswordRight(dbUser *user.User, newUser *user.User) (Error, bool) {
	if newUser.Password != dbUser.Password {
		data := Error{Message: WrongPasswordMsg}
		return data, false
	}
	return Error{}, true
}

func isUserDataValid(newUser *user.User) (Error, bool) {
	if newUser.Email == "" {
		data := Error{Message: EmptyEmailMsg}
		return data, false
	}
	if !helpers.IsValidEmail(newUser.Email) {
		data := Error{Message: InvalidEmailMsg}
		return data, false
	}
	if newUser.Password == "" {
		data := Error{Message: EmptyPasswordMsg}
		return data, false
	}
	return Error{}, true
}

func SetOverdueCookie(cntx echo.Context, session *http.Cookie) {
	session.Path = "/"
	session.Expires = time.Now().AddDate(0, 0, -2)
	cntx.SetCookie(session)
}

func (h *UserHandler) Logout(cntx echo.Context) error {
	session, err := cntx.Cookie("session_id")
	if err == http.ErrNoCookie {
		data := Error{Message: UserUnauthorizedMsg}
		return cntx.JSON(http.StatusUnauthorized, data)
	}

	err = h.SessionManager.Delete(session.Value)
	if err != nil {
		data := Error{Message: err.Error()}
		return cntx.JSON(http.StatusInternalServerError, data)
	}
	SetOverdueCookie(cntx, session)

	data := Result{"ok"}
	return cntx.JSON(http.StatusOK, data)
}

type SessionResponse struct {
	Status string `json:"status"`
}

func (h *UserHandler) CheckSession(cntx echo.Context) error {
	session, err := cntx.Cookie("session_id")
	if err == http.ErrNoCookie {
		data := SessionResponse{Status: "unauthorized"}
		return cntx.JSON(http.StatusUnauthorized, data)
	}

	cookie, has := h.SessionManager.Get(session.Value)
	if !has {
		data := SessionResponse{Status: "unauthorized"}
		return cntx.JSON(http.StatusUnauthorized, data)
	}
	isValid := h.SessionManager.IsValid(cookie)
	if !isValid {
		data := SessionResponse{Status: "unauthorized"}
		return cntx.JSON(http.StatusUnauthorized, data)
	}

	data := SessionResponse{Status: "authorized"}
	return cntx.JSON(http.StatusOK, data)
}

func (uh *UserHandler) SetAvatar(cntx echo.Context) error {
	cookie, err := cntx.Cookie("session_id")
	if err != nil {
		data := Error{Message: "user isn't authorized"}
		return cntx.JSON(http.StatusUnauthorized, data)
	}
	session, valid := uh.GetValidSession(cookie.Value)
	if !valid {
		data := Error{Message: "session is invalid"}
		return cntx.JSON(http.StatusUnauthorized, data)
	}

	// Get image
	image, err := cntx.FormFile("avatar")
	if err != nil {
		data := Error{Message: err.Error()}
		return cntx.JSON(http.StatusBadRequest, data)
	}
	imageFile, err := image.Open()
	if err != nil {
		data := Error{Message: err.Error()}
		return cntx.JSON(http.StatusBadRequest, data)
	}
	defer imageFile.Close()

	// Check content type of image
	fileHeader := make([]byte, 512)
	if _, err := imageFile.Read(fileHeader); err != nil {
		data := Error{Message: "error reading image file"}
		return cntx.JSON(http.StatusBadRequest, data)
	}
	fileExtension, allowed := helpers.IsAllowedImageContentType(fileHeader)
	if !allowed {
		data := Error{Message: "this content type is prohibited"}
		return cntx.JSON(http.StatusBadRequest, data)
	}
	imageFile.Seek(0, 0)

	curUser, _ := uh.UserRepo.Get(session.UserID)
	fileName := helpers.GetUniqFileName(curUser, fileExtension)
	newImageFilePath := "/avatars/" + fileName // TODO: Make it via config
	mode := int(0777)

	// Save image to storage
	fileInStorage, err := os.OpenFile("."+newImageFilePath, os.O_WRONLY|os.O_CREATE, os.FileMode(mode))
	if err != nil {
		log.Println("error in creating image file: ", err)
		data := Error{Message: "internal server error"}
		return cntx.JSON(http.StatusInternalServerError, data)
	}
	defer fileInStorage.Close()

	if _, err := io.Copy(fileInStorage, imageFile); err != nil {
		log.Println("error in copying image file: ", err)
		data := Error{Message: "internal server error"}
		return cntx.JSON(http.StatusInternalServerError, data)
	}

	// Update user avatar
	prevAvatarPath := curUser.Avatar
	if prevAvatarPath != "" {
		_ = os.Remove(prevAvatarPath)
	}
	curUser.Avatar = newImageFilePath

	type ImagePath struct {
		Avatar string `json:"avatar"`
	}
	data := ImagePath{Avatar: fileName}
	return cntx.JSON(http.StatusOK, data)
}
