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
		Decode(*http.Request, Candidate) error
	}

	// Encoder implementations should encode values from a candidate
	// to a ResponseWriter. A bool should be returned to indicate whether
	// the encoder wrote a response
	Encoder interface {
		Encode(http.ResponseWriter, *http.Request, Candidate) (bool, error)
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

// Encoder registers an encoder with the parcel factory
func (self *Factory) Encoder(encoder Encoder) {
	self.encoders = append(self.encoders, encoder)
}

// Decoder registers a decoder with the parcel factory
func (self *Factory) Decoder(decoder Decoder) {
	self.decoders = append(self.decoders, decoder)
}

// Parcel creates a parcel for the given http context
func (self *Factory) Parcel(rw http.ResponseWriter, r *http.Request) *Parcel {
	return &Parcel{RW: rw, R: r, factory: self}
}

// Parcel

// Encode encodes candidate by passing through registered encoders on
// parent factory. Encoding will cease as soon as an encoder has responded
// with a `written` result of `true`. If no encoders write to the response,
// an ResponseNotWrittenError is returned.
func (self *Parcel) Encode(code int, c Candidate) (err error) {
	self.RW.WriteHeader(code)

	var written bool

	for _, encoder := range self.factory.encoders {
		if written, err = encoder.Encode(self.RW, self.R, c); err != nil || written == true {
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
func (self *Parcel) Decode(c Candidate) (err error) {
	for _, decoder := range self.factory.decoders {
		if err = decoder.Decode(self.R, c); err != nil {
			return
		}
	}

	return
}
