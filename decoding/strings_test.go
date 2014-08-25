package decoding

import (
	"net/http"
	"testing"
)

type (
	TestStrings struct {
		Name  string `test:"name"`
		Email string `test:"email"`
		Age   int    `test:"age"`
	}

	TestStringer struct {
		values map[string]string
	}
)

func (self *TestStringer) Len(r *http.Request) int {
	return len(self.values)
}

func (self *TestStringer) Get(r *http.Request, name string) string {
	return self.values[name]
}

func TestStringsDecoder(t *testing.T) {

	stringer := &TestStringer{
		values: map[string]string{
			"email": "example@test.com",
			"name":  "Tyler",
			"age":   "22",
		},
	}

	decoder := Strings(stringer, "test")

	req, _ := http.NewRequest("POST", "/test", nil)

	candidate := new(TestStrings)

	if err := decoder.Decode(req, candidate); err != nil {
		t.Fatal(err)
	}

	if candidate.Age != 22 {
		t.Fatal("Age", candidate.Age)
	}

	if candidate.Email != "example@test.com" {
		t.Fatal("Email", candidate.Email)
	}

	if candidate.Name != "Tyler" {
		t.Fatal("Name", candidate.Name)
	}
}
