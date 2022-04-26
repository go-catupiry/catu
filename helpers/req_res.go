package helpers

import "net/http"

type FakeResponseWriter struct {
}

func (r *FakeResponseWriter) Header() http.Header {
	return nil
}

func (r *FakeResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}
func (r *FakeResponseWriter) WriteHeader(statusCode int) {
}
