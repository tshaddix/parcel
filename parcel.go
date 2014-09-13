// Package parcel provides encoding and
// decoding of candidate values from http
// sources
package parcel

import (
	"errors"
	"net/http"
)

type (
	// Candidate is a shortcut for candidate targets
	Candidate interface{}

	// Decoder implementations should decode values from a
	// request to a candidate
	Decoder interface {
		Decode(*http.Request, interface{}) error
	}

	// Encoder implementations should encode values from a candidate
	// to a ResponseWriter. A bool should be returned to indicate whether
	// the encoder wrote a response
	Encoder interface {
		Encode(http.ResponseWriter, *http.Request, interface{}, int) (bool, error)
	}

	// Parcel is a simple reference structure that
	// enables easy encoding and decoding
	Parcel struct {
		RW      http.ResponseWriter
		R       *http.Request
		factory *Factory
	}

	// Factory stores the implementation details of available
	// and configured encoders and decoders
	Factory struct {
		encoders []Encoder
		decoders []Decoder
	}
)

var (
	// Error returned when no response encoding is written to an http.ResponseWriter
	ResponseNotWrittenError = errors.New("Response was not written: No encoder indicated a written response")
)

// NewFactory creates a new Parcel factory
func NewFactory() *Factory {
	f := new(Factory)
	f.encoders = make([]Encoder, 0)
	f.decoders = make([]Decoder, 0)
	return f
}

// UseEncoder registers an encoder with the parcel factory
func (f *Factory) UseEncoder(encoder Encoder) {
	f.encoders = append(f.encoders, encoder)
}

// UseDecoder registers a decoder with the parcel factory
func (f *Factory) UseDecoder(decoder Decoder) {
	f.decoders = append(f.decoders, decoder)
}

// Use is a convience function to register encoders and
// decoders.
func (f *Factory) Use(i interface{}) {
	if decoder, ok := i.(Decoder); ok {
		f.UseDecoder(decoder)
	}

	if encoder, ok := i.(Encoder); ok {
		f.UseEncoder(encoder)
	}
}

// Parcel creates a parcel for the given http context
func (f *Factory) Parcel(rw http.ResponseWriter, r *http.Request) *Parcel {
	return &Parcel{RW: rw, R: r, factory: f}
}

// NewCodec is an alias to Parcel()
func (f *Factory) NewCodec(rw http.ResponseWriter, r *http.Request) *Parcel {
	return f.Parcel(rw, r)
}

// Parcel

// Encode encodes candidate by passing through registered encoders on
// parent factory. Encoding will cease as soon as an encoder has responded
// with a `written` result of `true`. If no encoders write to the response,
// an ResponseNotWrittenError is returned.
func (p *Parcel) Encode(code int, c Candidate) (err error) {
	var written bool

	for _, encoder := range p.factory.encoders {
		if written, err = encoder.Encode(p.RW, p.R, c, code); err != nil || written == true {
			return
		}
	}

	if !written {
		err = ResponseNotWrittenError
	}

	return
}

// Decode decodes candidate by passing through registered decoders on
// parent factory. If any decoder returns an error, the chain is stopped
// and the error is returned
func (p *Parcel) Decode(c Candidate) (err error) {
	for _, decoder := range p.factory.decoders {
		if err = decoder.Decode(p.R, c); err != nil {
			return
		}
	}

	return
}
