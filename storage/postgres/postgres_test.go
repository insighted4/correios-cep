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

package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	if shouldSkip() {
		t.SkipNow()
	}

	pg, err := newTestPostgres(t)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, pg)
	assert.NotNil(t, pg.db)
}

func TestClose(t *testing.T) {
	if shouldSkip() {
		t.SkipNow()
	}

	pg, err := newTestPostgres(t)
	if err != nil {
		t.Error(err)
	}

	assert.NotPanics(t, func() { pg.Close() })
}

func TestCheck(t *testing.T) {
	if shouldSkip() {
		t.SkipNow()
	}

	pg, err := newTestPostgres(t)
	if err != nil {
		t.Error(err)
	}

	if err := pg.Check(context.Background()); err != nil {
		t.Error(err)
	}
}
