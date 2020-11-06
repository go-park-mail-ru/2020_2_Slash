package models

type Content struct {
	ContentID        uint64      `json:"content_id"`
	Name             string      `json:"name"`
	OriginalName     string      `json:"original_name"`
	Description      string      `json:"description"`
	ShortDescription string      `json:"short_description"`
	Year             int         `json:"year"`
	Images           string      `json:"images"`
	Type             string      `json:"type"`
	Countries        []*Country  `json:"counties"`
	Genres           []*Genre    `json:"genres"`
	Actors           []*Actor    `json:"actors"`
	Directors        []*Director `json:"directors"`
}

func (c *Content) ReplaceBy(other *Content) {
	if other.Name != "" {
		c.Name = other.Name
	}
	if other.OriginalName != "" {
		c.OriginalName = other.OriginalName
	}
	if other.Description != "" {
		c.Description = other.Description
	}
	if other.ShortDescription != "" {
		c.ShortDescription = other.ShortDescription
	}
	if other.Year != 0 {
		c.Year = other.Year
	}
	if other.Images != "" {
		c.Images = other.Images
	}
	if other.Type == "movie" || other.Type == "tv_show" {
		c.Type = other.Type
	}
	if len(other.Countries) > 0 {
		c.Countries = other.Countries
	}
	if len(other.Genres) > 0 {
		c.Genres = other.Genres
	}
	if len(other.Actors) > 0 {
		c.Actors = other.Actors
	}
	if len(other.Directors) > 0 {
		c.Directors = other.Directors
	}
}
