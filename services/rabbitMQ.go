package services

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectMQ() (*amqp.Connection,error){
	MQ_URI := os.Getenv("MQ_URI")
	conn,err := amqp.Dial(MQ_URI)
	if err!= nil {
		return nil,err
	}

	return conn,nil
}

func CreateChannel(MQConnection *amqp.Connection, ChannelName string) (*amqp.Channel,error){
	ch,err := MQConnection.Channel()
	if err!=nil {
		return nil,err
	}

	return ch,nil
}

func QueueMessage() {

}
