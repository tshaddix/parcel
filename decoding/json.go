package decoding

import (
	"encoding/json"
	"net/http"

	"github.com/tshaddix/go-parcel"
)

type (
	JsonDecoder struct{}
)

func Json() *JsonDecoder {
	return new(JsonDecoder)
}

func (self *JsonDecoder) Decode(r *http.Request, candidate parcel.Candidate) (err error) {

	if r.Header.Get("Content-Type") == parcel.MimeJson {
		err = json.NewDecoder(r.Body).Decode(candidate)

		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			err = &RequestDecodeError{
				FromType: ute.Value,
				ToType:   ute.Type.Name(),
				Err:      err,
			}
		}
	}

	return
}
