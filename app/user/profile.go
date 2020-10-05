package user

type UserProfile struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

func (u *User) GetProfile() *UserProfile {
	return &UserProfile{
		Nickname: u.Nickname,
		Email:    u.Email,
		Avatar:   u.Avatar,
	}
}
