package queue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"weave-module/config"
)

var Connection *amqp.Connection
var Channel *amqp.Channel

type QueueNames struct {
	NotificationQueue string
	EmailQueue        string
	AnalyticsQueue    string
	ProcessingQueue   string
}

var Queues = QueueNames{
	NotificationQueue: "notifications",
	EmailQueue:        "emails",
	AnalyticsQueue:    "analytics",
	ProcessingQueue:   "processing",
}

func Connect(cfg *config.Config) error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		cfg.Queue.User,
		cfg.Queue.Password,
		cfg.Queue.Host,
		cfg.Queue.Port,
		cfg.Queue.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	Connection = conn
	Channel = ch

	// Declare all queues
	err = declareQueues()
	if err != nil {
		return fmt.Errorf("failed to declare queues: %w", err)
	}

	return nil
}

func declareQueues() error {
	queues := []string{
		Queues.NotificationQueue,
		Queues.EmailQueue,
		Queues.AnalyticsQueue,
		Queues.ProcessingQueue,
	}

	for _, queueName := range queues {
		_, err := Channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
	}

	return nil
}

func Close() error {
	if Channel != nil {
		Channel.Close()
	}
	if Connection != nil {
		return Connection.Close()
	}
	return nil
}

// Message types
type NotificationMessage struct {
	UserID  string      `json:"user_id"`
	Type    string      `json:"type"`
	Title   string      `json:"title"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type EmailMessage struct {
	To      string            `json:"to"`
	Subject string            `json:"subject"`
	Body    string            `json:"body"`
	Data    map[string]string `json:"data,omitempty"`
}

type AnalyticsMessage struct {
	Event     string      `json:"event"`
	UserID    string      `json:"user_id,omitempty"`
	WeaveID   string      `json:"weave_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

type ProcessingMessage struct {
	Type    string      `json:"type"`
	WeaveID string      `json:"weave_id,omitempty"`
	UserID  string      `json:"user_id,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Publisher functions
func PublishNotification(msg NotificationMessage) error {
	return publish(Queues.NotificationQueue, msg)
}

func PublishEmail(msg EmailMessage) error {
	return publish(Queues.EmailQueue, msg)
}

func PublishAnalytics(msg AnalyticsMessage) error {
	return publish(Queues.AnalyticsQueue, msg)
}

func PublishProcessing(msg ProcessingMessage) error {
	return publish(Queues.ProcessingQueue, msg)
}

func publish(queueName string, message interface{}) error {
	if Channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = Channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message to queue %s: %w", queueName, err)
	}

	return nil
}

// Consumer functions
func ConsumeNotifications(handler func(NotificationMessage) error) error {
	return consume(Queues.NotificationQueue, func(body []byte) error {
		var msg NotificationMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return err
		}
		return handler(msg)
	})
}

func ConsumeEmails(handler func(EmailMessage) error) error {
	return consume(Queues.EmailQueue, func(body []byte) error {
		var msg EmailMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return err
		}
		return handler(msg)
	})
}

func ConsumeAnalytics(handler func(AnalyticsMessage) error) error {
	return consume(Queues.AnalyticsQueue, func(body []byte) error {
		var msg AnalyticsMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return err
		}
		return handler(msg)
	})
}

func ConsumeProcessing(handler func(ProcessingMessage) error) error {
	return consume(Queues.ProcessingQueue, func(body []byte) error {
		var msg ProcessingMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return err
		}
		return handler(msg)
	})
}

func consume(queueName string, handler func([]byte) error) error {
	if Channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	msgs, err := Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer for queue %s: %w", queueName, err)
	}

	go func() {
		for d := range msgs {
			err := handler(d.Body)
			if err != nil {
				log.Printf("Error processing message from queue %s: %v", queueName, err)
				d.Nack(false, true) // Requeue on error
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}