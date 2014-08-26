package encoding

import (
	"net/http"
)

type (
	// QueryStringer is an implementation
	// strings.Stringer for query strings
	QueryStringer struct{}
)

// Len returns the length of the query strings
// in the request URL
func (self *QueryStringer) Len(r *http.Request) int {
	return len(r.URL.Query())
}

// Get is a simple pass through to Get for
// request URL query strings
func (self *QueryStringer) Get(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

// QueryDecode returns a StringsCodec configured
// with a QueryStringer and tag-name "query"
func Query() *StringsCodec {
	return Strings(new(QueryStringer), "query")
}
