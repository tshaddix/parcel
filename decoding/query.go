package decoding

import (
	"net/http"
)

type (
	QueryStringer struct {
	}
)

func Query() *StringsDecoder {
	sd := &StringsDecoder{
		new(QueryStringer),
		"query",
	}

	return sd
}

func (self *QueryStringer) Len(r *http.Request) int {
	return len(r.URL.Query())
}

func (self *QueryStringer) Get(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}
