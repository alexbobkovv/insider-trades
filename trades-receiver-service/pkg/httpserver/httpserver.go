package httpserver

import "net/http"

func New(handler http.Handler, port string) *http.Server {
	srv := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	return srv
}
