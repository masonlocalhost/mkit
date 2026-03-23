package httpclient

import (
	"net/http"
	"time"
)

type Response struct {
	Body        any
	StatusCode  int
	RawResponse *http.Response
}

type Request struct {
	client       *Client
	url          string
	headers      http.Header
	queries      map[string]string
	requestBody  any
	responseBody any
	timeout      time.Duration
	retryCount   int
}

func (r *Request) SetBody(body any) *Request {
	r.requestBody = body
	r.headers["Content-Type"] = []string{"application/json"}

	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.headers[key] = append(r.headers[key], value)

	return r
}

func (r *Request) SetQuery(key, value string) *Request {
	r.queries[key] = value

	return r
}

func (r *Request) SetHeaders(header http.Header) *Request {
	r.headers = header

	return r
}

func (r *Request) SetQueries(queries map[string]string) *Request {
	r.queries = queries

	return r
}

func (r *Request) SetTimeout(timeout time.Duration) *Request {
	r.timeout = timeout

	return r
}

func (r *Request) SetRetry(retryCount int) *Request {
	r.retryCount = retryCount

	return r
}

func (r *Request) SetResult(body any) *Request {
	r.responseBody = body

	return r
}
