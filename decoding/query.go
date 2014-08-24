package decoding

import (
	"net/http"
)

type (
	QueryStringer struct {
	}

	QueryDecoder struct {
		sd *StringsDecoder
	}
)

func Query() *QueryDecoder {
	sd := &StringsDecoder{
		new(QueryStringer),
		"query",
	}

	return &QueryDecoder{sd}
}

func (self *QueryStringer) Len(r *http.Request) int {
	return len(r.URL.Query())
}

func (self *QueryStringer) Get(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

func (self *QueryDecoder) Decode(r *http.Request, candidate parcel.Candidate) (err error) {
	return self.sd.Decode(r, candidate)
}
