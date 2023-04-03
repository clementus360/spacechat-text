package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/clementus360/spacechat-text/config"
	"github.com/clementus360/spacechat-text/controllers"
	"github.com/clementus360/spacechat-text/services"
	"github.com/gorilla/mux"
)

var poolInstance *services.ConnectionPool
var once sync.Once

func GetConnectionPool() *services.ConnectionPool {
	once.Do(func() {
		poolInstance = &services.ConnectionPool{
			MaxConnections:           5,
			MaxChannelsPerConnection: 10,
		}
		if err := poolInstance.InitializePool(); err != nil {
			panic(fmt.Sprintf("Failed to initialize the connection pool: %v", err))
		}
	})

	return poolInstance
}

func main() {
	config.LoadEnv()

	pool := GetConnectionPool()

	router := mux.NewRouter()

	router.HandleFunc("/api/socket/{phoneNumber}", controllers.SocketHandler(pool))

	err := http.ListenAndServe(":3002", router)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
