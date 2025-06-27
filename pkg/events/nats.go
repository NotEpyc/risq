package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"risq_backend/pkg/logger"

	"github.com/nats-io/nats.go"
)

// NATSPublisher handles event publishing
type NATSPublisher struct {
	conn *nats.Conn
}

// NATSSubscriber handles event subscription
type NATSSubscriber struct {
	conn *nats.Conn
}

// EventHandler represents a function that handles events
type EventHandler func(ctx context.Context, data []byte) error

// NewNATSPublisher creates a new NATS publisher
func NewNATSPublisher(natsURL string) (*NATSPublisher, error) {
	opts := []nats.Option{
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Warnf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Warn("NATS connection closed")
		}),
	}

	conn, err := nats.Connect(natsURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	logger.Infof("Connected to NATS server at %s", natsURL)

	return &NATSPublisher{conn: conn}, nil
}

// NewNATSSubscriber creates a new NATS subscriber
func NewNATSSubscriber(natsURL string) (*NATSSubscriber, error) {
	opts := []nats.Option{
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Warnf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Warn("NATS connection closed")
		}),
	}

	conn, err := nats.Connect(natsURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	logger.Infof("Connected to NATS server at %s", natsURL)

	return &NATSSubscriber{conn: conn}, nil
}

// Publish publishes an event to a subject
func (p *NATSPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	if err := p.conn.Publish(subject, payload); err != nil {
		return fmt.Errorf("failed to publish event to subject %s: %w", subject, err)
	}

	logger.Debugf("Published event to subject: %s", subject)
	return nil
}

// PublishAndWait publishes an event and waits for acknowledgment
func (p *NATSPublisher) PublishAndWait(ctx context.Context, subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	if err := p.conn.Publish(subject, payload); err != nil {
		return fmt.Errorf("failed to publish event to subject %s: %w", subject, err)
	}

	// Flush to ensure message is sent
	if err := p.conn.Flush(); err != nil {
		return fmt.Errorf("failed to flush NATS connection: %w", err)
	}

	logger.Debugf("Published and flushed event to subject: %s", subject)
	return nil
}

// Subscribe subscribes to a subject with a handler
func (s *NATSSubscriber) Subscribe(subject string, handler EventHandler) (*nats.Subscription, error) {
	sub, err := s.conn.Subscribe(subject, func(msg *nats.Msg) {
		ctx := context.Background()

		logger.Debugf("Received event on subject: %s", subject)

		if err := handler(ctx, msg.Data); err != nil {
			logger.Errorf("Error handling event on subject %s: %v", subject, err)
			return
		}

		logger.Debugf("Successfully handled event on subject: %s", subject)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}

	logger.Infof("Subscribed to subject: %s", subject)
	return sub, nil
}

// SubscribeQueue subscribes to a subject with queue group for load balancing
func (s *NATSSubscriber) SubscribeQueue(subject, queue string, handler EventHandler) (*nats.Subscription, error) {
	sub, err := s.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		ctx := context.Background()

		logger.Debugf("Received queued event on subject: %s, queue: %s", subject, queue)

		if err := handler(ctx, msg.Data); err != nil {
			logger.Errorf("Error handling queued event on subject %s, queue %s: %v", subject, queue, err)
			return
		}

		logger.Debugf("Successfully handled queued event on subject: %s, queue: %s", subject, queue)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %s with queue %s: %w", subject, queue, err)
	}

	logger.Infof("Subscribed to subject: %s with queue: %s", subject, queue)
	return sub, nil
}

// Close closes the NATS connections
func (p *NATSPublisher) Close() {
	if p.conn != nil {
		p.conn.Close()
		logger.Info("NATS publisher connection closed")
	}
}

func (s *NATSSubscriber) Close() {
	if s.conn != nil {
		s.conn.Close()
		logger.Info("NATS subscriber connection closed")
	}
}

// Health check for NATS connection
func (p *NATSPublisher) IsConnected() bool {
	return p.conn != nil && p.conn.IsConnected()
}

func (s *NATSSubscriber) IsConnected() bool {
	return s.conn != nil && s.conn.IsConnected()
}
