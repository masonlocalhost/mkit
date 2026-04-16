package pubsub

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/wagslane/go-rabbitmq"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type state int

const (
	stateIdle = iota
	stateRunning
)

// TransmissionType 2 modes: Broadcast to each service, BroadcastAll to all service instances
type TransmissionType int

const (
	Broadcast    TransmissionType = 1 // Broadcast to each service
	BroadcastAll TransmissionType = 2 // Broadcast to all service instances, include common queue of each service
)

// Service is used to publish and subscribe proto messages
type Service struct {
	serviceName     string
	logger          *slog.Logger
	conn            *rabbitmq.Conn
	subscriberMap   map[string]*MessageChannel
	ctx             context.Context
	state           state
	errGroup        errgroup.Group
	payloadRegistry map[string]proto.Message
	publisher       *rabbitmq.Publisher
	exchangeName    string
	sync.Mutex
}

func NewService(
	ctx context.Context, conn *rabbitmq.Conn, serviceName string, exchangeName string, logger *slog.Logger,
	payloadRegistry map[string]proto.Message,
) (*Service, error) {
	publisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		return nil, fmt.Errorf("cant create rabbitmq publisher: %w", err)
	}

	return &Service{
		serviceName:     serviceName,
		logger:          logger,
		conn:            conn,
		subscriberMap:   make(map[string]*MessageChannel),
		ctx:             ctx,
		state:           0,
		payloadRegistry: payloadRegistry,
		publisher:       publisher,
		exchangeName:    exchangeName,
	}, nil
}

func (s *Service) Subscribe(topic string, transmissionType TransmissionType, handler handler) error {
	if s.state != stateIdle {
		return fmt.Errorf("cannot add subscriber when pubsub service is running")
	}

	s.Lock()
	msgChannel, ok := s.subscriberMap[topic]
	if ok {
		msgChannel.AddHandler(handler, transmissionType)
	} else {
		if _, ok := s.payloadRegistry[topic]; !ok {
			return fmt.Errorf("no payload registry for topic '%s', please register it in 'service.payloadRegistry()'", topic)
		}

		mCh, err := NewMsgChannel(s.logger, s.serviceName, s.exchangeName, topic, s.payloadRegistry[topic])
		if err != nil {
			return fmt.Errorf("failed to create message channel for topic '%s': %w", topic, err)
		}
		s.subscriberMap[topic] = mCh
		s.subscriberMap[topic].AddHandler(handler, transmissionType)
	}
	s.Unlock()

	return nil
}

// Publish broadcasts to all consumers bound to the exchange with the given topic.
func (s *Service) Publish(ctx context.Context, topic string, payload proto.Message) error {
	msg, err := protojson.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal collection update event: %w", err)
	}

	if err := s.publisher.PublishWithContext(
		ctx, msg, []string{getRoutingKey(s.exchangeName, topic, BroadcastAll)}, rabbitmq.WithPublishOptionsExchange(s.exchangeName),
	); err != nil {
		return fmt.Errorf("failed to publish mq message: %w", err)
	}

	return nil
}

func (s *Service) Start() error {
	if s.state == stateRunning {
		s.logger.Error("pubsub service already running")
		return nil
	}
	group, ctx := errgroup.WithContext(s.ctx)

	s.Lock()
	for topic, msgChannel := range s.subscriberMap {
		group.Go(func() error {
			return msgChannel.StartConsuming(ctx, s.conn)
		})
		s.logger.Info("Started consuming", "topic", topic)
	}
	s.Unlock()

	s.state = stateRunning
	return group.Wait()
}

func (s *Service) Stop() {
	s.Lock()
	for _, msgChannel := range s.subscriberMap {
		msgChannel.StopConsuming()
	}
	s.Unlock()

	s.logger.Info("Stopped all pubsub consumers")
}
