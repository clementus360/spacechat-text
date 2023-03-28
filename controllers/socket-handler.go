package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clementus360/spacechat-text/models"
	"github.com/clementus360/spacechat-text/services"
	"github.com/clementus360/spacechat-text/utils"
	"github.com/gorilla/mux"
)

func SocketHandler(res http.ResponseWriter, req *http.Request) {
	// Get the userId(Phone number) from thhe request parameters
	vars := mux.Vars(req)
	phoneNumber := vars["phoneNumber"]
	rdb := services.ConnectRedis()

	// Create a new websocket connection
	conn,err := services.WebsocketConnection(res, req)
	if err!=nil {
		utils.HandleError(err, "Failed to establish socket connection", res, http.StatusInternalServerError)
		return
	}

	// Store the socket using the userId(Phone number) as the key
	err = services.StoreSocket(conn, phoneNumber, rdb, context.Background())
	if err!=nil {
		utils.HandleError(err, "Failed to save socket to redis", res, http.StatusInternalServerError)
		return
	}

	// Make sure that the socket will close to avoid leaks
	defer func() {
		fmt.Println("is this running?")
		err=conn.Close()
		if err!=nil {
			fmt.Println("Failed to close websocket connection")
		}
	}()

	// Listen for messages on the socket connection
	for {

		// Receive a message from remote client and handle errors
		_,msg,err := conn.ReadMessage()
		if err!=nil {
			utils.HandleError(err,"Failed to read message", res, http.StatusInternalServerError)
			return
		}

		// Parse the json message into a Message struct
		var message models.Message
		err = json.Unmarshal(msg, &message)
		if err!=nil {
			utils.HandleError(err, "Failed to parse json", res, http.StatusInternalServerError)
			return
		}

		// Print the message (for now)
		fmt.Println(message)
	}

}
