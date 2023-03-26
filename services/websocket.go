package services

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func WebsocketConnection(res http.ResponseWriter, req *http.Request) (*websocket.Conn, error){
	var upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			fmt.Println(r.Header.Get("Origin"))
			return r.Header.Get("Origin") == "http://127.0.0.1:5500"
		},
	}

	conn,err:=upgrader.Upgrade(res,req,nil)
	if err!=nil {
		return nil,err
	}

	return conn,nil
}
