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

package errors

import "net/http"

// Kind enums.
const (
	KindNotFound       = http.StatusNotFound
	KindBadRequest     = http.StatusBadRequest
	KindUnexpected     = http.StatusInternalServerError
	KindAlreadyExists  = http.StatusConflict
	KindRateLimit      = http.StatusTooManyRequests
	KindNotImplemented = http.StatusNotImplemented
	KindRedirect       = http.StatusMovedPermanently
)

// IsNotFoundErr helper function for KindNotFound.
func IsNotFoundErr(err error) bool {
	return Kind(err) == KindNotFound
}
