package main

import (
	"bytes"
	"net/http"
	"testing"
)

func BenchmarkPost100Re(b *testing.B) {
	var resp *http.Response
	var err error

	for n := 0; n < b.N; n++ {
		resp, err = http.DefaultClient.Post("https://reqres.in/api/register", "application/json", bytes.NewBuffer([]byte("{\"email\":\"test@test.com\",\"password\":\"password\"}")))
		if err != nil {
			b.Error(err)
		}
		resp.Body.Close()
	}
}

func BenchmarkPost100(b *testing.B) {
	for i := 0; i < 100; i++ {
		resp, err := http.DefaultClient.Post("https://reqres.in/api/register", "application/json", bytes.NewBuffer([]byte("{\"email\":\"test@test.com\",\"password\":\"password\"}")))
		if err != nil {
			b.Error(err)
		}
		resp.Body.Close()
	}
}
