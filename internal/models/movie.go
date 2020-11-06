package models

type Movie struct {
	ID    uint64 `json:"id"`
	Video string `json:"video"`
	Content
}
