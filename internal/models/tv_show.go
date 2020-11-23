package models

type TVShow struct {
	ID      uint64 `json:"id"`
	Seasons int    `json:"seasons"`
	Content
}
