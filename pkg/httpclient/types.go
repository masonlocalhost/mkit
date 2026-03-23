package httpclient

import "time"

type Config struct {
	DefaultTimeout    time.Duration
	DefaultRetryCount int
}
