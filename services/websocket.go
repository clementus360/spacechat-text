package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/clementus360/spacechat-text/models"
	"github.com/clementus360/spacechat-text/utils"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

// Upgrade an http connection to a websocket connection
func WebsocketConnection(res http.ResponseWriter, req *http.Request) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// CheckOrigin: func(r *http.Request) bool {
		// 	fmt.Println(r.Header.Get("Origin"))
		// 	return r.Header.Get("Origin") == "http://127.0.0.1:5500"
		// },
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Listen for messages on the socket connection
func ReceiveMessage(conn *websocket.Conn, res http.ResponseWriter, pool *ConnectionPool, rdb *redis.Client) {
	for {
		// Receive a message from remote client and handle errors
		_, msg, err := conn.ReadMessage()
		if err != nil {
			utils.HandleError(err, "Failed to read message", res, http.StatusInternalServerError)
			return
		}

		// Parse the json message into a Message struct
		var message models.Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			utils.HandleError(err, "Failed to parse json", res, http.StatusInternalServerError)
			return
		}

		// // Send messages directly
		// conn2, err := GetSocket(message.Receiver, rdb, context.Background())
		// if err != nil {
		// 	utils.HandleError(err, "Failed to get Connection from redis", res, http.StatusInternalServerError)
		// 	return
		// }

		// err = conn2.WriteJSON(message)
		// if err != nil {
		// 	utils.HandleError(err, "Failed to send message through socket", res, http.StatusInternalServerError)
		// 	return
		// }

		// Queue message
		err = QueueMessage(pool, &message)
		if err != nil {
			utils.HandleError(err, fmt.Sprintf("Failed to queue message: %v", err), res, http.StatusInternalServerError)
			return
		}

		fmt.Println("Message has been queued")

		// Print the message (for now)
		fmt.Println(message)
	}
}

func QueueMessage(pool *ConnectionPool, message *models.Message) error {
	// Get RabbitMQ channnel from pool
	channel, err := pool.GetChannel()
	if err != nil {
		return err
	}

	defer pool.ReleaseChannel(channel)

	// Declare exchange
	err = channel.ExchangeDeclare(
		"chat_messages",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Check if the queue does not exist to create a new one
	q, err := channel.QueueDeclare(message.Sender, false, false, false, false, nil)
	if err != nil {
		return err
	}

	// Bind queue to exchange
	err = channel.QueueBind(q.Name, "user."+message.Sender, "chat_messages", false, nil)
	if err != nil {
		return err
	}

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Queue Message
	err = channel.PublishWithContext(ctx,
		"chat_messages",
		"user."+message.Receiver,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message.Payload),
		},
	)

	if err != nil {
		fmt.Println("testgo")
		return err
	}

	return nil
}
