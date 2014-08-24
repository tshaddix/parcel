package decoding

import (
	"net/http"

	"github.com/tshaddix/go-parcel"
)

type (
	QueryDecoder struct {
	}
)

func (self *QueryDecoder) Decode(r *http.Request, candidate parcel.Candidate) (err error) {

}
