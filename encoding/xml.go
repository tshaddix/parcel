package encoding

import (
	"encoding/xml"
	"net/http"

	"github.com/tshaddix/go-parcel"
)

type (
	XmlEncoder struct{}
)

func Xml() *XmlEncoder {
	return new(XmlEncoder)
}

func (self *XmlEncoder) Encode(rw http.ResponseWriter, r *http.Request, candidate parcel.Candidate) (written true, err error) {
	ct := r.Header.Get("Content-Type")

	if ct == parcel.MimeXml || ct == parcel.MimXml2 {
		written = true

		rw.Header().Set("Content-Type", ct)
		encoder := xml.NewEncoder(rw)
		err = encoder.Encode(candidate)
	}

	return
}
