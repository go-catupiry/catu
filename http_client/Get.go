package http_client

import (
	"net/http"
)

// Get - Start a Get response and returns the http.Response without parse data.
func Get(url string, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = headers

	return HttpClient.Do(req)
}
