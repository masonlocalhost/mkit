package rabbitmq

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	rmq "github.com/wagslane/go-rabbitmq"
)

// Those utils is used for some edge cases that cannot be done with "github.com/wagslane/go-rabbitmq"

func createChannel(c *rmq.Conn) (*amqp091.Channel, error) {
	conn, err := getAMQPConn(c)
	if err != nil {
		return nil, fmt.Errorf("cannot get amqp conn: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot create amqp chan: %w", err)
	}

	return ch, nil
}
func closeChannel(ch *amqp091.Channel) {
	if err := ch.Close(); err != nil {
		logrus.Warn("error closing channel")
	}
}

func QueuePurge(c *rmq.Conn, name string, noWait bool) (int, error) {
	ch, err := createChannel(c)
	if err != nil {
		return 0, fmt.Errorf("cannot create amqp chan: %w", err)
	}
	defer closeChannel(ch)

	return ch.QueuePurge(name, noWait)
}

func QueueDelete(c *rmq.Conn, name string, ifUnused bool, ifEmpty bool, noWait bool) (int, error) {
	ch, err := createChannel(c)
	if err != nil {
		return 0, fmt.Errorf("cannot create amqp chan: %w", err)
	}
	defer closeChannel(ch)

	return ch.QueueDelete(name, ifUnused, ifEmpty, noWait)
}

func QueueDeclare(c *rmq.Conn, name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) error {
	ch, err := createChannel(c)
	if err != nil {
		return fmt.Errorf("cannot create amqp chan: %w", err)
	}
	defer closeChannel(ch)

	_, err = ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)

	return err
}
