package encoding

import (
	"encoding/xml"
	"net/http"
)

type (
	// XmlEncoder is a XML implementation
	// of parcel.Encoder and parcel.Decoder
	XmlCodec struct{}
)

// Xml returns a new XmlCodec
func Xml() *XmlCodec {
	return new(XmlCodec)
}

// Encode will encode the candidate as a XML response
// given the request content-type is set to "application/xml"
// or "text/xml"
func (_ *XmlCodec) Encode(rw http.ResponseWriter, r *http.Request, candidate interface{}, code int) (written bool, err error) {
	ct := r.Header.Get("Content-Type")

	if ct == MimeXml || ct == MimeXml2 {
		written = true

		rw.Header().Set("Content-Type", ct)
		rw.WriteHeader(code)

		encoder := xml.NewEncoder(rw)
		err = encoder.Encode(candidate)
	}

	return
}

// Decode simply wraps "encoding/xml" decoder
// implementation by processing any request with
// content-type set to "application/xml" or "text/xml"
func (_ *XmlCodec) Decode(r *http.Request, candidate interface{}) (err error) {
	ct := r.Header.Get("Content-Type")

	if (ct == MimeXml || ct == MimeXml2) && r.ContentLength != 0 {
		err = xml.NewDecoder(r.Body).Decode(candidate)
	}

	return
}
