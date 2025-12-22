package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type WebhookClient interface {
	NotifyBookingService(event string, paymentID uuid.UUID, bookingID uuid.UUID) error
}

type webhookClient struct {
	bookingWebhookURL string
}

type PaymentWebhookPayload struct {
	Event     string `json:"event"`
	PaymentID string `json:"payment_id"`
	BookingID string `json:"booking_id"`
}

func NewWebhookClient() WebhookClient {
	bookingURL := os.Getenv("BOOKING_WEBHOOK_URL")
	if bookingURL == "" {
		bookingURL = "http://localhost:3001/api/v1/bookings/webhook/payment"
	}
	return &webhookClient{bookingWebhookURL: bookingURL}
}

func (c *webhookClient) NotifyBookingService(event string, paymentID uuid.UUID, bookingID uuid.UUID) error {
	payload := PaymentWebhookPayload{
		Event:     event,
		PaymentID: paymentID.String(),
		BookingID: bookingID.String(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %v", err)
	}

	req, err := http.NewRequest("POST", c.bookingWebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook to booking service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("booking service webhook returned status %d", resp.StatusCode)
	}

	return nil
}
