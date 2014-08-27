package encoding

import (
	"encoding/json"
	"mime"
	"net/http"
)

type (
	// JSONCodec is a JSON implementation
	// of parcel.Decoder and parcel.Encoder
	JSONCodec struct{}
)

// JSON returns a new Json Encoder/Decoder
func JSON() *JSONCodec {
	return new(JSONCodec)
}

// Decode simple wraps "encoding/json" decoder
// implementation by processing any request with
// content-type set to "application/json"
func (*JSONCodec) Decode(r *http.Request, candidate interface{}) (err error) {

	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		return
	}

	if mt == MimeJSON && r.ContentLength != 0 {
		err = json.NewDecoder(r.Body).Decode(candidate)
	}

	return
}

// Encode will encode the candidate as a JSON response given
// the request content-type is set to "application/json"
func (*JSONCodec) Encode(rw http.ResponseWriter, r *http.Request, candidate interface{}, code int) (written bool, err error) {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		return
	}

	if mt != MimeJSON {
		return
	}

	written = true

	rw.Header().Set("Content-Type", MimeJSON)
	rw.WriteHeader(code)

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(candidate)

	return
}
