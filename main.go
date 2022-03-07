package main

import (
	"context"
	"encoding/json"

	rmq "github.com/memnix/rabbitmq-tools"

	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var mg MongoInstance

func main() {
	err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		fmt.Println("Disconnect")
		err := mg.Client.Disconnect(context.TODO())
		if err != nil {
			return
		}
	}()

	connection := new(rmq.RabbitMQConnection)

	err = connection.InitConnection("amqp://guest:guest@localhost:5672/", "logs")
	if err != nil {
		return
	}

	fmt.Println("Connected")

	defer func() {
		err := connection.CloseConnection()
		if err != nil {
			failOnError(err, "Failed to close connection")
		}
	}()

	defer func() {
		err := connection.CloseChannel()
		if err != nil {
			failOnError(err, "Failed to close channel")
		}
	}()

	queueMap := make(map[string]string)
	queueMap["error"] = "error"
	queueMap["info"] = "info"

	queues := []rmq.Queue{{
		Keys: []string{"error.#"},
		Name: "error",
	}, {
		Keys: []string{"info.#"},
		Name: "info",
	}}

	err = connection.AddQueues(queues)
	if err != nil {
		return
	}

	forever := make(chan bool)

	deliveries, err := connection.Consume()
	if err != nil {
		panic(err)
	}

	for q, d := range deliveries {
		go func(q string, delivery <-chan amqp.Delivery) {
			for d := range delivery {
				log.Printf("Received a message: %s from %s", d.Body, q)
				logObject := new(Log)
				err := json.Unmarshal(d.Body, &logObject)
				collection := mg.Db.Collection(queueMap[connection.GetQueue(q).Name])
				_, err = collection.InsertOne(context.TODO(), logObject)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}(q, d)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
