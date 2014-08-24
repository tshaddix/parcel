package parcel

import (
	"net/http"
)

const (
	MimeJson = "application/json"
	MimeXml  = "application/xml"
	MimeXml2 = "text/xml"
)

type (
	Decoder interface {
		Decode(*http.Request, interface{}) error
	}

	Encoder interface {
		Encode(http.ResponseWriter, *http.Request, interface{}) (bool, error)
	}

	Parcel struct {
		rw      http.ResponseWriter
		r       *http.Request
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
	f.encoders = make([]Encoder)
	f.decoders = make([]Decoder)
	return f
}

// Encoder registers an encoder with the parcel factory
func (self *Factory) Encoder(encoder Encoder) {
	self.encoders = append(self.encoders, encoder)
}

// Decoder registers a decoder with the parcel factory
func (self *Factory) Decoder(decoder decoder) {
	self.decoders = append(self.decoders, decoder)
}

// Parcel creates a parcel for the given http context
func (self *Factory) Parcel(rw http.ResponseWriter, r *http.Request) *Parcel {
	return &Parcel{rw: rw, r: r, factory: self}
}

// Parcel

// Encode encodes candidate by passing through registered encoders on
// parent factory. Encoding will cease as soon as an encoder has responded
// with a `written` result of `true`.
func (self *Parcel) Encode(code int, c Candidate) (err error) {
	self.rw.WriteHeader(code)

	var written bool

	for _, encoder := range self.factory.encoders {
		if written, err = self.Encode(self.rw, self.r, c); err != nil || written == true {
			return
		}
	}
}

// Decode decodes candidate by passing through registered decoders on
// parent factory.
func (self *Parcel) Decode(c Candidate) (err error) {
	for _, decoder := range self.factory.decoders {
		if err = decoder.Decode(self.r, c); err != nil {
			return
		}
	}
}
