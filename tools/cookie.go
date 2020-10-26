package tools

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"net/http"
)

func CreateCookie(sess *models.Session) *http.Cookie {
	return &http.Cookie{
		Name:     consts.SessionName,
		Value:    sess.Value,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  sess.ExpiresAt,
		HttpOnly: true,
	}
}
