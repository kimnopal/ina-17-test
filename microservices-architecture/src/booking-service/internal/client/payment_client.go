package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type PaymentClient interface {
	CreatePayment(bookingID uint, amount float64) (*PaymentResponse, error)
	GetPaymentStatus(paymentID string) (*PaymentResponse, error)
}

type paymentClient struct {
	baseURL string
}

type PaymentResponse struct {
	ID        string  `json:"id"`
	BookingID uint    `json:"booking_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
}

type PaymentRequest struct {
	BookingID uint    `json:"booking_id"`
	Amount    float64 `json:"amount"`
}

func NewPaymentClient() PaymentClient {
	baseURL := os.Getenv("PAYMENT_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3002"
	}
	return &paymentClient{baseURL: baseURL}
}

func (c *paymentClient) CreatePayment(bookingID uint, amount float64) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments", c.baseURL)

	reqBody := PaymentRequest{
		BookingID: bookingID,
		Amount:    amount,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("payment service returned status %d: %s", resp.StatusCode, string(body))
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, err
	}

	return &paymentResp, nil
}

func (c *paymentClient) GetPaymentStatus(paymentID string) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments/%s", c.baseURL, paymentID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment service returned status %d: %s", resp.StatusCode, string(body))
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, err
	}

	return &paymentResp, nil
}
