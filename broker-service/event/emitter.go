package event

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)
	return declareExchange(channel)
}

func (e *Emitter) Push(event, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)
	log.Println("Pushing to channel")
	ctx := context.Background()

	err = channel.PublishWithContext(ctx, "logs_topic", severity, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(event),
	})
	if err != nil {
		return err
	}
	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}
	if err := emitter.setup(); err != nil {
		return Emitter{}, err
	}

	return emitter, nil

}
