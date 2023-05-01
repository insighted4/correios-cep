package net

import (
	"net/http"

	"go.opencensus.io/plugin/ochttp"
)

func NewClient() *http.Client {
	client := &http.Client{
		Transport: &ochttp.Transport{
			Base: http.DefaultTransport,
		},
	}

	return client
}
