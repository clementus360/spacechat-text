package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Connect to a redis instance using a URI from environment variables
func ConnectRedis() *redis.Client {
	REDIS_URI := os.Getenv("REDIS_URI")
	rdb := redis.NewClient(&redis.Options{
		Addr:     REDIS_URI,
		Password: "",
		DB:       0,
	})

	return rdb
}

// Storing the socker remote address in redis with the userID(Phone number) as the key
func StoreSocket(conn *websocket.Conn, userId string, rdb *redis.Client, ctx context.Context) error {

	// Getting the remote address string from the websocket connectionn
	SerializedSocket, err := SerializeSocket(conn)
	if err != nil {
		return err
	}

	// Storing the remote address
	err = rdb.Set(ctx, userId, SerializedSocket, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// Getting the remote socket from redis and retrieve the websocket connection
func GetSocket(userId string, rdb *redis.Client, ctx context.Context) (*websocket.Conn, error) {

	// Get the connection string from redis
	SerializedSocket, err := rdb.Get(ctx, userId).Result()
	if err != nil {
		return nil, err
	}

	DeserializedSocket, err := DeserializeSocket([]byte(SerializedSocket))
	if err != nil {
		return nil, err
	}

	fmt.Println(DeserializedSocket)

	return DeserializedSocket, nil
}

// Delete a socket from redis
func DeleteSocket(userId string, rdb *redis.Client, ctx context.Context) error {
	err := rdb.Del(ctx, userId).Err()
	if err != nil {
		return err
	}

	return nil
}

func SerializeSocket(conn *websocket.Conn) ([]byte, error) {

	connData := make(map[string]interface{})
	connData["localAddr"] = conn.LocalAddr().String()
	connData["remoteAddr"] = conn.RemoteAddr().String()
	connData["subprotocol"] = conn.Subprotocol()

	return json.Marshal(connData)
}

func DeserializeSocket(data []byte) (*websocket.Conn, error) {

	var connData map[string]string
	err := json.Unmarshal(data, &connData)
	if err != nil {
		return nil, err
	}

	RemoteAddress := "ws://" + connData["remoteAddr"]
	u, err := url.Parse(RemoteAddress)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Add("Sec-WebSocket-Protocol", connData["subprotocol"])

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return nil, err
	}

	return conn, nil

}
