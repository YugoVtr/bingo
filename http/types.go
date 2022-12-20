package http

type Logger interface {
	Log(...interface{})
	Logf(string, ...interface{})
}

type HTTPError error

type ServerConfig struct {
	TCPAddress string
}
