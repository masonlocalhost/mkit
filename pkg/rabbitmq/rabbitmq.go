package rabbitmq

import (
	"fmt"
	"mkit/pkg/config"
	"reflect"
	"sync"
	"unsafe"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	rmq "github.com/wagslane/go-rabbitmq"
)

func New(config *config.App, logger *logrus.Logger) (*rmq.Conn, error) {
	var (
		rabbitCfg = config.RabbitMQ
		connUrl   = fmt.Sprintf(
			"amqp://%s:%s@%s:%d",
			rabbitCfg.User, rabbitCfg.Password, rabbitCfg.Host, rabbitCfg.Port,
		)
	)

	conn, err := rmq.NewConn(
		connUrl,
		rmq.WithConnectionOptionsLogger(NewLogger(logger, logrus.ErrorLevel)), // default disable logs for lower levels
	)

	// Try getAMQPConn for early error return
	if _, err := getAMQPConn(conn); err != nil {
		return nil, fmt.Errorf("cannot getAMQPConn from conn: %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf("cannot init rabbitmq service: %w", err)
	}

	return conn, nil
}

func getAMQPConn(c *rmq.Conn) (*amqp091.Connection, error) {
	if c == nil {
		return nil, fmt.Errorf("nil Conn")
	}

	// 1. Navigate to the ConnectionManager
	val := reflect.ValueOf(c).Elem()
	mgrField := val.FieldByName("connectionManager")
	if !mgrField.IsValid() || mgrField.IsNil() {
		return nil, fmt.Errorf("manager nil")
	}

	// Create an addressable version of the manager pointer
	mgrPtr := reflect.NewAt(mgrField.Type(), unsafe.Pointer(mgrField.UnsafeAddr())).Elem()
	mgrStruct := mgrPtr.Elem()

	// 2. Access the Mutex (connectionMu)
	muField := mgrStruct.FieldByName("connectionMu")
	if !muField.IsValid() || muField.IsNil() {
		return nil, fmt.Errorf("mutex nil or not found")
	}

	// 3. Convert the private mutex to an interface we can use
	// connectionMu is *sync.RWMutex, which satisfies sync.Locker
	muPtr := reflect.NewAt(muField.Type(), unsafe.Pointer(muField.UnsafeAddr())).Elem()
	locker, ok := muPtr.Interface().(sync.Locker)
	if !ok {
		return nil, fmt.Errorf("field is not a sync.Locker")
	}

	// 4. LOCK the manager while we grab the connection
	locker.Lock()
	defer locker.Unlock()

	// 5. Now safely grab the connection
	connField := mgrStruct.FieldByName("connection")
	if !connField.IsValid() || connField.IsNil() {
		return nil, fmt.Errorf("connection nil")
	}

	// Bypass protection for the connection field itself
	connVal := reflect.NewAt(connField.Type(), unsafe.Pointer(connField.UnsafeAddr())).Elem()

	return connVal.Interface().(*amqp091.Connection), nil
}
