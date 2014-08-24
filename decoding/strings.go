package decoding

type (
	Stringer interface {
		Get(*http.Request, string) string
		Len(*http.Request) int
	}

	StringsDecoder struct {
		stringer Stringer
		tagName  string
	}
)

func (self *StringsDecoder) Decode(r *http.Request, candidate parcel.Candidate) (err error) {
	// Shortcut for no strings
	if self.stringer.Len() == 0 {
		return
	}

	ty := reflect.TypeOf(candidate).Elem()
	value := reflect.ValueOf(candidate).Elem()

	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		tag := field.Tag.Get(self.tagName)

		// tag exists
		if tag != "" {
			qs := self.stringer.Get(tag)

			// string exists
			if qs != "" {
				kind := field.Type.Kind()

				if kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 || kind == reflect.Int {
					v, e := strconv.ParseInt(qs, 10, 64)

					if e != nil {
						err = &RequestDecodeError{FromType: "string", ToType: "integer"}
						return
					}

					value.Field(i).SetInt(v)
				} else if kind == reflect.Float32 || kind == reflect.Float64 {
					v, e := strconv.ParseFloat(qs, 10, 64)

					if e != nil {
						err = &RequestDecodeError{FromType: "string", ToType: "floating point"}
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
					err = &RequestDecodeError{FromType: "string", ToType: field.Type.Name()}
					return
				}
			}
		}
	}
}
