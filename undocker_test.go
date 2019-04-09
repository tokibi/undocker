package undocker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupMockAPI() (mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Here is a base URL")
	})
	server := httptest.NewServer(mux)
	return mux, server.URL, server.Close
}
