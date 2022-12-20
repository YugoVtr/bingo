package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
