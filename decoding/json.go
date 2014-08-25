package decoding

import (
	"encoding/json"
	"net/http"
)

type (
	// JsonDecoder is a JSON implementation
	// of parcel.Decoder
	JsonDecoder struct{}
)

// Json returns a new JsonDecoder
func Json() *JsonDecoder {
	return new(JsonDecoder)
}

// Decoder simple wraps "encoding/json" decoder
// implementation by processing any request with
// content-type set to "application/json"
func (self *JsonDecoder) Decode(r *http.Request, candidate interface{}) (err error) {

	if r.Header.Get("Content-Type") == MimeJson {
		err = json.NewDecoder(r.Body).Decode(candidate)
	}

	return
}
