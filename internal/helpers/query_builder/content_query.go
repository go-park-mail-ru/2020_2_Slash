package query_builder

import (
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func BuildValuesQuery(valInd, valCount int) string {
	var values []string
	for i := 0; i < valCount; i++ {
		curInd := fmt.Sprintf("$%d", valInd+i)
		values = append(values, curInd)
	}
	valuesQuery := strings.Join(values, ", ")
	return fmt.Sprintf("(%s)", valuesQuery)
}

func BuildFilterQuery(entity string, valInd, valCount int) string {
	entityTable := fmt.Sprintf("content_%s", entity)             // content_genre
	entityID := fmt.Sprintf("%s.%s_id", entityTable, entity)     // content_genre.genre_id
	entityContentID := fmt.Sprintf("%s.content_id", entityTable) // content_genre.content_id

	valuesQuery := BuildValuesQuery(valInd, valCount)
	selectQuery := `
		SELECT %s
		FROM %s
		WHERE %s IN %s
		GROUP BY %s
		HAVING COUNT(%s)=%d`

	subQuery := fmt.Sprintf(
		selectQuery,
		entityContentID,
		entityTable,
		entityID,
		valuesQuery,
		entityContentID,
		entityContentID,
		valCount,
	)
	return fmt.Sprintf("AND c.id IN (%s)", subQuery)
}

func BuildYearCondition(valInd, valCount int) string {
	var conditions []string
	for i := 0; i < valCount; i++ {
		curCondition := fmt.Sprintf("c.year=$%d", valInd+i)
		conditions = append(conditions, curCondition)
	}
	resultCondition := strings.Join(conditions, " OR ")
	return fmt.Sprintf("AND (%s)", resultCondition)
}

func GetContentJoinFiltersByParams(values []interface{}, params *models.ContentFilter) (string, []interface{}) {
	var filters []string

	if params.Year != nil {
		filter := BuildYearCondition(len(values)+1, len(params.Year))
		filters = append(filters, filter)
		for _, year := range params.Year {
			values = append(values, year)
		}
	}

	if params.Genre != nil {
		filter := BuildFilterQuery("genre", len(values)+1, len(params.Genre))
		filters = append(filters, filter)
		for _, genre := range params.Genre {
			values = append(values, genre)
		}
	}

	if params.Country != nil {
		filter := BuildFilterQuery("country", len(values)+1, len(params.Country))
		filters = append(filters, filter)
		for _, country := range params.Country {
			values = append(values, country)
		}
	}

	if params.Actor != nil {
		filter := BuildFilterQuery("actor", len(values)+1, len(params.Actor))
		filters = append(filters, filter)
		for _, actor := range params.Actor {
			values = append(values, actor)
		}
	}

	if params.Director != nil {
		filter := BuildFilterQuery("director", len(values)+1, len(params.Director))
		filters = append(filters, filter)
		for _, director := range params.Director {
			values = append(values, director)
		}
	}

	filtersQuery := strings.Join(filters, " ")
	return filtersQuery, values
}
