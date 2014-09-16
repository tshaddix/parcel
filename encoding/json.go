package encoding

import (
	"encoding/json"
	"mime"
	"net/http"
	"strings"
)

type (
	// JSONCodec is a JSON implementation
	// of parcel.Decoder and parcel.Encoder
	JSONCodec struct {
		indent string
	}
)

// JSON returns a new JSON Encoder/Decoder
func JSON() *JSONCodec {
	return &JSONCodec{
		indent: "",
	}
}

// JSONIndent returns a new JSON Encoder/Decoder
// with the marshalled JSON indented by amt
func JSONIndent(amt int) *JSONCodec {
	return &JSONCodec{
		indent: strings.Repeat(" ", amt),
	}
}

// Decode simple wraps "encoding/json" decoder
// implementation by processing any request with
// content-type set to "application/json"
func (*JSONCodec) Decode(r *http.Request, candidate interface{}) (err error) {

	if !(r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") {
		return
	}

	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		return
	}

	if mt == MimeJSON {
		err = json.NewDecoder(r.Body).Decode(candidate)
	}

	return
}

// Encode will encode the candidate as a JSON response given
// the request content-type is set to "application/json"
func (jc *JSONCodec) Encode(rw http.ResponseWriter, candidate interface{}) (err error) {
	var output []byte

	if jc.indent != "" {
		output, err = json.MarshalIndent(candidate, "", jc.indent)
	} else {
		output, err = json.Marshal(candidate)
	}

	if err != nil {
		return
	}

	_, err = rw.Write(output)

	return
}

func (*JSONCodec) ContentType() string {
	return MimeJSON
}

func (*JSONCodec) Encodes() []string {
	return []string{MimeJSON}
}
