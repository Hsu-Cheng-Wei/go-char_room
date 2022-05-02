package services

import (
	"chatRoom/enviroment"
	"chatRoom/models"
	"chatRoom/utilities"
	"github.com/streadway/amqp"
)

var MqConn *amqp.Connection

func init() {
	conn := Connect()

	MqConn = conn
	ensureExchange()
}

func Connect() *amqp.Connection {
	conn, err := amqp.Dial(enviroment.MqCon)

	utilities.FailOnError(err, "Failed to connect to RabbitMQ")

	return conn
}

func ensureExchange() {
	ch, err := MqConn.Channel()
	defer ch.Close()

	utilities.FailOnError(err, "Failed to open a channel")

	if err := ch.ExchangeDeclare(
		models.ChatRoomExchange,
		amqp.ExchangeTopic,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		utilities.FailOnError(err, "Failed to declare a exchange")
	}
}
