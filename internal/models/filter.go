package models

type ContentFilter struct {
	Year     int   `query:"year"`
	Genre    int   `query:"genre"`
	Country  int   `query:"country"`
	Actor    int   `query:"actor"`
	Director int   `query:"director"`
	IsFree   *bool `query:"is_free"`
}
