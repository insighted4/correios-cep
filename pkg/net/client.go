// Copyright 2023 The Correios CEP Admin Authors
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
