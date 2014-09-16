// Package parcel provides encoding and
// decoding of candidate values from http
// sources
package parcel

import (
	"errors"
	"mime"
	"net/http"
	"strings"
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
	// to a ResponseWriter.
	Encoder interface {
		Encodes() []string
		ContentType() string
		Encode(http.ResponseWriter, interface{}) error
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
		defaultEncoder Encoder
		encoders       map[string]Encoder
		decoders       []Decoder
	}
)

var (
	// Error returned when no response encoding is written to an http.ResponseWriter
	ResponseNotWrittenError = errors.New("Response was not written: No encoder indicated a written response")
)

// NewFactory creates a new Parcel factory
func NewFactory() *Factory {
	f := new(Factory)
	f.encoders = make(map[string]Encoder)
	f.decoders = make([]Decoder, 0)
	return f
}

// UseDefaultEncoder will set an encoder as a fallback if no Accept header is set
func (f *Factory) UseDefaultEncoder(encoder Encoder) {
	f.defaultEncoder = encoder
}

// UseEncoder registers an encoder with the parcel factory
func (f *Factory) UseEncoder(encoder Encoder) {
	for _, encodes := range encoder.Encodes() {
		f.encoders[encodes] = encoder
	}
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
func (p *Parcel) Encode(code int, c Candidate) error {
	encoder := p.GetEncoder()

	if encoder != nil {
		p.RW.Header().Set("Content-Type", encoder.ContentType())
		p.RW.WriteHeader(code)
		return encoder.Encode(p.RW, c)
	} else {
		return ResponseNotWrittenError
	}
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

// GetEncoder will return the encoder that will be used for encoding
// based on parcel's content negotiation rules
func (p *Parcel) GetEncoder() Encoder {
	// Content negotiation
	accepts := parseAccept(p.R.Header.Get("Accept"))

	if len(accepts) == 0 {
		// Try content-type retrieval
		if p.R.Method == "POST" || p.R.Method == "PUT" || p.R.Method == "PATCH" {
			mt, _, err := mime.ParseMediaType(p.R.Header.Get("Content-Type"))

			if err == nil {
				accepts = append(accepts, mt)
			}
		}
	}

	for _, a := range accepts {
		accepted := p.factory.encoders[a]

		if accepted != nil {
			return accepted
		}
	}

	if p.factory.defaultEncoder != nil {
		return p.factory.defaultEncoder
	} else {
		return nil
	}
}

// parseAccept will parce and accept header and return the accepted mimetypes
func parseAccept(accept string) []string {
	parts := strings.Split(accept, ",")
	for i, part := range parts {
		index := strings.IndexByte(part, ';')
		if index >= 0 {
			part = part[0:index]
		}
		part = strings.TrimSpace(part)
		parts[i] = part
	}
	return parts
}
