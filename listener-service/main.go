package main

import (
	"fmt"
	"github.com/kleczynski/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer func(rabbitConn *amqp.Connection) {
		err := rabbitConn.Close()
		if err != nil {

		}
	}(rabbitConn)

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages ")
	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}
	// watch the queue and consume events
	if err := consumer.Listen([]string{"log.INFO", "log.ERROR", "log.WARNING"}); err != nil {
		log.Println(err)
	}

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("Error: RabbitMQ not yet ready")
			counts++
		} else {
			connection = c
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}
	log.Println("Connected to RabbitMQ")
	return connection, nil
}
