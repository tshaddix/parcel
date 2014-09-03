parcel
======

Encoding and decoding ease for golang webapps. Out of the box support for JSON, XML, and Query Strings.

This package aims to provide a middleware inspired codec system for golang webapps. Using `parcel` involves creating a parcel factory and then using a parcel built from an `http.ResponseWriter` and `*http.Request`.

## Factory

A parcel factory stores the encoder and decoder chains used to process requests. That's it! Setting up a parcel factory goes something like this:

```go
// New factory
factory = parcel.NewFactory()

// Encoders/Decoders will be called in the order
// they are registered. The setup below:
// Request ->
// 1. Query Strings
// 2. Json
// 3. Xml
// Response ->
// 1. Json
// 2. Xml
// Notice that the Query codec only provides a decoder,
// so it will not be added to response chain
factory.Use(encoding.Query())
factory.Use(encoding.JSON())
factory.Use(encoding.XML())
```

The `Use()` function will register any decoders or encoders as part of the middleware chain. The middleware will run in the order registered when processing a request or writing a response.

Some of these codecs (such as `JsonCodec` and `XmlCodec`) will register both a decoder and an encoder when using `Use()`. To register just an encoder or a decoder, you can use `UseDecoder()` or `UseEncoder()`.

```go
// Just register the JSON encoder in the encoder chain
factory.UseEncoder(encoding.JSON())
```

## Parcel

After configuring a factory, you can use it to build a `*parcel.Parcel` for encoding and decoding.

```go
// factory configured before

func myHandler(rw http.ResponseWriter, r *http.Request) {
	p := factory.Parcel(rw, r)
}
```

### Decode(interface{}) error

`Decode` takes an interface, and runs all registered decoders on the factory, in order. It is important to note that if two decoders have a reference to the same property, the last decoder to run will be the final assignement of the property.

```go
// Run decoders and populate personStruct
err := p.Decode(&personStruct)
```

Errors during decoding are returned "raw" from whichever decoder returned the error. Errors will halt the decoding process immediately.

### Encode(int, interface{}) error

`Encode` will take an http status code and interface, and does the opposite of a decoder (it writes a response). The main difference between a decoder and an encoder is the handling of the next encoder in the chain. If an encoder indicates that it has written to the response, the chain is terminated. It is up to the encoder to determine whether it should write to the request or not. The included `XmlCodec` and `JsonCodec` use the request `Content-Type` header to determine the encoding.

```go
err := p.Encode(http.StatusOK, &personStruct)
```

Errors work the same as the decoding process: An error will be returned "raw" from the encoder and the processing chain will be halted.

## In the Box

Included with `parcel` are a few codecs which should be useful. You can pick and choose which codecs to use/extend. They can be found under `parcel/encoding`.

### JSONCodec

`JSONCodec` uses the `encoding/json` package to encode and/or decode JSON bodies. The `JSONCodec` uses the `Content-Type` header to determine whether to encode or decode a given parcel. Only request indicating `application/json` as the content type will be processed.

```go
jsonCodec := encoding.JSON()
```

### XMLCodec

`XMLCodec` acts much like the JSON version by wrapping the `encoding/xml` package. Requests indicating the content type of `application/xml` or `text/xml` will be the only requests processed by the codec.

```go
xmlCodec := encoding.XML()
```

### QueryCodec (a configured StringsCodec)

`QueryCodec` configures a decoder implementation that parses query strings from a request.

```go
queryCodec := encoding.Query()
```

Use the `query` tag to indicate which fields should be processed by the query codec on your structs. Fields without a query tag will be ignored by this codec.

```go
type myStruct struct {
	Token string  `query:"token"`
	Ids   []int64 `query:"id"` //?id=1&id=2&id=3
}
```

### StringsCodec
`StringsCodec` is a simple codec that uses reflection to manipulate a string source into a target destination property and type. Refer to the [godoc](https://godoc.org/github.com/tshaddix/parcel) for more information on how to use this codec (or look at the example below).

## Adding a Codec

Here is an example using [gorilla mux](https://github.com/gorilla/mux) params and the helpful `StringsCodec`:

```go

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tshaddix/parcel"
	"github.com/tshaddix/parcel/encoding"
)

// decoding.Stringer implementation
type MuxParamStringer struct {}

// MuxParams is a shortcut for building a param
// decoder off of the strings decoder 
func MuxParams() *encoding.StringsCodec {
	return encoding.Strings(
		new(MuxParamStringer),
		"param", // process fields in format `param:"name"`
	)
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

factory.Use(MuxParams())

```

For more advanced needs of building a codec, refer to the [godoc](https://godoc.org/github.com/tshaddix/parcel).

## TODO
- Comments on decoding/encoding packages
- Update README with better examples and abilities.
- Add more test cases (encoding, bad cases, etc)