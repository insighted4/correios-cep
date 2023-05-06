package net

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	"go.opencensus.io/plugin/ochttp"
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:10.0) Gecko/20100101 Firefox/10.0"

func NewClient() *resty.Client {
	client := resty.New().
		SetHeader("User-Agent", UserAgent).
		SetTransport(
			&ochttp.Transport{
				Base: http.DefaultTransport,
			},
		)

	return client
}
