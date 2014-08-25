package encoding

import (
	"encoding/json"
	"net/http"
)

const (
	MimeJson = "application/json"
)

type (
	JsonEncoder struct{}
)

func Json() *JsonEncoder {
	return new(JsonEncoder)
}

func (self *JsonEncoder) Encode(rw http.ResponseWriter, r *http.Request, candidate interface{}) (written bool, err error) {
	if r.Header.Get("Content-Type") == MimeJson {
		written = true

		rw.Header().Set("Content-Type", MimeJson)
		encoder := json.NewEncoder(rw)
		err = encoder.Encode(candidate)
	}

	return
}
