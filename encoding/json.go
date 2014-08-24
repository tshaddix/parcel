import (
	"encoding/json"
	"net/http"

	"github.com/tshaddix/go-parcel"
)

type (
	JsonEncoder struct{}
)

func Json() *JsonEncoder {
	return new(JsonEncoder)
}

func (self *JsonEncoder) Encode(rw http.ResponseWriter, r *http.Request, candidate parcel.Candidate) (written bool, err error) {
	if r.Header.Get("Content-Type") == parcel.MimeJson {
		written = true

		rw.Header().Set("Content-Type", parcel.MimeJson)
		encoder := json.NewEncoder(rw)
		err = encoder.Encode(candidate)
	}

	return
}