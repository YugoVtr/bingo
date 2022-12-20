package http

import "net/http"

func NewServer(config ServerConfig) *http.Server {
	return &http.Server{
		Addr:    config.TCPAddress,
		Handler: newRouter(),
	}
}
