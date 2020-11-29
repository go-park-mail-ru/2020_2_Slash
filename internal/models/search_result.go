package models

type SearchResult struct {
	TVShows []*TVShow `json:"tv_shows"`
	Movies  []*Movie  `json:"movies"`
	Actors  []*Actor  `json:"actors"`
}
