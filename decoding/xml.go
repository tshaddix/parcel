package decoding

import (
	"encoding/xml"
	"net/http"
)

type (
	// XmlDecoder providers an XML implementation
	// of parcel.Decoder
	XmlDecoder struct {
	}
)

// Xml returns a new XmlDecoder
func Xml() *XmlDecoder {
	return new(XmlDecoder)
}

// Decode simply wraps "encoding/xml" decoder
// implementation by processing any request with
// content-type set to "application/xml" or "text/xml"
func (self *XmlDecoder) Decode(r *http.Request, candidate interface{}) (err error) {
	ct := r.Header.Get("Content-Type")

	if ct == MimeXml || ct == MimeXml2 {
		err = xml.NewDecoder(r.Body).Decode(candidate)
	}

	return
}
