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

var pool *services.ConnectionPool
var once sync.Once

func GetConnectionPool() *services.ConnectionPool {
	once.Do(func() {
		pool = &services.ConnectionPool{
			MaxConnections:           1,
			MaxChannelsPerConnection: 5,
		}
		if err := pool.InitializePool(); err != nil {
			panic(fmt.Sprintf("Failed to initialize the connection pool: %v", err))
		}
	})

	return pool
}

func main() {
	config.LoadEnv()

	router := mux.NewRouter()

	router.HandleFunc("/api/socket/{phoneNumber}", controllers.SocketHandler)

	err := http.ListenAndServe(":3002", router)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
