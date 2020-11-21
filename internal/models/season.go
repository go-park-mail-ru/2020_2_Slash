package models

type Season struct {
	ID             uint64     `json:"id"`
	Number         int        `json:"number"`
	EpisodesNumber int        `json:"episodes_number"`
	TVShowID       uint64     `json:"tv_show_id"`
	Episodes       []*Episode `json:"episodes"`
}
