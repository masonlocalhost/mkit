package httpclient

import "fmt"

var (
	ErrRequestTimeout           = fmt.Errorf("request timeout")
	ErrMaxRetryReached          = fmt.Errorf("max retries reached")
	ErrCannotDecodeResponseBody = fmt.Errorf("cannot decode response body")
	errStatusCodeShouldRetry    = fmt.Errorf("status code sould retry")
)
