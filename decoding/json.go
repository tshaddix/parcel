package decoding

import (
	"encoding/json"
	"net/http"
)

const (
	MimeJson = "application/json"
)

type (
	JsonDecoder struct{}
)

func Json() *JsonDecoder {
	return new(JsonDecoder)
}

func (self *JsonDecoder) Decode(r *http.Request, candidate interface{}) (err error) {

	if r.Header.Get("Content-Type") == MimeJson {
		err = json.NewDecoder(r.Body).Decode(candidate)
	}

	return
}
