package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clementus360/spacechat-text/services"
	"github.com/clementus360/spacechat-text/utils"
	"github.com/gorilla/mux"
)

func SocketHandler(pool *services.ConnectionPool) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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

		// Delete socket from redis if the connection is closed
		conn.SetCloseHandler(func(code int, text string) error {
			err := services.DeleteSocket(phoneNumber, rdb, context.Background())
			if err != nil {
				utils.HandleError(err, "Failed to delete socket", res, http.StatusInternalServerError)
				return err
			}

			return nil
		})

		// Create a channel to prevent the socket from closing
		done := make(chan struct{})

		go func() {
			defer close(done)
			services.ReceiveMessage(conn, res, pool, rdb)
		}()

		// Make sure that the socket will close after the channel is closed to avoid leaks
		<-done
		err = conn.Close()
		if err != nil {
			fmt.Println("Failed to close websocket connection")
			return
		}

	}
}
