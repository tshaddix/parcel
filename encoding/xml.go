package encoding

import (
	"encoding/xml"
	"mime"
	"net/http"
)

type (
	// XMLCodec is a XML implementation
	// of parcel.Encoder and parcel.Decoder
	XMLCodec struct{}
)

// XML returns a new XMLCodec
func XML() *XMLCodec {
	return new(XMLCodec)
}

// Encode will encode the candidate as a XML response
// given the request content-type is set to "application/xml"
// or "text/xml"
func (*XMLCodec) Encode(rw http.ResponseWriter, r *http.Request, candidate interface{}, code int) (written bool, err error) {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		return
	}

	if mt != MimeXML && mt != MimeXML2 {
		return
	}

	written = true

	rw.Header().Set("Content-Type", mt)
	rw.WriteHeader(code)

	encoder := xml.NewEncoder(rw)
	err = encoder.Encode(candidate)

	return
}

// Decode simply wraps "encoding/xml" decoder
// implementation by processing any request with
// content-type set to "application/xml" or "text/xml"
func (*XMLCodec) Decode(r *http.Request, candidate interface{}) (err error) {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		return
	}

	if (mt != MimeXML && mt != MimeXML2) || r.ContentLength == 0 {
		return
	}

	err = xml.NewDecoder(r.Body).Decode(candidate)

	return
}
