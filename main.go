package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/clementus360/spacechat-text/config"
	"github.com/clementus360/spacechat-text/controllers"
	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnv()

	router := mux.NewRouter()

	router.HandleFunc("/api/socket", controllers.SocketHandler)
	fmt.Println("test")

	err:=http.ListenAndServe(":3002", router)
	if err!=nil {
		log.Fatal("Failed to start server: ",err)
	}
}
