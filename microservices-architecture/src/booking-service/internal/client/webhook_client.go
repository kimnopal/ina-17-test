package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// WebhookClient interface for sending webhook notifications to other services
// TODO: [RECOMMENDATION] Replace webhook with Message Broker (RabbitMQ/Kafka)
// Benefits:
// - Guaranteed delivery with retry mechanism
// - Better scalability and decoupling
// - Event replay capability for debugging
// - Dead letter queue for failed messages
//
// Example with RabbitMQ:
// publisher.Publish("booking.events", BookingEvent{...})
//
// Example with Kafka:
// producer.Send("booking-events-topic", BookingMessage{...})
type WebhookClient interface {
	NotifyPaymentService(event string, bookingID uuid.UUID, status string) error
}

type webhookClient struct {
	paymentWebhookURL string
}

// BookingWebhookPayload represents the webhook payload sent to payment service
type BookingWebhookPayload struct {
	Event     string `json:"event"`
	BookingID string `json:"booking_id"`
	Status    string `json:"status"`
}

// NewWebhookClient creates a new webhook client instance
func NewWebhookClient() WebhookClient {
	paymentURL := os.Getenv("PAYMENT_WEBHOOK_URL")
	if paymentURL == "" {
		paymentURL = "http://localhost:3002/api/v1/payments/webhook/booking"
	}
	return &webhookClient{paymentWebhookURL: paymentURL}
}

// NotifyPaymentService sends a webhook notification to payment service
// TODO: [RECOMMENDATION] Replace with Message Broker for better reliability
// Example with RabbitMQ:
//
//	publisher.Publish("booking.events", BookingEvent{
//	    Event:     event,
//	    BookingID: bookingID,
//	    Status:    status,
//	})
//
// Example with Kafka:
//
//	producer.Send("booking-events-topic", BookingMessage{...})
func (c *webhookClient) NotifyPaymentService(event string, bookingID uuid.UUID, status string) error {
	payload := BookingWebhookPayload{
		Event:     event,
		BookingID: bookingID.String(),
		Status:    status,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %v", err)
	}

	req, err := http.NewRequest("POST", c.paymentWebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook to payment service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment service webhook returned status %d", resp.StatusCode)
	}

	return nil
}
