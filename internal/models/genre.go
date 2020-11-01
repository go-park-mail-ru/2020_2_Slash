package models

type Genre struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type ContentGenre struct {
	ContentID uint64 `json:"content_id"`
	GenreID   uint64 `json:"genre_id"`
}
