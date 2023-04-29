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
