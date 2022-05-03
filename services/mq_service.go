package services

import (
	"chatRoom/enviroment"
	"chatRoom/models"
	"chatRoom/utilities"
	"crypto/tls"
	"github.com/streadway/amqp"
)

var MqConn *amqp.Connection

func init() {
	conn := Connect()

	MqConn = conn
	ensureExchange()
}

func Connect() *amqp.Connection {
	//conn, err := amqp.Dial(enviroment.MqCon)

	conn, err := amqp.DialTLS(enviroment.MqCon, &tls.Config{
		ServerName: enviroment.MqServer,
	})

	/*
		conn, err := amqp.DialTLS(enviroment.MqCon, &tls.Config{
			ServerName:
		})*/

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
