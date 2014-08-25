package decoding

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type (
	Stringer interface {
		Get(*http.Request, string) string
		Len(*http.Request) int
	}

	StringsDecoder struct {
		stringer Stringer
		tagName  string
	}

	StringsDecodeError struct {
		FromType string
		ToType   string
	}
)

func Strings(stringer Stringer, name string) *StringsDecoder {
	return &StringsDecoder{stringer, name}
}

func (self *StringsDecodeError) Error() string {
	return "StringsDecodeError: Can not convert type " + self.FromType + " to type " + self.ToType
}

func (self *StringsDecoder) Decode(r *http.Request, candidate interface{}) (err error) {
	// Shortcut for no strings
	if self.stringer.Len(r) == 0 {
		return
	}

	value := reflect.ValueOf(candidate).Elem()
	ty := value.Type()

	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		tag := field.Tag.Get(self.tagName)

		// tag exists
		if tag != "" {
			qs := self.stringer.Get(r, tag)

			// string exists
			if qs != "" {
				if value.Field(i).CanSet() {

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
