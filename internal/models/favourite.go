package models

import "time"

type Favourite struct {
	UserID    uint64    `json:"-"`
	ContentID uint64    `json:"content_id"`
	Created   time.Time `json:"-"`
}
