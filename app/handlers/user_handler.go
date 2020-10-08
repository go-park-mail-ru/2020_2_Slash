package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2020_2_Slash/app/helpers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/session"
	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
)

type UserHandler struct {
	UserRepo       *user.UserRepo
	SessionManager *session.SessionManager
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		UserRepo:       user.NewUserRepo(),
		SessionManager: session.NewSessionManager(),
	}
}

type UserInput struct {
	Nickname         string `json:"nickname"`
	Email            string `json:"email"`
	Password         string `json:"password,omitempty"`
	RepeatedPassword string `json:"repeated_password,omitempty"`
}

type Error struct {
	Message string `json:"error"`
}

type Result struct {
	Message string `json:"result"`
}

func WriteResponse(w http.ResponseWriter, body interface{}, status int) {
	res, err := json.Marshal(body)
	if err != nil {
		log.Println("Error in decoding responce data: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	w.Write(res)
}

func WriteErrorResponse(w http.ResponseWriter, msg string, status int) {
	log.Println(msg)
	data := Error{Message: msg}
	WriteResponse(w, data, status)
	return
}

func CreateUser(userInput *UserInput) (*user.User, error) {
	if userInput.Email == "" || userInput.Password == "" || userInput.RepeatedPassword == "" {
		return nil, errors.New("not enough input data")
	}
	if !helpers.IsValidEmail(userInput.Email) {
		return nil, errors.New("email is invalid")
	}
	if userInput.Password != userInput.RepeatedPassword {
		return nil, errors.New("passwords don't match")
	}
	if userInput.Nickname == "" {
		// If nickname wasn't sent, use email before @
		nickname := strings.Split(userInput.Email, "@")[0]
		userInput.Nickname = nickname
	}
	user := &user.User{
		Nickname: userInput.Nickname,
		Email:    userInput.Email,
		Password: userInput.Password,
	}
	return user, nil
}

func CreateCookie(session *session.Session) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	}
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	userInput := &UserInput{}
	err := decoder.Decode(userInput)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := CreateUser(userInput)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = uh.UserRepo.Register(user)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusConflict)
		return
	}
	log.Println("Registred user: ", user)

	session := uh.SessionManager.Create(user)
	cookie := CreateCookie(session)
	http.SetCookie(w, cookie)
	data := Result{Message: "ok"}
	WriteResponse(w, data, http.StatusCreated)
}

func (uh *UserHandler) GetValidSession(cookieVal string) (*session.Session, bool) {
	session, has := uh.SessionManager.Get(cookieVal)
	if !has || !uh.SessionManager.IsValid(session) || !uh.UserRepo.Exists(session.UserID) {
		return nil, false
	}
	return session, true
}

func (uh *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		WriteErrorResponse(w, "user isn't authorized", http.StatusUnauthorized)
		return
	}
	session, valid := uh.GetValidSession(cookie.Value)
	if !valid {
		WriteErrorResponse(w, "session is invalid", http.StatusUnauthorized)
		return
	}

	curUser, _ := uh.UserRepo.Get(session.UserID)
	userProfile := curUser.GetProfile()
	WriteResponse(w, userProfile, http.StatusOK)
}

func (uh *UserHandler) ChangeUserProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		WriteErrorResponse(w, "user isn't authorized", http.StatusUnauthorized)
		return
	}
	session, valid := uh.GetValidSession(cookie.Value)
	if !valid {
		WriteErrorResponse(w, "session is invalid", http.StatusUnauthorized)
		return
	}

	inputData := make(map[string]string)
	err = json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		log.Println("Error in decoding input data: ", err)
		data := Error{Message: err.Error()}
		WriteResponse(w, data, http.StatusBadRequest)
		return
	}

	curUser, _ := uh.UserRepo.Get(session.UserID)

	if nickname, has := inputData["nickname"]; has && nickname != "" {
		curUser.Nickname = nickname
	}
	if email, has := inputData["email"]; has && helpers.IsValidEmail(email) {
		err := uh.UserRepo.UpdateEmail(curUser.ID, email)
		if err != nil {
			log.Println("Email already exists")
			data := Error{Message: "email already exists"}
			WriteResponse(w, data, http.StatusBadRequest)
			return
		}
	}

	data := Result{Message: "ok"}
	WriteResponse(w, data, http.StatusOK)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newUser, err := getUserFromRequest(r)
	if err != nil {
		log.Println("Error in decoding user data: ", err)
		data := Error{Message: err.Error()}
		WriteResponse(w, data, http.StatusBadRequest)
		return
	}

	if err, ok := isUserDataValid(newUser); !ok {
		WriteResponse(w, err, http.StatusBadRequest)
		return
	}

	dbUser, ok := h.UserRepo.GetByEmail(newUser.Email)
	if !ok {
		data := Error{Message: WrongEmailMsg}
		WriteResponse(w, data, http.StatusBadRequest)
		return
	}
	if err, ok := isPasswordRight(dbUser, newUser); !ok {
		WriteResponse(w, err, http.StatusBadRequest)
		return
	}

	// save session to db
	session := h.SessionManager.Create(dbUser)
	// set cookie in browser
	cookie := CreateCookie(session)
	http.SetCookie(w, cookie)

	data := NewLoginResponse(dbUser.ID, dbUser.Nickname, dbUser.Avatar)
	WriteResponse(w, data, http.StatusOK)
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

func SetOverdueCookie(w http.ResponseWriter, session *http.Cookie) {
	session.Path = "/"
	session.Expires = time.Now().AddDate(0, 0, -2)
	http.SetCookie(w, session)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		data := Error{Message: UserUnauthorizedMsg}
		WriteResponse(w, data, http.StatusUnauthorized)
		return
	}

	h.SessionManager.Delete(session.Value)

	SetOverdueCookie(w, session)

	data := Result{"ok"}
	WriteResponse(w, data, http.StatusOK)
}

type SessionResponse struct {
	Status string `json:"status"`
}

func (h *UserHandler) CheckSession(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		data := SessionResponse{Status: "unauthorized"}
		WriteResponse(w, data, http.StatusUnauthorized)
		return
	}

	cookie, has := h.SessionManager.Get(session.Value)
	if !has {
		data := SessionResponse{Status: "unauthorized"}
		WriteResponse(w, data, http.StatusUnauthorized)
		return
	}
	isValid := h.SessionManager.IsValid(cookie)
	if !isValid {
		data := SessionResponse{Status: "unauthorized"}
		WriteResponse(w, data, http.StatusUnauthorized)
		return
	}

	data := SessionResponse{Status: "authorized"}
	WriteResponse(w, data, http.StatusOK)
	return
}

func (uh *UserHandler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		WriteErrorResponse(w, "user isn't authorized", http.StatusUnauthorized)
		return
	}
	session, valid := uh.GetValidSession(cookie.Value)
	if !valid {
		WriteErrorResponse(w, "session is invalid", http.StatusUnauthorized)
		return
	}

	// Get form
	var maxMemory int64 = 10 << 20 // 10Mb
	r.Body = http.MaxBytesReader(w, r.Body, maxMemory+1024)
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		if err.Error() == "http: request body too large" {
			WriteErrorResponse(w, "image too large", http.StatusRequestEntityTooLarge)
			return
		}
		WriteErrorResponse(w, "avatar field is expected", http.StatusBadRequest)
		return
	}

	// Get image
	imageFile, _, err := r.FormFile("avatar")
	if err != nil {
		WriteErrorResponse(w, "avatar field is expected", http.StatusBadRequest)
		return
	}
	defer imageFile.Close()

	// Check content type of image
	fileHeader := make([]byte, 512)
	if _, err := imageFile.Read(fileHeader); err != nil {
		WriteErrorResponse(w, "error reading image file", http.StatusBadRequest)
		return
	}
	fileExtension, allowed := helpers.IsAllowedImageContentType(fileHeader)
	if !allowed {
		WriteErrorResponse(w, "this content type is prohibited", http.StatusBadRequest)
		return
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer fileInStorage.Close()

	if _, err := io.Copy(fileInStorage, imageFile); err != nil {
		log.Println("error in copying image file: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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
	WriteResponse(w, data, http.StatusOK)
}
