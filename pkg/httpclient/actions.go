package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/avast/retry-go/v4"
	"net/http"
	"time"
)

func (r *Request) Post(ctx context.Context, url string) (*Response, error) {
	return r.doHTTP(ctx, url, http.MethodPost)
}

func (r *Request) Get(ctx context.Context, url string) (*Response, error) {
	return r.doHTTP(ctx, url, http.MethodGet)
}

func (r *Request) prepareHTTPRequest(ctx context.Context, url string, method string) (*http.Request, error) {
	var (
		bodyReader *bytes.Reader = nil
		req        *http.Request
		err        error
	)

	if r.requestBody != nil {
		bodyBytes, err := json.Marshal(r.requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to parse request body: %w", err)
		}

		bodyReader = bytes.NewReader(bodyBytes)
		req, err = http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	for k, vs := range r.headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	q := req.URL.Query()
	for k, v := range r.queries {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return req, err
}

func (r *Request) doHTTP(ctx context.Context, url string, method string) (*Response, error) {
	client := r.client.GetHTTPClient()
	rCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	req, err := r.prepareHTTPRequest(rCtx, url, method)
	if err != nil {
		return nil, fmt.Errorf("failed to init http request: %w", err)
	}

	res, err := doWithRetry(client, req, r.retryCount)
	if err != nil {
		return nil, err
	}

	response := &Response{
		RawResponse: res,
		StatusCode:  res.StatusCode,
		Body:        r.responseBody,
	}

	if r.responseBody == nil {
		return response, nil
	}

	if res.StatusCode == http.StatusOK && response.Body != nil {
		if err := json.NewDecoder(res.Body).Decode(&response.Body); err != nil {
			return nil, errors.Join(ErrCannotDecodeResponseBody, err)
		}
	}

	return response, nil
}

func doWithRetry(client *http.Client, req *http.Request, retryCount int) (*http.Response, error) {
	res, err := retry.DoWithData(
		func() (*http.Response, error) {
			res, err := client.Do(req)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return nil, ErrRequestTimeout
				}

				return res, err
			}
			if shouldRetry(err, res) {
				return res, errStatusCodeShouldRetry
			}

			return res, err
		},
		retry.RetryIf(func(err error) bool {
			return errors.Is(err, errStatusCodeShouldRetry)
		}),
		retry.Attempts(uint(retryCount)),
		retry.Delay(3*time.Second),
		retry.MaxJitter(1*time.Second),
	)
	if err != nil {
		var retryErr retry.Error
		if errors.As(err, &retryErr) {
			lastErr := retryErr.WrappedErrors()[len(retryErr.WrappedErrors())-1]
			if res == nil {
				return nil, fmt.Errorf("error while making http request, last error is: %w", lastErr)
			}
			if errors.Is(lastErr, errStatusCodeShouldRetry) {
				// retried until max, ignore
				err = nil
			}
		}

		return nil, err
	}

	return res, nil
}

func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		return !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded))
	}

	if resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout {
		return true
	}
	return false
}
