package encoding

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var TypeMatchError = errors.New("No matching type")

// StrSet is a utility function for setting a typed value
// by converting a string
func StrSet(kind reflect.Kind, qs string, field reflect.Value) error {
	switch kind {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v, e := strconv.ParseInt(qs, 10, 64)

		if e != nil {
			return e
		}

		field.SetInt(v)
	case reflect.Float32, reflect.Float64:
		v, e := strconv.ParseFloat(qs, 10)

		if e != nil {
			return e
		}

		field.SetFloat(v)
	case reflect.String:
		field.SetString(qs)
	case reflect.Bool:
		v := strings.ToLower(qs)

		if v == "true" {
			field.SetBool(true)
		} else {
			field.SetBool(false)
		}
	default:
		return TypeMatchError
	}

	return nil
}
