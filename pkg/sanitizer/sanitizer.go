package sanitizer

import (
	"github.com/microcosm-cc/bluemonday"
	"reflect"
)

func Sanitize(i interface{}) {
	elem := reflect.ValueOf(i).Elem()
	policy := bluemonday.UGCPolicy()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.Kind() == reflect.String {
			sanitized := policy.Sanitize(field.String())
			field.SetString(sanitized)
		}
	}
}
