package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type UserClient interface {
	GetAuthenticatedUser(authToken string) (*UserResponse, error)
}

type userClient struct {
	baseURL string
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserServiceResponse struct {
	Message string       `json:"message"`
	Data    UserResponse `json:"data"`
}

func NewUserClient() UserClient {
	baseURL := os.Getenv("USER_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	return &userClient{baseURL: baseURL}
}

func (c *userClient) GetAuthenticatedUser(authToken string) (*UserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/auth", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned status %d: %s", resp.StatusCode, string(body))
	}

	var userResp UserServiceResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, err
	}

	if userResp.Data.ID == "" {
		return nil, errors.New("invalid user data received")
	}

	return &userResp.Data, nil
}
