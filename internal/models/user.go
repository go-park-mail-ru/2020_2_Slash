package models

type User struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Avatar   string `json:"avatar"`
}

func (u *User) Sanitize() *User {
	u.Password = ""
	return u
}
