package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"order-service/internal/config"
	"time"

	"github.com/nats-io/nats.go"
)

type EventPublisher interface {
	PublishOrderCreated(event OrderCreatedEvent) error
	Close()
}

type NATSPublisher struct {
	conn *nats.Conn
}

func NewNATSPublisher(cfg *config.Config) (*NATSPublisher, error) {
	startTime := time.Now()
	log.Printf("[%s] Connecting to NATS at %s",
		startTime.Format(time.RFC3339Nano), cfg.NATS.URL)

	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		log.Printf("[%s] Failed to connect to NATS: %v [latency: %v]",
			time.Now().Format(time.RFC3339Nano), err, time.Since(startTime))
		return nil, err
	}

	log.Printf("[%s] Successfully connected to NATS [latency: %v]",
		time.Now().Format(time.RFC3339Nano), time.Since(startTime))
	return &NATSPublisher{
		conn: nc,
	}, nil
}

func (p *NATSPublisher) PublishOrderCreated(event OrderCreatedEvent) error {
	startTime := time.Now()
	log.Printf("[%s] Preparing to publish order.created event for order ID: %s with %d items",
		startTime.Format(time.RFC3339Nano), event.OrderID, len(event.Items))

	if event.Timestamp == 0 {
		event.Timestamp = time.Now().UnixNano()
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("[%s] Error marshalling order created event: %v [latency: %v]",
			time.Now().Format(time.RFC3339Nano), err, time.Since(startTime))
		return err
	}

	messageSizeKB := float64(len(eventBytes)) / 1024.0

	publishTime := time.Now()
	err = p.conn.Publish(SubjectOrderCreated, eventBytes)
	if err != nil {
		log.Printf("[%s] Error publishing order created event: %v [latency: %v]",
			time.Now().Format(time.RFC3339Nano), err, time.Since(startTime))
		return err
	}

	log.Printf("[%s] Published order.created event for order ID: %s (%.2f KB) to subject: %s [latency: %v]",
		time.Now().Format(time.RFC3339Nano), event.OrderID, messageSizeKB,
		SubjectOrderCreated, time.Since(startTime))

	fmt.Printf("EVENT_PUBLISHED,order_id=%s,timestamp=%d,items=%d,size=%.2fKB,latency_ms=%.2f\n",
		event.OrderID, publishTime.UnixNano(), len(event.Items),
		messageSizeKB, float64(time.Since(startTime).Microseconds())/1000.0)

	return nil
}

func (p *NATSPublisher) Close() {
	closeTime := time.Now()

	if p.conn != nil {
		p.conn.Close()
		log.Printf("[%s] NATS publisher connection closed",
			closeTime.Format(time.RFC3339Nano))
	}
}
