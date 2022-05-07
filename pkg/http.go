package pkg

import (
	"fmt"
	"net/http"
	"time"
)

var httpClient = http.Client{Timeout: time.Duration(30 * 1_000_000_000)}

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("Status %d: Message: %s: %v", r.StatusCode, http.StatusText(r.StatusCode), r.Err)
}
