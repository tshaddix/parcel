package encoding

import (
	"net/http"
	"testing"
)

type (
	TestQuery struct {
		SortBy       string `query:"sort-by"`
		Limit        int    `query:"limit"`
		ShowInactive bool   `query:"show-inactive"`
		Ids          []int  `query:"id"`
	}
)

func TestQueryDecoder(t *testing.T) {

	decoder := Query()

	req, _ := http.NewRequest("GET", "/list?sort-by=name&limit=100&show-inactive=true&id=1&id=2", nil)

	candidate := new(TestQuery)

	if err := decoder.Decode(req, candidate); err != nil {
		t.Fatal(err)
	}

	if candidate.SortBy != "name" {
		t.Fail()
	}

	if candidate.Limit != 100 {
		t.Fail()
	}

	if candidate.ShowInactive != true {
		t.Fail()
	}

	if len(candidate.Ids) != 2 {
		t.Fail()
	}
}
