package db

import (
	"log"
	"os"

	amqp "github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type RabbitConfig struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitConnection() (RabbitConfig, error) {
	url := os.Getenv("RABBITMQ_SERVER")

	log.Println(url)
	conn, err := amqp.Dial(url)
	if err != nil {
		return RabbitConfig{}, err
	}
	chann, err := conn.Channel()
	if err != nil {
		return RabbitConfig{}, err
	}
	return RabbitConfig{Connection: conn, Channel: chann}, nil
}
	
func (rabbitCon RabbitConfig) Close() {
	rabbitCon.Connection.Close()
	rabbitCon.Channel.Close()
}
