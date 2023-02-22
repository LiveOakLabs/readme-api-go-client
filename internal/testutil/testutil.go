// testutil is an internal package with utilities for testing the readme package.
package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

// RoundTripFunc is used to pass requests through to evaluate the response.
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip is used to pass requests through to evaluate the response.
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls.
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// APITestResponse represents a mock HTTP response for tests.
type APITestResponse struct {
	URL     string
	Status  int
	Body    string
	Headers http.Header
}

// New wraps NewClient and the NewTestClient round trip client configured for tests.
//
// Tests use this to mock an API response.
func (r *APITestResponse) New(t *testing.T) *readme.Client {
	client, _ := readme.NewClient("hunter2", "http://readme-test.local/api/v1")
	client.HTTPClient = NewTestClient(func(req *http.Request) *http.Response {
		if r.URL != "" {
			assert.Equal(t, req.URL.String(), r.URL)
		}

		return &http.Response{
			StatusCode: r.Status,
			Body:       io.NopCloser(bytes.NewBufferString(r.Body)),
			Header:     r.Headers,
		}
	})

	return client
}

// JsonToStruct is a helper for parsing JSON to a specified struct interface.
func JsonToStruct(t *testing.T, jsonString string, object any) {
	err := json.Unmarshal([]byte(jsonString), &object)
	if err != nil {
		t.Errorf("error converting json string to struct: %v\n", jsonString)
	}
}

// StructToJson is a helper for parsing a struct to a JSON string.
func StructToJson(t *testing.T, object any) string {
	res, err := json.Marshal(object)
	if err != nil {
		t.Errorf("error converting struct to json: %v\n", err)
	}

	return string(res)
}
