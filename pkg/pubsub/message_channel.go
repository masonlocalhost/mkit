package pubsub

import (
	"context"
	"fmt"
	rabbitmq2 "mkit/pkg/rabbitmq"
	"time"

	"buf.build/go/protovalidate"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func getRoutingKey(exchange string, topic string, transmissionType TransmissionType) string {
	return fmt.Sprintf("%s.%s.%d", exchange, topic, transmissionType)
}

// handler can receive proto message or rabbitmq message if needed
type handler func(msg proto.Message, d rabbitmq.Delivery)

type MessageChannel struct {
	serviceName          string
	broadcastConsumer    *rabbitmq.Consumer
	broadcastAllConsumer *rabbitmq.Consumer
	broadcastHandlers    []handler
	broadcastAllHandlers []handler
	logger               *logrus.Logger
	ctx                  context.Context
	topic                string
	unmarshalPayload     proto.Message
	validator            protovalidate.Validator
	exchangeName         string
}

func NewMsgChannel(
	logger *logrus.Logger, serviceName, exchangeName, topic string, unmarshalPayload proto.Message,
) (*MessageChannel, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to init proto validator: %w", err)
	}
	proto.Reset(unmarshalPayload)

	return &MessageChannel{
		serviceName:      serviceName,
		logger:           logger,
		topic:            topic,
		unmarshalPayload: unmarshalPayload,
		validator:        validator,
		exchangeName:     exchangeName,
	}, nil
}

func (m *MessageChannel) AddHandler(handler handler, transmissionType TransmissionType) {
	switch transmissionType {
	case Broadcast:
		m.broadcastHandlers = append(m.broadcastHandlers, handler)
	case BroadcastAll:
		m.broadcastAllHandlers = append(m.broadcastAllHandlers, handler)
	}
}

func (m *MessageChannel) StartConsuming(rctx context.Context, conn *rabbitmq.Conn) error {
	errGroup, _ := errgroup.WithContext(rctx)
	// start all instance consumer
	if len(m.broadcastAllHandlers) > 0 {
		baConsumer, err := rabbitmq.NewConsumer(conn, "",
			rabbitmq.WithConsumerOptionsExchangeName(m.exchangeName),
			rabbitmq.WithConsumerOptionsRoutingKey(getRoutingKey(m.exchangeName, m.topic, BroadcastAll)),
			rabbitmq.WithConsumerOptionsExchangeDeclare,
			rabbitmq.WithConsumerOptionsExchangeKind("topic"),
			rabbitmq.WithConsumerOptionsQueueExclusive,
			rabbitmq.WithConsumerOptionsQueueAutoDelete,
			rabbitmq.WithConsumerOptionsLogger(rabbitmq2.NewLogger(m.logger, logrus.ErrorLevel)),
		)
		if err != nil {
			return fmt.Errorf("cant start all instances consumer: %w", err)
		}

		errGroup.Go(func() error {
			return baConsumer.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
				return m.ProcessMsg(d, m.broadcastAllHandlers)
			})
		})

		m.broadcastAllConsumer = baConsumer
	}

	// start one instance consumer
	if len(m.broadcastHandlers) > 0 {
		bConsumer, err := rabbitmq.NewConsumer(conn, fmt.Sprintf("%s.%s.%s-queue", m.exchangeName, m.topic, m.serviceName),
			rabbitmq.WithConsumerOptionsExchangeName(m.exchangeName),
			rabbitmq.WithConsumerOptionsRoutingKey(fmt.Sprintf("%s.%s.*", m.exchangeName, m.topic)), // subscribe to all transmission types
			rabbitmq.WithConsumerOptionsExchangeDeclare,
			rabbitmq.WithConsumerOptionsExchangeKind("topic"),
			rabbitmq.WithConsumerOptionsQueueArgs(
				map[string]any{
					"x-expires": (24 * time.Hour).Milliseconds(), // auto remove in 24 hour if no connection
				},
			),
		)
		if err != nil {
			return fmt.Errorf("cant start single service consumer: %w", err)
		}

		errGroup.Go(func() error {
			return bConsumer.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
				return m.ProcessMsg(d, m.broadcastHandlers)
			})
		})

		m.broadcastAllConsumer = bConsumer
	}

	return errGroup.Wait()
}

func (m *MessageChannel) StopConsuming() {
	if m.broadcastConsumer != nil {
		m.broadcastConsumer.Close()
	}
	if m.broadcastAllConsumer != nil {
		m.broadcastAllConsumer.Close()
	}
}

func (m *MessageChannel) ProcessMsg(d rabbitmq.Delivery, handlers []handler) rabbitmq.Action {
	payload := proto.Clone(m.unmarshalPayload)

	if err := protojson.Unmarshal(d.Body, payload); err != nil {
		m.logger.Errorf("Failed to unmarshal proto message for topic '%s' (maybe wrong message payload is published): %v", m.topic, err)
	}

	if err := m.validator.Validate(payload); err != nil {
		m.logger.Warnf("Pubsub proto message validation error for topic '%v', ingored: %v", m.topic, err)
	}

	for _, h := range handlers {
		go h(payload, d)
	}

	return 0
}
