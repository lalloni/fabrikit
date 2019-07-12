package status

import (
	"net/http"
)

const (
	OK = 200 // RFC 7231, 6.3.1

	BadRequest = 400 // RFC 7231, 6.5.1
	Forbidden  = 403 // RFC 7231, 6.5.3
	NotFound   = 404 // RFC 7231, 6.5.4

	Error          = 500 // RFC 7231, 6.6.1
	NotImplemented = 501 // RFC 7231, 6.6.2
)

func Text(status int) string {
	return http.StatusText(status)
}
