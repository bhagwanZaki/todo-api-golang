package main

import (
	"common/env"
	"common/queue"
	"common/logger"
	"encoding/json"
	"log"
)

type EmailType struct {
	To      string
	Subject string
	Body    string
}

func main() {
	DEBUG := false
	envErr := env.LoadEnv(DEBUG)

	if envErr != nil {
		log.Fatalln("Failed to env file ",envErr.Error())
		return
	}

	ch := queue.InitializeQueue()
	defer ch.Close()
	q := queue.StartQueue(ch, "email")
	InitSMTP()

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	logger.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			requeueCount := 0
			if count, ok := d.Headers["x-requeue-count"]; ok {
				requeueCount = count.(int)
			}

			log.Printf("Received a message: %s", d.Body)

			var email EmailType
			err := json.Unmarshal(d.Body, &email)

			if err != nil {
				logger.ReQueueError(err, "Failed to decode json")
				d.Nack(false, true)
			}

			err = SendEmail(email.To, email.Subject, email.Body)

			if err != nil {
				requeueCount++
				d.Headers["x-requeue-count"] = requeueCount

				if requeueCount == 3 {
					logger.DeadTaskError(err, "Failed to send message adding to DLQ : ")
					d.Nack(false, false)
				} else {
					logger.ReQueueError(err, "Failed to send message adding to REQUEUE : ")
					d.Nack(false, true)
				}
			}

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
