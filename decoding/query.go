package decoding

import (
	"net/http"
)

type (
	QueryStringer struct{}
)

func (self *QueryStringer) Len(r *http.Request) int {
	return len(r.URL.Query())
}

func (self *QueryStringer) Get(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

func Query() *StringsDecoder {
	return &StringsDecoder{new(QueryStringer), "query"}
}
