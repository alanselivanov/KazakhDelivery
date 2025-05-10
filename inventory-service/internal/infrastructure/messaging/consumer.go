package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"inventory-service/internal/config"

	"github.com/nats-io/nats.go"
)

type MessageHandler func(context.Context, *OrderCreatedEvent) error

type EventConsumer interface {
	SubscribeToOrderCreated(handler MessageHandler) error
	PublishToDLQ(subject string, event []byte) error
	Close()
}

type NATSConsumer struct {
	conn         *nats.Conn
	subscription *nats.Subscription
}

func NewNATSConsumer(cfg *config.Config) (*NATSConsumer, error) {
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
	return &NATSConsumer{
		conn: nc,
	}, nil
}

func (c *NATSConsumer) SubscribeToOrderCreated(handler MessageHandler) error {
	startTime := time.Now()
	log.Printf("[%s] Subscribing to subject: %s",
		startTime.Format(time.RFC3339Nano), SubjectOrderCreated)

	sub, err := c.conn.Subscribe(SubjectOrderCreated, func(msg *nats.Msg) {
		receiveTime := time.Now()
		log.Printf("[%s] Received message from subject: %s (%.2f KB)",
			receiveTime.Format(time.RFC3339Nano),
			msg.Subject, float64(len(msg.Data))/1024.0)

		var event OrderCreatedEvent
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			log.Printf("[%s] Error unmarshalling message: %v",
				time.Now().Format(time.RFC3339Nano), err)

			fmt.Printf("EVENT_UNMARSHALLING_ERROR,subject=%s,timestamp=%d,error=%s\n",
				msg.Subject, receiveTime.UnixNano(), err.Error())

			c.PublishToDLQ(SubjectDeadLetter, msg.Data)
			return
		}

		msgLatency := float64(receiveTime.UnixNano()-event.Timestamp) / 1e6

		log.Printf("[%s] Received order.created event for order ID: %s with %d items (message latency: %.2f ms)",
			receiveTime.Format(time.RFC3339Nano), event.OrderID, len(event.Items), msgLatency)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Printf("EVENT_RECEIVED,order_id=%s,timestamp=%d,items=%d,latency_ms=%.2f\n",
			event.OrderID, receiveTime.UnixNano(), len(event.Items), msgLatency)

		procStartTime := time.Now()
		err = handler(ctx, &event)
		procDuration := time.Since(procStartTime)

		if err != nil {
			log.Printf("[%s] Error handling order created event: %v [processing_time: %v]",
				time.Now().Format(time.RFC3339Nano), err, procDuration)

			fmt.Printf("EVENT_PROCESSING_ERROR,order_id=%s,timestamp=%d,error=%s,proc_time_ms=%.2f\n",
				event.OrderID, time.Now().UnixNano(), err.Error(), float64(procDuration.Microseconds())/1000.0)

			c.PublishToDLQ(SubjectDeadLetter, msg.Data)
		} else {
			log.Printf("[%s] Successfully processed order created event for order ID: %s [processing_time: %v]",
				time.Now().Format(time.RFC3339Nano), event.OrderID, procDuration)

			fmt.Printf("EVENT_PROCESSED,order_id=%s,timestamp=%d,proc_time_ms=%.2f\n",
				event.OrderID, time.Now().UnixNano(), float64(procDuration.Microseconds())/1000.0)
		}
	})

	if err != nil {
		log.Printf("[%s] Error subscribing to subject %s: %v [latency: %v]",
			time.Now().Format(time.RFC3339Nano), SubjectOrderCreated, err, time.Since(startTime))
		return err
	}

	log.Printf("[%s] Successfully subscribed to subject: %s [latency: %v]",
		time.Now().Format(time.RFC3339Nano), SubjectOrderCreated, time.Since(startTime))

	c.subscription = sub
	return nil
}

func (c *NATSConsumer) PublishToDLQ(subject string, event []byte) error {
	startTime := time.Now()
	log.Printf("[%s] Publishing message to DLQ: %s (%.2f KB)",
		startTime.Format(time.RFC3339Nano), subject, float64(len(event))/1024.0)

	err := c.conn.Publish(subject, event)
	if err != nil {
		log.Printf("[%s] Error publishing to DLQ: %v [latency: %v]",
			time.Now().Format(time.RFC3339Nano), err, time.Since(startTime))
		return err
	}

	log.Printf("[%s] Successfully published message to DLQ: %s [latency: %v]",
		time.Now().Format(time.RFC3339Nano), subject, time.Since(startTime))

	fmt.Printf("DLQ_PUBLISHED,subject=%s,timestamp=%d,size=%.2fKB,latency_ms=%.2f\n",
		subject, startTime.UnixNano(), float64(len(event))/1024.0,
		float64(time.Since(startTime).Microseconds())/1000.0)

	return nil
}

func (c *NATSConsumer) Close() {
	closeTime := time.Now()

	if c.subscription != nil {
		c.subscription.Unsubscribe()
		log.Printf("[%s] Unsubscribed from NATS", closeTime.Format(time.RFC3339Nano))
	}

	if c.conn != nil {
		c.conn.Close()
		log.Printf("[%s] NATS consumer connection closed", closeTime.Format(time.RFC3339Nano))
	}
}
