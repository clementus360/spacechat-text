package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/clementus360/spacechat-text/models"
	"github.com/clementus360/spacechat-text/utils"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Upgrade an http connection to a websocket connection
func WebsocketConnection(res http.ResponseWriter, req *http.Request) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		// Add a blacklist so I can block some IPs
		CheckOrigin: func(r *http.Request) bool {
			disallowedOrigins := []string{"http://example.com", "https://example2.com"}
			origin := r.Header.Get("Origin")
			for _, disallowed := range disallowedOrigins {
				if disallowed == origin {
					return false
				}
			}
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
func ReceiveMessage(conn *websocket.Conn, res http.ResponseWriter, pool *ConnectionPool) {
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

		fmt.Println("received message:", message)

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

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Marshal message
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Queue Message
	err = channel.PublishWithContext(ctx,
		"chat_messages",
		"user."+message.Receiver,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		},
	)

	if err != nil {
		fmt.Println("testgo")
		return err
	}

	return nil
}

func RelayMessage(queue string, conn *websocket.Conn, pool *ConnectionPool) error {
	channel, err := pool.GetChannel()
	if err != nil {
		return err
	}

	defer pool.ReleaseChannel(channel)

	// Create a channel to signal when the WebSocket connection is closed
	connClosed := make(chan struct{})
	var closeOnce sync.Once

	// Set up event listeners for the WebSocket connection
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("WebSocket connection closed")
		closeOnce.Do(func() {
			close(connClosed)
		})
		return nil
	})

	// Start a goroutine to read messages from the queue
	msgs, err := channel.Consume(queue, queue, false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			err := conn.WriteMessage(websocket.TextMessage, msg.Body)
			if err != nil {
				fmt.Println("WebSocket write error:", err)
				break
			}
			msg.Ack(false)
			fmt.Println("relayed message:", msg)
		}
		closeOnce.Do(func() {
			close(connClosed)
		})
	}()

	// Start a goroutine to check for WebSocket connection closure
	go func() {
		<-connClosed
		fmt.Println("WebSocket connection closed. Cancelling consumer...")
		channel.Cancel(queue, false)
	}()

	return nil
}
