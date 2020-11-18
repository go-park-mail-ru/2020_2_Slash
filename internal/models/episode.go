package models

type Episode struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Video       string `json:"video"`
	Description string `json:"description"`
	Poster      string `json:"poster"`
	SeasonID    uint64 `json:"season_id"`
}
