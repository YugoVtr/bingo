package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     Logger
}

func NewClient(baseURL string, log Logger) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: time.Second},
		logger:     log,
	}
}

func (cli Client) HealthCheck() error {
	url, method := fmt.Sprintf("%s/healthcheck", cli.baseURL), http.MethodGet
	cli.logger.Log(method, url)

	res, err := cli.httpClient.Get(url)
	if err != nil {
		return HTTPError(fmt.Errorf("problem reaching %s %s, %w", method, url, err))
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return HTTPError(fmt.Errorf("status error %d from %s %s", res.StatusCode, method, url))
	}

	return nil
}

func (cli Client) WebSocket() error {
	url := fmt.Sprintf("ws%s/ws", strings.TrimPrefix(cli.baseURL, "http"))
	cli.logger.Log("WebSocket", url)

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return HTTPError(fmt.Errorf("could not open a ws connection on %s %v", url, err))
	}
	defer ws.Close()

	return nil
}
