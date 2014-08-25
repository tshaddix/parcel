package decoding

import (
	"net/http"
	"os"
	"testing"
)

type (
	TestXml struct {
		Name       string  `xml:"name"`
		Email      string  `xml:"email"`
		IsAdmin    bool    `xml:"is-admin"`
		Age        int     `xml:"age"`
		HourlyRate float32 `xml:"hourly-rate"`
	}
)

func TestXmlDecoder(t *testing.T) {
	bodyReader, err := os.Open("../test/test.xml")

	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest("POST", "/person", bodyReader)

	if err != nil {
		t.Fatal(err)
	}

	r.Header.Set("Content-Type", "application/xml")

	decoder := Xml()

	candidate := new(TestXml)

	if err = decoder.Decode(r, candidate); err != nil {
		t.Fatal(err)
	}

	if candidate.Name != "Tyler Shaddix" {
		t.Fatal("Name", candidate.Name)
	}

	if candidate.Age != 23 {
		t.Fatal("Age", candidate.Age)
	}

	if candidate.Email != "tyler@example.com" {
		t.Fatal("Email", candidate.Email)
	}

	if candidate.IsAdmin != true {
		t.Fatal("IsAdmin", candidate.IsAdmin)
	}

	if candidate.HourlyRate != 22.50 {
		t.Fatal("HourlyRate", candidate.HourlyRate)
	}
}
