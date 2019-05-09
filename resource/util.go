package resource

import (
	"gopkg.in/bluesuncorp/validator.v5"
	"reflect"
)

func hasExistValidateError(err error, field string) bool {
	if value, ok := err.(*validator.StructErrors); ok {
		if tagErr := value.Errors[field]; tagErr != nil {
			return true
		}

	}
	return false
}

func getJsonTagValue(t interface{}, field string) (string, bool) {
	structType := reflect.TypeOf(t)
	stField, exist := structType.FieldByName(field)
	if exist {
		tagValue := stField.Tag.Get("json")
		if tagValue != "" {
			return tagValue, true
		}
	}
	return "", false
}
