package models

type Rating struct {
	UserID    uint64 `json:"-"`
	ContentID uint64 `json:"content_id"`
	Likes     bool   `json:"likes"`
}
