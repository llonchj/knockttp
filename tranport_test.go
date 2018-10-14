package knockttp_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/llonchj/knockttp"
)

func TestTransport(t *testing.T) {
	sd, err := knockttp.NewRoutesFromFile("fixtures/www.example.com/_server.json")
	if err != nil {
		t.Fatal(err)
	}
	network := knockttp.Network{
		Services: knockttp.Services{
			"www.example.com": *sd,
		},
	}
	data := knockttp.Data{}

	client := http.Client{
		Transport: knockttp.NewTransport(network, data),
	}

	req, err := http.NewRequest("GET", "http://www.example.com/redirect", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("response status mismatch: '%d'", res.StatusCode)
	}

	b := bytes.NewBuffer([]byte{})
	b.ReadFrom(res.Body)
	if b.String() != "Hello World!" {
		t.Fatalf("response body mismatch: '%s'", b.String())
	}
}
