package encoding

import (
	"encoding/json"
	"net/http"
)

type (
	// JsonEncoder is a JSON implementation
	// of parcel.Encoder
	JsonEncoder struct{}

	// JsonDecoder is a JSON implementation
	// of parcel.Decoder
	JsonDecoder struct{}
)

// JsonEncode returns a new JsonDecoder
func JsonDecode() *JsonDecoder {
	return new(JsonDecoder)
}

// JsonEncode returns a new JsonEncoder
func JsonEncode() *JsonEncoder {
	return new(JsonEncoder)
}

// Decode simple wraps "encoding/json" decoder
// implementation by processing any request with
// content-type set to "application/json"
func (self *JsonDecoder) Decode(r *http.Request, candidate interface{}) (err error) {

	if r.Header.Get("Content-Type") == MimeJson && r.ContentLength > 0 {
		err = json.NewDecoder(r.Body).Decode(candidate)
	}

	return
}

// Encode will encode the candidate as a JSON response given
// the request content-type is set to "application/json"
func (self *JsonEncoder) Encode(rw http.ResponseWriter, r *http.Request, candidate interface{}, code int) (written bool, err error) {
	if r.Header.Get("Content-Type") == MimeJson {
		written = true

		rw.Header().Set("Content-Type", MimeJson)
		rw.WriteHeader(code)

		encoder := json.NewEncoder(rw)
		err = encoder.Encode(candidate)
	}

	return
}
