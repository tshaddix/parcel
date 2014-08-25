package encoding

import (
	"encoding/xml"
	"net/http"
)

type (
	// XmlEncoder is a XML implementation
	// of parcel.Encoder
	XmlEncoder struct{}
)

// Xml returns a new XML encoder
func Xml() *XmlEncoder {
	return new(XmlEncoder)
}

// Encode will encode the candidate as a XML response
// given the request content-type is set to "application/xml"
// or "text/xml"
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
