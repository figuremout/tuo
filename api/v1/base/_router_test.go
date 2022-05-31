package base

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v1/ping")
	if err != nil {
		t.Fatal(err)
	}
	test.CheckAPIResp(resp, t)
}
