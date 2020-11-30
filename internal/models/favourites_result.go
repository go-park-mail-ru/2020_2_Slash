package models

type FavouritesResult struct {
	Movies  []*Movie  `json:"movies"`
	TVShows []*TVShow `json:"tv_shows"`
}
