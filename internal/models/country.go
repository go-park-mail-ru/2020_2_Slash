package models

type Country struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type ContentCountry struct {
	ContentID uint64 `json:"content_id"`
	CountryID uint64 `json:"country_id"`
}
