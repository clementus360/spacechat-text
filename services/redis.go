package services

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Connect to a redis instance using a URI from environment variables
func ConnectRedis() *redis.Client{
	REDIS_URI:=os.Getenv("REDIS_URI")
	rdb := redis.NewClient(&redis.Options{
		Addr: REDIS_URI,
		Password: "",
		DB: 0,
	})

	return rdb
}

// Storing the socker remote address in redis with the userID(Phone number) as the key
func StoreSocket(conn *websocket.Conn, userId string, rdb *redis.Client, ctx context.Context) error{

	// Getting the remote address string from the websocket connectionn
	remoteAddress := conn.RemoteAddr().String()

	// Storing the remote address
	err := rdb.Set(ctx, userId, remoteAddress, 0).Err()
	if err!=nil {
		return err
	}

	return nil
}

// Getting the remote socket from redis and retrieve the websocket connection
func GetSocket(userId string, rdb *redis.Client, ctx context.Context, req *http.Request) (*websocket.Conn,error) {

	// Get the connection string from redis
	connString,err := rdb.Get(ctx, userId).Result()
	if err!=nil {
		return nil,err
	}

	// Create a websocket dialer used to dial the connection string
	dialer := websocket.Dialer{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

	// Dial the connection string
	conn,res,err := dialer.Dial("ws://"+connString, nil)
	if err!=nil {
		return nil,err
	}

	fmt.Println(res.Body," : ",res.Status)

	return conn,nil
}
