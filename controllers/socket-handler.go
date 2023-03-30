package controllers

import (
	"context"
	"fmt"
	"net/http"

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
	conn, err := services.WebsocketConnection(res, req)
	if err != nil {
		utils.HandleError(err, "Failed to establish socket connection", res, http.StatusInternalServerError)
		return
	}

	// Store the socket using the userId(Phone number) as the key
	err = services.StoreSocket(conn, phoneNumber, rdb, context.Background())
	if err != nil {
		utils.HandleError(err, "Failed to save socket to redis", res, http.StatusInternalServerError)
		return
	}

	// Make sure that the socket will close to avoid leaks
	defer func() {
		err = conn.Close()
		if err != nil {
			fmt.Println("Failed to close websocket connection")
		}
	}()

	go services.ReceiveMessage(conn, res)

}
