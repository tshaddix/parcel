go-parcel
=========

Encoding and Decoding service for Golang web apps. Out of the box support for JSON, XML, and Querystrings.

```go
package main

import (
	"net/http"

	"github.com/tshaddix/go-parcel"
	"github.com/tshaddix/go-parcel/encoding"
	"github.com/tshaddix/go-parcel/decoding"
)

func main() {
	factory := parcel.NewFactory()

	// Encoders will be called in the
	// order they are registered
	factory.Encoder(encoding.Json())
	factory.Encoder(encoding.Xml())

	// Decoders will be called in the order
	// they are registered
	factory.Decoder(decoding.Query())
	factory.Decoder(decoding.Json())
	factory.Decoder(decoding.Xml())

	myHandler := func(rw http.ResponseWriter, r *http.Request){
		p := factory.Parcel(rw, r)

		// Decode
		err := p.Decode(&someStruct)

		// Decode will now be populated with
		// matching querystrings and json/xml
		// values from request body

		// Do stuff

		// Encode in appropriate response format
		// (xml to xml, json to json)
		err = p.Encode(http.StatusCreated, &someStruct)
	}

	// Build web server
}
```

## BYOD/E

Bring your own decoder/encoder. Anything with the function `Decode(*http.Request, parcel.Candidate) error` is a viable decoder.

Here is an example using [gorilla mux](https://github.com/gorilla/mux) params and the helpful `StringsDecoder` (used internally in `QueryDecoder`):

```go

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tshaddix/go-parcel"
	"github.com/tshaddix/go-parcel/decoding"
)

// decoding.Stringer implementation
type MuxParamStringer struct {}

// MuxParams is a shortcut for building a param
// decoder off of the strings decoder 
func MuxParams() *decoding.StringsDecoder {
	return &decoding.StringsDecoder{
		new(MuxParamStringer),
		"param", // process fields in format `param:"name"`
	}
}

// Len returns length of strings source
func (self *MuxParamStringer) Len(r *http.Request) int {
	return len(mux.Vars(r))
}

// Get returns the string value of a named parameter
func (self *MuxParamStringer) Get(r *http.Request, name string) string {
	return mux.Vars(r)[name]
}

// Later on

factory.Decoder(MuxParams())

```

Encoders have the function: `Encode(http.ResponseWriter, *http.Request, parcel.Candidate) (true, error)` the boolean result should indicate whether the encoding process wrote to the ResponseWriter. A `true` result will stop going down the line of remaining encoders.

## TODO
- Add more test cases