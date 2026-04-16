package rabbitmq

import (
	"fmt"
	"log/slog"
	"mkit/pkg/config"
	"reflect"
	"sync"
	"unsafe"

	"github.com/rabbitmq/amqp091-go"
	rmq "github.com/wagslane/go-rabbitmq"
)

func New(config *config.App, logger *slog.Logger) (*rmq.Conn, error) {
	var (
		rabbitCfg = config.RabbitMQ
		connUrl   = fmt.Sprintf(
			"amqp://%s:%s@%s:%d",
			rabbitCfg.User, rabbitCfg.Password, rabbitCfg.Host, rabbitCfg.Port,
		)
	)

	conn, err := rmq.NewConn(
		connUrl,
		rmq.WithConnectionOptionsLogger(NewLogger(logger, slog.LevelError)),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot init rabbitmq service: %w", err)
	}

	if _, err := getAMQPConn(conn); err != nil {
		return nil, fmt.Errorf("cannot getAMQPConn from conn: %v", err)
	}

	return conn, nil
}

func getAMQPConn(c *rmq.Conn) (*amqp091.Connection, error) {
	if c == nil {
		return nil, fmt.Errorf("nil Conn")
	}

	val := reflect.ValueOf(c).Elem()
	mgrField := val.FieldByName("connectionManager")
	if !mgrField.IsValid() || mgrField.IsNil() {
		return nil, fmt.Errorf("manager nil")
	}

	mgrPtr := reflect.NewAt(mgrField.Type(), unsafe.Pointer(mgrField.UnsafeAddr())).Elem()
	mgrStruct := mgrPtr.Elem()

	muField := mgrStruct.FieldByName("connectionMu")
	if !muField.IsValid() || muField.IsNil() {
		return nil, fmt.Errorf("mutex nil or not found")
	}

	muPtr := reflect.NewAt(muField.Type(), unsafe.Pointer(muField.UnsafeAddr())).Elem()
	locker, ok := muPtr.Interface().(sync.Locker)
	if !ok {
		return nil, fmt.Errorf("field is not a sync.Locker")
	}

	locker.Lock()
	defer locker.Unlock()

	connField := mgrStruct.FieldByName("connection")
	if !connField.IsValid() || connField.IsNil() {
		return nil, fmt.Errorf("connection nil")
	}

	connVal := reflect.NewAt(connField.Type(), unsafe.Pointer(connField.UnsafeAddr())).Elem()

	return connVal.Interface().(*amqp091.Connection), nil
}
