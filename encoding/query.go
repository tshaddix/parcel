package encoding

import (
	"net/http"
	"reflect"
)

type (
	// QueryCodec is a Query implementation
	// of parcel.Decoder
	QueryCodec struct{}

	// QueryTypeError is used to indicate
	// an error during the decoding process of
	// a Query
	QueryTypeError struct {
		ToType string
	}
)

// Query returns a new QueryCodec
func Query() *QueryCodec {
	return new(QueryCodec)
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
				if err := StrSet(field.Type.Elem().Kind(), entry, sl.Index(a)); err != nil {
					return &QueryTypeError{field.Type.Elem().Name()}
				}
			}

			value.Field(i).Set(sl)
		default:
			if err := StrSet(kind, qs, value.Field(i)); err != nil {
				return &QueryTypeError{field.Type.Name()}
			}
		}
	}

	return nil
}

// Error provides the error implementation for
// a QueryTypeError
func (e *QueryTypeError) Error() string {
	return "QueryTypeError: Can not convert type query string to type " + e.ToType
}
