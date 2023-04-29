package errors

import (
	"strings"
)

// IsObjectNotFoundErr returns true if the Go command line
// hints at an object not found.
func IsObjectNotFoundErr(err error) bool {
	return strings.Contains(err.Error(), "remote: Object not found")
}
