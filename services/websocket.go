package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clementus360/spacechat-text/models"
	"github.com/clementus360/spacechat-text/utils"
	"github.com/gorilla/websocket"
)

// Upgrade an http connection to a websocket connection
func WebsocketConnection(res http.ResponseWriter, req *http.Request) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			fmt.Println(r.Header.Get("Origin"))
			return r.Header.Get("Origin") == "http://127.0.0.1:5500"
		},
	}

	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Listen for messages on the socket connection
func ReceiveMessage(conn *websocket.Conn, res http.ResponseWriter) {
	for {
		// Receive a message from remote client and handle errors
		_, msg, err := conn.ReadMessage()
		if err != nil {
			utils.HandleError(err, "Failed to read message", res, http.StatusInternalServerError)
			return
		}

		// Parse the json message into a Message struct
		var message models.Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			utils.HandleError(err, "Failed to parse json", res, http.StatusInternalServerError)
			return
		}

		// Print the message (for now)
		fmt.Println(message)
	}
}
