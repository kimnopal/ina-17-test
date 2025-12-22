package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type BookingClient interface {
	GetBookingByID(bookingID uuid.UUID) (*BookingResponse, error)
}

type bookingClient struct {
	baseURL string
}

type BookingResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	EventID     uuid.UUID  `json:"event_id"`
	TicketID    uuid.UUID  `json:"ticket_id"`
	Quantity    int        `json:"quantity"`
	TotalAmount float64    `json:"total_amount"`
	Status      string     `json:"status"`
	ExpiredAt   *time.Time `json:"expired_at"`
}

type BookingServiceResponse struct {
	Message string          `json:"message"`
	Data    BookingResponse `json:"data"`
}

func NewBookingClient() BookingClient {
	baseURL := os.Getenv("BOOKING_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3001"
	}
	return &bookingClient{baseURL: baseURL}
}

func (c *bookingClient) GetBookingByID(bookingID uuid.UUID) (*BookingResponse, error) {
	url := fmt.Sprintf("%s/api/v1/bookings/%s", c.baseURL, bookingID.String())

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to booking service: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("booking service returned status %d: %s", resp.StatusCode, string(body))
	}

	var bookingResp BookingServiceResponse
	if err := json.Unmarshal(body, &bookingResp); err != nil {
		return nil, err
	}

	return &bookingResp.Data, nil
}
