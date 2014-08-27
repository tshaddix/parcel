package encoding

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type (
	// Stringer implementations should
	// provide a way to fetch source strings
	// as a map (key, value)
	Stringer interface {
		Len(*http.Request) int
		Get(*http.Request, string) string
	}

	// StringsCodec is a parcel.Decoder
	// implementation that decodes from a
	// string source
	StringsCodec struct {
		sr Stringer

		// tagName indicates the tag value
		// used to find matching keys
		// (e.g. "query")
		tagName string
	}

	// StringsDecodeError is used to indicate
	// an error during the decoding process of
	// Stringer
	StringsDecodeError struct {
		FromType string
		ToType   string
	}
)

// Strings returns a new configured StringsCodec
func Strings(s Stringer, name string) *StringsCodec {
	return &StringsCodec{s, name}
}

// Error provides the error implementation for
// a StringsDecodeError
func (e *StringsDecodeError) Error() string {
	return "StringsDecodeError: Can not convert type " + e.FromType + " to type " + e.ToType
}

// Decode uses a StringMapper to match keys to
// candidate tags and convert values to the
// appropriate type
func (s *StringsCodec) Decode(r *http.Request, candidate interface{}) error {
	// Shortcut for no strings
	if s.sr.Len(r) == 0 {
		return nil
	}

	// Value and type of candidate
	value := reflect.ValueOf(candidate).Elem()
	ty := value.Type()

	// Iterate candidate fields
	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		tag := field.Tag.Get(s.tagName)

		// tag does not exists
		if tag == "" {
			continue
		}

		qs := s.sr.Get(r, tag)

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
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			v, e := strconv.ParseInt(qs, 10, 64)

			if e != nil {
				return &StringsDecodeError{FromType: "string", ToType: "integer"}
			}

			value.Field(i).SetInt(v)
		case reflect.Float32, reflect.Float64:
			v, e := strconv.ParseFloat(qs, 10)

			if e != nil {
				return &StringsDecodeError{FromType: "string", ToType: "floating point"}
			}

			value.Field(i).SetFloat(v)
		case reflect.String:
			value.Field(i).SetString(qs)
		case reflect.Bool:
			v := strings.ToLower(qs)

			if v == "true" {
				value.Field(i).SetBool(true)
			} else {
				value.Field(i).SetBool(false)
			}
		default:
			return &StringsDecodeError{FromType: "string", ToType: field.Type.Name()}
		}
	}

	return nil
}
