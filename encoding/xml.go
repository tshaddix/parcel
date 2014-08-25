package encoding

import (
	"encoding/xml"
	"net/http"
)

const (
	MimeXml  = "application/xml"
	MimeXml2 = "text/xml"
)

type (
	XmlEncoder struct{}
)

func Xml() *XmlEncoder {
	return new(XmlEncoder)
}

func (self *XmlEncoder) Encode(rw http.ResponseWriter, r *http.Request, candidate interface{}) (written bool, err error) {
	ct := r.Header.Get("Content-Type")

	if ct == MimeXml || ct == MimeXml2 {
		written = true

		rw.Header().Set("Content-Type", ct)
		encoder := xml.NewEncoder(rw)
		err = encoder.Encode(candidate)
	}

	return
}
