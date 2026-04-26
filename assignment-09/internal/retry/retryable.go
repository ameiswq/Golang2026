package retry

import (
	"errors"
	"net"
	"net/http"
	"net/url"
)

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true
		}
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			return true
		}
		return true
	}
	if resp == nil {
		return false
	}
	switch resp.StatusCode {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	case http.StatusUnauthorized, http.StatusNotFound:
		return false
	default:
		return false
	}
}
