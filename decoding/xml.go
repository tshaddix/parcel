package decoding

import (
	"encoding/xml"
	"net/http"
)

const (
	MimeXml  = "application/xml"
	MimeXml2 = "text/xml"
)

type (
	XmlDecoder struct {
	}
)

func Xml() *XmlDecoder {
	return new(XmlDecoder)
}

func (self *XmlDecoder) Decode(r *http.Request, candidate interface{}) (err error) {
	ct := r.Header.Get("Content-Type")

	if ct == MimeXml || ct == MimeXml2 {
		err = xml.NewDecoder(r.Body).Decode(candidate)
	}

	return
}
