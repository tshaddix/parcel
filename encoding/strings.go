package encoding

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type (
	// Stringer implementations should
	// provide a way to fetch the length
	// of the string source and get a string
	// by a key
	Stringer interface {
		Get(*http.Request, string) string
		Len(*http.Request) int
	}

	// StringsCodec is a parcel.Decoder
	// implementation that decodes from a
	// string source
	StringsCodec struct {
		stringer Stringer

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

// Strings returns a new configured Strings
func Strings(stringer Stringer, name string) *StringsCodec {
	return &StringsCodec{stringer, name}
}

// Error provides the error implementation for
// a StringsDecodeError
func (self *StringsDecodeError) Error() string {
	return "StringsDecodeError: Can not convert type " + self.FromType + " to type " + self.ToType
}

// Decode uses a Stringer to match keys to
// candidate tags and convert values to the
// appropriate type
func (self *StringsCodec) Decode(r *http.Request, candidate interface{}) (err error) {
	// Shortcut for no strings
	if self.stringer.Len(r) == 0 {
		return
	}

	// Value and type of candidate
	value := reflect.ValueOf(candidate).Elem()
	ty := value.Type()

	// Iterate candidate fields
	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		tag := field.Tag.Get(self.tagName)

		// tag exists
		if tag != "" {
			qs := self.stringer.Get(r, tag)

			// string exists
			if qs != "" {
				if value.Field(i).CanSet() {

					// Do conversion

					kind := field.Type.Kind()

					if kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 || kind == reflect.Int {
						v, e := strconv.ParseInt(qs, 10, 64)

						if e != nil {
							err = &StringsDecodeError{FromType: "string", ToType: "integer"}
							return
						}

						value.Field(i).SetInt(v)
					} else if kind == reflect.Float32 || kind == reflect.Float64 {
						v, e := strconv.ParseFloat(qs, 10)

						if e != nil {
							err = &StringsDecodeError{FromType: "string", ToType: "floating point"}
							return
						}

						value.Field(i).SetFloat(v)
					} else if kind == reflect.String {
						value.Field(i).SetString(qs)
					} else if kind == reflect.Bool {
						v := strings.ToLower(qs)

						if v == "true" {
							value.Field(i).SetBool(true)
						} else {
							value.Field(i).SetBool(false)
						}
					} else {
						err = &StringsDecodeError{FromType: "string", ToType: field.Type.Name()}
						return
					}
				}
			}
		}
	}

	return
}
