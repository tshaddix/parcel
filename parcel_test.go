package parcel

import (
	"net/http"
	"os"
	"testing"

	"github.com/tshaddix/parcel/encoding"
)

type (
	TestPerson struct {
		Name         string  `xml:"name" json:"name"`
		Email        string  `xml:"email" json:"email"`
		IsAdmin      bool    `xml:"is-admin" json:"isAdmin"`
		Age          int     `xml:"age" json:"age"`
		HourlyRate   float32 `xml:"hourly-rate" json:"hourlyRate"`
		AccessToken  string  `query:"access-token"`
		ShowMatching bool    `query:"show-matching"`
	}
)

func TestDecoding(t *testing.T) {
	factory := NewFactory()
	factory.Use(encoding.JSON())
	factory.Use(encoding.XML())
	factory.Use(encoding.Query())

	jsonReader, err := os.Open("test/test.json")

	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest("POST", "/person?access-token=123456&show-matching=true", jsonReader)

	if err != nil {
		t.Fatal(err)
	}

	stat, err := jsonReader.Stat()

	if err != nil {
		t.Fatal(err)
	}

	r.ContentLength = stat.Size()

	r.Header.Set("Content-Type", "application/json")

	p := factory.Parcel(nil, r)

	person1 := new(TestPerson)

	if err = p.Decode(person1); err != nil {
		t.Fatal(err)
	}

	if person1.Name != "Tyler Shaddix" {
		t.Fatal("Name", person1.Name)
	}

	if person1.Email != "tyler@example.com" {
		t.Fatal("Email", person1.Email)
	}

	if person1.AccessToken != "123456" {
		t.Fatal("AccessToken", person1.AccessToken)
	}
}
