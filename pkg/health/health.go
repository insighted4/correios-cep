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

package health

import (
	"context"
	"time"
)

type Checker interface {
	Check(ctx context.Context) error
}

// NewCustomHealthCheckFunc returns a new health check function.
func NewCustomHealthCheckFunc(checker Checker, now func() time.Time) func(context.Context) (details interface{}, err error) {
	return func(ctx context.Context) (details interface{}, err error) {
		if err := checker.Check(ctx); err != nil {
			return nil, err
		}

		return nil, nil
	}
}
