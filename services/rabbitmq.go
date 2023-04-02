package services

import (
	"errors"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnectionPool struct {
	MaxConnections           int
	MaxChannelsPerConnection int
	Connections              []*amqp.Connection
	Channels                 []*amqp.Channel
	Mutex                    sync.Mutex
}

func ConnectMQ() (*amqp.Connection, error) {
	MQ_URI := os.Getenv("MQ_URI")
	conn, err := amqp.Dial(MQ_URI)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateChannel(MQConnection *amqp.Connection) (*amqp.Channel, error) {
	ch, err := MQConnection.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (pool *ConnectionPool) InitializePool() error {
	pool.Connections = make([]*amqp.Connection, 0, pool.MaxConnections)
	pool.Channels = make([]*amqp.Channel, 0, pool.MaxChannelsPerConnection)

	for i := 0; i <= pool.MaxConnections; i++ {
		conn, err := ConnectMQ()
		if err != nil {
			return nil
		}

		pool.Connections = append(pool.Connections, conn)

		for j := 0; j <= pool.MaxChannelsPerConnection; j++ {
			channel, err := CreateChannel(conn)

			if err != nil {
				return nil
			}

			pool.Channels = append(pool.Channels, channel)
		}
	}

	return nil
}

func (pool *ConnectionPool) GetChannel() (*amqp.Channel, error) {
	pool.Mutex.Lock()
	defer pool.Mutex.Unlock()

	if len(pool.Channels) == 0 {
		return nil, errors.New("no channels available in the pool")
	}

	channel := pool.Channels[0]
	pool.Channels = pool.Channels[1:]

	return channel, nil
}

func (pool *ConnectionPool) ReleaseChannel(channel *amqp.Channel) {
	pool.Mutex.Lock()
	defer pool.Mutex.Unlock()

	pool.Channels = append(pool.Channels, channel)
}
