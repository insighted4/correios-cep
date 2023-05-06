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

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPagination(t *testing.T) {
	p := NewPagination(20, 5)
	assert.Equal(t, 20, p.Limit)
	assert.Equal(t, 100, p.Offset)
}

func TestNewPaginationWithDefault(t *testing.T) {
	p := NewPagination(PaginationLimit+1, 0)
	assert.Equal(t, PaginationLimit, p.Limit)
	assert.Equal(t, 0, p.Offset)
}
