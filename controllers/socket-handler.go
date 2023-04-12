package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/clementus360/spacechat-text/services"
	"github.com/clementus360/spacechat-text/utils"
	"github.com/gorilla/mux"
)

func SocketHandler(pool *services.ConnectionPool) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Get the userId(Phone number) from the request parameters
		vars := mux.Vars(req)
		phoneNumber := vars["phoneNumber"]

		// Initialize redis connection
		rdb := services.ConnectRedis()

		AUTH_URI := os.Getenv("AUTH_URI")

		ticket := req.URL.Query().Get("ticket")

		resp, err := http.Post(fmt.Sprintf("%v/authorize/%v", AUTH_URI, phoneNumber), "application/x-www-form-urlencoded", strings.NewReader(ticket))
		if err != nil {
			utils.HandleError(err, "Failed to connect to authorization service", res, http.StatusInternalServerError)
			return
		}

		fmt.Println(resp.StatusCode, "=?", http.StatusOK)

		if resp.StatusCode != http.StatusOK {
			utils.HandleError(fmt.Errorf("failed to authorize the user"), "Failed to authorize the user", res, http.StatusUnauthorized)
			return
		}

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
			fmt.Println("closed?")
			err := services.DeleteSocket(phoneNumber, rdb, context.Background())
			if err != nil {
				utils.HandleError(err, "Failed to delete socket", res, http.StatusInternalServerError)
				return err
			}

			return nil
		})

		// Initialize a rabbitmq Queue
		queue, err := services.InitializeQueue(phoneNumber)
		if err != nil {
			utils.HandleError(err, "Failed to initialize the queue", res, http.StatusInternalServerError)
			return
		}

		// Create a channel to prevent the socket from closing
		done := make(chan struct{})

		go services.RelayMessage(queue, conn, pool)

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
