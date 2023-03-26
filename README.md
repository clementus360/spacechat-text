# Go Messaging Chat Service using WebSockets and RabbitMQ

The Go Messaging Chat Service is a simple messaging application that enables users to send messages to each other. It is built using Go, WebSockets, and RabbitMQ. When a user sends a message, it will contain the receiver's identifier. The server will then send the message to the recipient. If the recipient is offline, RabbitMQ will be used to queue the messages until the user comes back online.

## Architecture

The Go Messaging Chat Service uses a client-server architecture. The server is built using Go and uses WebSockets to communicate with clients. When a client sends a message, the server receives the message and sends it to the recipient. If the recipient is offline, the message is queued using RabbitMQ. When the recipient comes back online, the server retrieves any queued messages from RabbitMQ and sends them to the recipient.

## Requirements

The following software is required to run the Go Messaging Chat Service:

**Go**<br />
**RabbitMQ**
