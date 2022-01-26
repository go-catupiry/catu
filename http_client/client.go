package http_client

import (
	"net/http"
	"time"

	"github.com/go-catupiry/catu/configuration"
)

// CustomHTTPClient - Custom http client required to make requests testable
type CustomHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	HttpClient CustomHTTPClient
)

func Init() {
	httpClientTimeout := configuration.GetInt64Env("HTTP_CLIENT_TIMEOUT", 120)

	timeout := time.Second * time.Duration(httpClientTimeout)
	HttpClient = &http.Client{Timeout: timeout}
}
