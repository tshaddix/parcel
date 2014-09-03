package encoding

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type (
	// QueryCodec is a Query implementation
	// of parcel.Decoder
	QueryCodec struct{}

	// QueryTypeError is used to indicate
	// an error during the decoding process of
	// a Query
	QueryTypeError struct {
		FromType string
		ToType   string
	}
)

// Query returns a new QueryCodec
func Query() *QueryCodec {
	return new(QueryCodec)
}

// set converts a string value into a target value's kind
func (q *QueryCodec) set(kind reflect.Kind, qs string, field reflect.Value) error {
	switch kind {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v, e := strconv.ParseInt(qs, 10, 64)

		if e != nil {
			return &QueryTypeError{FromType: "string", ToType: "integer"}
		}

		field.SetInt(v)
	case reflect.Float32, reflect.Float64:
		v, e := strconv.ParseFloat(qs, 10)

		if e != nil {
			return &QueryTypeError{FromType: "string", ToType: "floating point"}
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
		return &QueryTypeError{FromType: "string", ToType: field.Type().Name()}
	}

	return nil
}

// Decode will convert query string values into appropriate type
// values
func (q *QueryCodec) Decode(r *http.Request, candidate interface{}) error {
	queries := r.URL.Query()

	// Shortcut for queries strings
	if len(queries) == 0 {
		return nil
	}

	// Value and type of candidate
	value := reflect.ValueOf(candidate).Elem()
	ty := value.Type()

	// Iterate candidate fields
	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		tag := field.Tag.Get("query")

		// tag does not exists
		if tag == "" {
			continue
		}

		qs := queries.Get(tag)

		// stringer value does not exist
		if qs == "" {
			continue
		}

		// value can not be set
		if !value.Field(i).CanSet() {
			continue
		}

		// Do conversion

		kind := field.Type.Kind()

		switch kind {
		case reflect.Slice:
			arr := queries[tag]
			arrLen := len(arr)

			sl := reflect.MakeSlice(field.Type, arrLen, arrLen)

			for a, entry := range arr {
				if err := q.set(field.Type.Elem().Kind(), entry, sl.Index(a)); err != nil {
					return err
				}
			}

			value.Field(i).Set(sl)
		default:
			if err := q.set(kind, qs, value.Field(i)); err != nil {
				return err
			}
		}
	}

	return nil
}

// Error provides the error implementation for
// a QueryTypeError
func (e *QueryTypeError) Error() string {
	return "QueryTypeError: Can not convert type " + e.FromType + " to type " + e.ToType
}
