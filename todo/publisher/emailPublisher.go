package publisher

import (
	"common/logger"
	"context"
	"encoding/json"
	"todoGoApi/types"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ch *amqp.Channel
var ctx context.Context

func InitEmailQueue(channel *amqp.Channel, contx context.Context) {
	ch = channel
	ctx = contx
}

func AddEmailTask(to string, subject string, msg string) error {
	task := types.EmailType{
		To:      to,
		Subject: subject,
		Body:    msg,
	}

	body, err := json.Marshal(task)

	if err != nil {
		logger.Logger(err.Error(), "[AddEmailTask]")
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",      // exchange
		"email", // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		logger.Logger(err.Error(), "[AddEmailTask] Failed To Publish Task :")
		return err
	}
	logger.Logger("Add task to queue","[AddEmailTask]")
	return nil
}
