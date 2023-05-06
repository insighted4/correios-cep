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
	"time"
)

type Address struct {
	CEP          string     `json:"cep" db:"cep"`
	State        string     `json:"state"  db:"state"`
	City         string     `json:"city" db:"location"`
	Neighborhood string     `json:"neighborhood" db:"neighborhood"`
	Location     string     `json:"location" db:"location"`
	Children     []*Address `json:"children" db:"children"`

	CreatedAt *time.Time `json:"created_at,omitempty,omitempty"  db:"cep"`
	UpdatedAt *time.Time `json:"updated_at,omitempty,omitempty"  db:"cep"`
}
