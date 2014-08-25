package encoding

import (
	"net/http"
	"os"
	"testing"
)

type (
	TestJson struct {
		Name       string  `json:"name"`
		Email      string  `json:"email"`
		IsAdmin    bool    `json:"isAdmin"`
		Age        int     `json:"age"`
		HourlyRate float32 `json:"hourlyRate"`
	}
)

func TestJsonDecoder(t *testing.T) {
	bodyReader, err := os.Open("../test/test.json")

	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest("POST", "/person", bodyReader)

	if err != nil {
		t.Fatal(err)
	}

	r.Header.Set("Content-Type", "application/json")

	decoder := JsonDecode()

	candidate := new(TestJson)

	if err = decoder.Decode(r, candidate); err != nil {
		t.Fatal(err)
	}

	if candidate.Name != "Tyler Shaddix" {
		t.Fail()
	}

	if candidate.Age != 23 {
		t.Fail()
	}

	if candidate.Email != "tyler@example.com" {
		t.Fail()
	}

	if candidate.IsAdmin != true {
		t.Fail()
	}

	if candidate.HourlyRate != 22.50 {
		t.Fail()
	}
}
