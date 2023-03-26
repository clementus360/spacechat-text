package controllers

import (
	"fmt"
	"net/http"

	"github.com/clementus360/spacechat-text/services"
	"github.com/clementus360/spacechat-text/utils"
)

func SocketHandler(res http.ResponseWriter, req *http.Request) {
	conn,err := services.WebsocketConnection(res, req)
	if err!=nil {
		utils.HandleError(err, "Failed to establish socket connection", res, http.StatusInternalServerError)
		return
	}

	defer func() {
		fmt.Println("is this running?")
		err=conn.Close()
		if err!=nil {
			fmt.Println("Failed to close websocket connection")
		}
	}()

	for {
		_,msg,err := conn.ReadMessage()
		if err!=nil {
			utils.HandleError(err,"Failed to read message", res, http.StatusInternalServerError)
			return
		}

		fmt.Println(string(msg))
	}

}
