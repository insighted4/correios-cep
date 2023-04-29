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
