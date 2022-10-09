package main

import (
	"context"
	"encoding/json"

	rmq "github.com/memnix/rabbitmq-tools"

	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var mg MongoInstance

func init() {
	LoadVar()
}

func main() {

	err := Connect()
	if err != nil {
		failOnError(err, "Failed to connect to mongoDB")
	}

	defer func() {
		log.Println("Disconnect from mongoDB")
		err := mg.Client.Disconnect(context.TODO())
		if err != nil {
			failOnError(err, "Failed to disconnect from mongoDB")
		}
	}()

	connection := new(rmq.RabbitMQConnection)

	err = connection.InitConnection(RabbitMQURL, LOGS_EXCHANGE)
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
	}

	log.Println("Connected to RabbitMQ")

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
	queueMap["warning"] = "warning"
	queueMap["info"] = "info"

	queues := []rmq.Queue{{
		Keys: []string{"error.#"},
		Name: "error",
	}, {
		Keys: []string{"info.#"},
		Name: "info",
	}, {
		Keys: []string{"warning.#"},
		Name: "warning",
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
				logObject := new(Log)
				err := json.Unmarshal(d.Body, &logObject)
				queue, err := connection.GetQueue(q)
				if err != nil {
					log.Printf("Failed to get queue %s : %s", q, err)
					continue
				}
				collection := mg.Db.Collection(queueMap[queue.Name])
				_, err = collection.InsertOne(context.TODO(), logObject)
				if err != nil {
					log.Printf("Error while inserting log to mongo: %s", err)
					continue
				}
			}
		}(q, d)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
