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

func StoreSocket(conn *websocket.Conn, userId string, rdb *redis.Client, ctx context.Context) error{

	err := rdb.Set(ctx, userId, conn.RemoteAddr().String(), 0).Err()
	if err!=nil {
		return err
	}

	fmt.Println(conn.RemoteAddr().String())

	return nil
}

func GetSocket(userId string, rdb *redis.Client, ctx context.Context, req *http.Request) (*websocket.Conn,error) {
	connString,err := rdb.Get(ctx, userId).Result()
	if err!=nil {
		return nil,err
	}

	dialer := websocket.Dialer{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

	fmt.Println("test")
	conn,res,err := dialer.Dial("ws://"+connString, nil)
	if err!=nil {
		return nil,err
	}

	fmt.Println(res.Body," : ",res.Status)

	return conn,nil
}
