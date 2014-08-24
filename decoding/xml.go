package decoding

type (
	XmlDecoder struct {
	}
)

func Xml() *XmlDecoder {
	return new(XmlDecoder)
}

func (self *XmlDecoder) Decode(r *http.Request, candidate parcel.Candidate) (err error) {
	if r.Header.Get("Content-Type") == MimeXml {
		err = xml.NewDecoder(r.Body).Decode(candidate)

		if ute, ok := err.(*xml.UnmarshalTypeError); ok {
			err = &RequestDecodeError{
				FromType: ute.Value,
				ToType:   ute.Type.Name(),
			}
		}
	}

	return
}
