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

// Decode simply wraps "encoding/xml" decoder
// implementation by processing any request with
// content-type set to "application/xml" or "text/xml"
func (*XMLCodec) Decode(r *http.Request, candidate interface{}) (err error) {
	if !(r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") {
		return
	}

	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		return
	}

	if mt != MimeXML && mt != MimeXML2 {
		return
	}

	err = xml.NewDecoder(r.Body).Decode(candidate)

	return
}

// Encode will encode the candidate as a XML response
// given the request content-type is set to "application/xml"
// or "text/xml"
func (*XMLCodec) Encode(rw http.ResponseWriter, candidate interface{}) (err error) {

	encoder := xml.NewEncoder(rw)
	err = encoder.Encode(candidate)

	return
}

func (*XMLCodec) ContentType() string {
	return MimeXML
}

func (*XMLCodec) Encodes() []string {
	return []string{MimeXML, MimeXML2}
}
