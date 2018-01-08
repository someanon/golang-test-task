package main

import (
	"net/http"
)

type client interface {
	get(url string) (*http.Response, error)
}

func newLoader() client {
	return httpClient{}
}

type httpClient struct{}

func (_ httpClient) get(url string) (*http.Response, error) {
	return http.Get(url)
}
