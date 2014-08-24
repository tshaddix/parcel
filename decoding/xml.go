package decoding

import (
	"encoding/xml"
	"net/http"

	"github.com/tshaddix/go-parcel"
)

type (
	XmlDecoder struct {
	}
)

func Xml() *XmlDecoder {
	return new(XmlDecoder)
}

func (self *XmlDecoder) Decode(r *http.Request, candidate parcel.Candidate) (err error) {
	ct := r.Header.Get("Content-Type")

	if ct == parcel.MimeXml || ct == parcel.MimeXml2 {
		err = xml.NewDecoder(r.Body).Decode(candidate)
		err = &RequestDecodeError{"", "", err}
	}

	return
}
