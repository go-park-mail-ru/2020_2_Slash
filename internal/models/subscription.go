package models

import "time"

type Subscription struct {
	ID      uint64    `json:"id"`
	UserID  uint64    `json:"user_id"`
	Expires time.Time `json:"expires"`
	Active  bool      `json:"active"`
}
