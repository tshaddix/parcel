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
		// the decoders are registered
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