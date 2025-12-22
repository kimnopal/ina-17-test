package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"user-service/internal/model"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(username, password string) (*model.UserResponse, error)
	Login(username, password string) (string, error)
	GetAuthenticatedUser(token string) (*model.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(username, password string) (*model.UserResponse, error) {
	// Check if username already exists
	existingUser, _ := s.repo.FindByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
	}, nil
}

func (s *userService) Login(username, password string) (string, error) {
	// Find user by username
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate token (simple UUID-based token)
	byteString := make([]byte, 100)
	rand.Read(byteString)
	token := fmt.Sprintf("%x", byteString)

	// Update user token in database
	user.Token = token
	if err := s.repo.Update(user); err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *userService) GetAuthenticatedUser(token string) (*model.UserResponse, error) {
	if token == "" {
		return nil, errors.New("authorization token is required")
	}

	user, err := s.repo.FindByToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	return &model.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
	}, nil
}
