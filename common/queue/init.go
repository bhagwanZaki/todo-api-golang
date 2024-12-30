package queue

import (
	"common/env"
	"common/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	dlxName = "dlx_exchange"
	dlqName = "dead_letter_queue"
)

func InitializeQueue() *amqp.Channel {
	// "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(env.SQS_URL)
	logger.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	logger.FailOnError(err, "Failed to open a channel")

	err = ch.Qos(
		1,     // prefetch count (only one unacknowledged message at a time)
		0,     // prefetch size (unlimited message size)
		false, // apply to this consumer only
	)
	logger.FailOnError(err,"Failed to set Qos")

	// Declare the DLX
	err = ch.ExchangeDeclare(
		dlxName,
		"direct", // or any type needed
		true,
		false,
		false,
		false,
		nil,
	)
	logger.FailOnError(err, "Failed to declare the DLX")

	// Declare the DLQ
	_, err = ch.QueueDeclare(
		dlqName,
		true,
		false,
		false,
		false,
		nil,
	)
	logger.FailOnError(err, "Failed to declare the DLQ")

	// Bind the DLQ to the DLX
	err = ch.QueueBind(
		dlqName,
		"dlx_key",
		dlxName,
		false,
		nil,
	)
	logger.FailOnError(err, "Failed to bind the DLQ to the DLX")

	return ch
}

func StartQueue(ch *amqp.Channel, queueName string) amqp.Queue {
	queueArgs := amqp.Table{
		"x-dead-letter-exchange":    dlxName,   // Set DLX
		"x-dead-letter-routing-key": "dlx_key", // Routing key for DLX
	}

	q, err := ch.QueueDeclare(
		queueName,   // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		queueArgs, // arguments
	)
	logger.FailOnError(err, "Failed to declare a queue")

	return q
}