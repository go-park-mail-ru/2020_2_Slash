package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

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
	w.WriteHeader(status)
	w.Write(res)
}

func CreateUser(userInput *UserInput) (*user.User, error) {
	if userInput.Email == "" || userInput.Password == "" || userInput.RepeatedPassword == "" {
		return nil, errors.New("Not enough input data")
	}
	if userInput.Password != userInput.RepeatedPassword {
		return nil, errors.New("Passwords don't match")
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

func createCookie(session *session.Session) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	userInput := &UserInput{}
	err := decoder.Decode(userInput)
	if err != nil {
		log.Println("Error in decoding user data: ", err)
		data := Error{Message: err.Error()}
		WriteResponse(w, data, http.StatusBadRequest)
		return
	}

	user, err := CreateUser(userInput)
	if err != nil {
		log.Println("Error in creating user: ", err)
		data := Error{Message: err.Error()}
		WriteResponse(w, data, http.StatusBadRequest)
		return
	}

	err = h.UserRepo.Register(user)
	if err != nil {
		log.Println("Error:", err)
		data := Error{Message: err.Error()}
		WriteResponse(w, data, http.StatusConflict)
		return
	}
	log.Println("Registred user: ", user)

	session := h.SessionManager.Create(user)
	cookie := createCookie(session)
	http.SetCookie(w, cookie)
	data := Result{Message: "ok"}
	WriteResponse(w, data, http.StatusCreated)
}
