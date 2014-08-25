package parcel

import (
	"net/http"
)

type (
	Decoder interface {
		Decode(*http.Request, interface{}) error
	}

	Encoder interface {
		Encode(http.ResponseWriter, *http.Request, interface{}) (bool, error)
	}

	Parcel struct {
		RW      http.ResponseWriter
		R       *http.Request
		factory *Factory
	}

	Factory struct {
		encoders []Encoder
		decoders []Decoder
	}

	Candidate interface{}
)

// Factory

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
// with a `written` result of `true`.
func (self *Parcel) Encode(code int, c Candidate) (err error) {
	self.RW.WriteHeader(code)

	var written bool

	for _, encoder := range self.factory.encoders {
		if written, err = encoder.Encode(self.RW, self.R, c); err != nil || written == true {
			return
		}
	}

	return
}

// Decode decodes candidate by passing through registered decoders on
// parent factory.
func (self *Parcel) Decode(c Candidate) (err error) {
	for _, decoder := range self.factory.decoders {
		if err = decoder.Decode(self.R, c); err != nil {
			return
		}
	}

	return
}
