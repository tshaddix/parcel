package encoding

import (
	"encoding/json"
	"net/http"
)

type (
	// JsonEncoder is a JSON implementation
	// of parcel.Encoder
	JsonEncoder struct{}
)

// Json returns a new JSON encoder
func Json() *JsonEncoder {
	return new(JsonEncoder)
}

// Encode will encode the candidate as a JSON response given
// the request content-type is set to "application/json"
func (self *JsonEncoder) Encode(rw http.ResponseWriter, r *http.Request, candidate interface{}) (written bool, err error) {
	if r.Header.Get("Content-Type") == MimeJson {
		written = true

		rw.Header().Set("Content-Type", MimeJson)
		encoder := json.NewEncoder(rw)
		err = encoder.Encode(candidate)
	}

	return
}
