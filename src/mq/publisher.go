package mq

import (
	"encoding/json"

	"log"
	"math/rand"
	"time"

	cnf "github.com/Tobias1R/gintonica/config"
	"github.com/streadway/amqp"
)

func Publisher() {
	conn, err := amqp.Dial(cnf.AMQPConnectionURL)
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	loops := 10
	var i = 0
	for i <= loops {
		rand.Seed(time.Now().UnixNano())

		addTask := AddTask{Number1: rand.Intn(999), Number2: rand.Intn(999)}
		body, err := json.Marshal(addTask)
		if err != nil {
			handleError(err, "Error encoding JSON")
		}

		err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})

		if err != nil {
			log.Fatalf("Error publishing message: %s", err)
		}

		log.Printf("AddTask: %d+%d", addTask.Number1, addTask.Number2)

		i++
	}

}
