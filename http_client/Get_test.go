package http_client

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockClient is the mock client
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	// GetDoFunc fetches the mock client's `Do` func
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

// Do is the mock client's `Do` func
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func TestGet(t *testing.T) {
	mockJSON := `{"online":true}`
	// r := ioutil.NopCloser(bytes.NewReader([]byte(mockJSON)))

	t.Run("Get should work with mock and return the mocked data", func(t *testing.T) {
		// GetDoFunc = func(*http.Request) (*http.Response, error) {
		// 	return &http.Response{
		// 		StatusCode: 200,
		// 		Body:       r,
		// 	}, nil
		// }

		url := "https://linkysystems.com/health.json"
		var headers http.Header

		resp, err := Get(url, headers)

		assert.NotNil(t, resp)
		assert.Nil(t, err)

		defer resp.Body.Close()

		rdrBody := io.Reader(resp.Body)
		bodyBytes, err := ioutil.ReadAll(rdrBody)

		assert.Nil(t, err)

		t.Log("Response body:", string(bodyBytes))

		assert.EqualValues(t, 200, resp.StatusCode)
		assert.EqualValues(t, mockJSON, string(bodyBytes))
	})
}
