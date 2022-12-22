package http

import (
	"fmt"
	"io"
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

func (cli Client) request(url string) ([]byte, error) {
	method := http.MethodGet
	defer cli.logger.Log(method, url)

	res, err := cli.httpClient.Get(url)
	if err != nil {
		return nil, HTTPError(fmt.Errorf("problem reaching %s %s, %w", method, url, err))
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, HTTPError(fmt.Errorf("status error %d from %s %s", res.StatusCode, method, url))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, HTTPError(fmt.Errorf("read response error %w", err))
	}

	return body, nil
}

func (cli Client) Home() error {
	url := fmt.Sprintf("%s/", cli.baseURL)

	_, err := cli.request(url)
	return err
}

func (cli Client) HealthCheck() error {
	url := fmt.Sprintf("%s/healthcheck", cli.baseURL)

	_, err := cli.request(url)
	return err
}

func (cli Client) PlayBingo() (*websocket.Conn, error) {
	url := fmt.Sprintf("ws%s/bingo/play", strings.TrimPrefix(cli.baseURL, "http"))
	cli.logger.Log("PlayBingo", url)

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, HTTPError(fmt.Errorf("could not open a ws connection on %s %v", url, err))
	}

	return ws, nil
}

func (cli Client) BingoNext(w io.Writer) error {
	url := fmt.Sprintf("%s/bingo/next", cli.baseURL)

	body, err := cli.request(url)
	if err != nil {
		return err
	}

	if _, err := w.Write(body); err != nil {
		return HTTPError(fmt.Errorf("write response error %w", err))
	}

	return nil
}
