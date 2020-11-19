package query_builder

import (
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func BuildContentJoinFilter(entity string, valInd int) string {
	entityTable := "content_" + entity                  // content_genre
	entityID := entityTable + "." + entity + "_id"      // content.genre_id
	entityContentID := entityTable + "." + "content_id" // content_genre.content_id

	// JOIN content_genre ON c.id=cg.content_id AND cg.genre_id=$1
	filter := "JOIN " + entityTable + " ON " + "c.id=" +
		entityContentID + " AND " + entityID + "=$" + strconv.Itoa(valInd)

	return filter
}

func GetContentJoinFiltersByParams(values []interface{}, params *models.ContentFilter) (string, []interface{}) {
	var filters []string

	if params.Genre != 0 {
		filter := BuildContentJoinFilter("genre", len(values)+1)
		filters = append(filters, filter)
		values = append(values, params.Genre)
	}

	if params.Country != 0 {
		filter := BuildContentJoinFilter("country", len(values)+1)
		filters = append(filters, filter)
		values = append(values, params.Country)
	}

	if params.Actor != 0 {
		filter := BuildContentJoinFilter("actor", len(values)+1)
		filters = append(filters, filter)
		values = append(values, params.Actor)
	}

	if params.Director != 0 {
		filter := BuildContentJoinFilter("director", len(values)+1)
		filters = append(filters, filter)
		values = append(values, params.Director)
	}

	filtersQuery := strings.Join(filters, " ")
	return filtersQuery, values
}

func GetContentWhereQueryByParams(values []interface{}, params *models.ContentFilter) (string, []interface{}) {
	var filters []string

	if params.Year != 0 {
		ind := len(values) + 1
		filter := `WHERE c.year=$` + strconv.Itoa(ind)
		filters = append(filters, filter)
		values = append(values, params.Year)
	}

	filtersQuery := strings.Join(filters, " ")
	return filtersQuery, values
}
