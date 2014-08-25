package encoding

import (
	"encoding/xml"
	"net/http"
)

type (
	// XmlEncoder is a XML implementation
	// of parcel.Encoder
	XmlEncoder struct{}

	// XmlDecoder providers an XML implementation
	// of parcel.Decoder
	XmlDecoder struct{}
)

// XmlEncode returns a new XmlEncoder
func XmlEncode() *XmlEncoder {
	return new(XmlEncoder)
}

// XmlDecode returns a new XmlDecoder
func XmlDecode() *XmlDecoder {
	return new(XmlDecoder)
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

// Decode simply wraps "encoding/xml" decoder
// implementation by processing any request with
// content-type set to "application/xml" or "text/xml"
func (self *XmlDecoder) Decode(r *http.Request, candidate interface{}) (err error) {
	ct := r.Header.Get("Content-Type")

	if (ct == MimeXml || ct == MimeXml2) && r.ContentLength > 0 {
		err = xml.NewDecoder(r.Body).Decode(candidate)
	}

	return
}
