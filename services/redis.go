package services

import (
	"context"
	"encoding/json"
	"fmt"
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

func StoreSocket(conn *websocket.Conn, userId string, rdb *redis.Client, ctx context.Context) error{

	conBytes,err := json.Marshal(conn)
	if err!=nil {
		return fmt.Errorf("failed to transform socket to json")
	}

	rdb.Set(ctx, userId, conBytes, 0)

	return nil
}

func GetSocket(userId string, rdb *redis.Client, ctx context.Context) (*websocket.Conn,error) {
	connString,err := rdb.Get(ctx, userId).Result()
	if err!=nil {
		return nil,err
	}

	var conn *websocket.Conn
	err=json.Unmarshal([]byte(connString), &conn)
	if err!=nil {
		return nil,err
	}

	return conn,nil
}
