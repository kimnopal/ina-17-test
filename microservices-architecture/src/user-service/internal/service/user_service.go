package service

import (
	"errors"
	"user-service/internal/auth"
	"user-service/internal/model"
	"user-service/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type UserService interface {
	CreateUser(username, password string) (*model.UserResponse, error)
	Login(username, password string) (*LoginResponse, error)
	GetAuthenticatedUser(userID string) (*model.UserResponse, error)
	RefreshToken(refreshToken string) (*RefreshResponse, error)
	Logout(refreshToken string) error
}

type userService struct {
	repo         repository.UserRepository
	refreshRepo  repository.RefreshTokenRepository
}

func NewUserService(repo repository.UserRepository, refreshRepo repository.RefreshTokenRepository) UserService {
	return &userService{
		repo:        repo,
		refreshRepo: refreshRepo,
	}
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

func (s *userService) Login(username, password string) (*LoginResponse, error) {
	// Find user by username
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate access token (JWT)
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Generate refresh token
	refreshTokenStr, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Store refresh token in database
	refreshToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: auth.GetRefreshTokenExpiry(),
	}

	if err := s.refreshRepo.Create(refreshToken); err != nil {
		return nil, errors.New("failed to store refresh token")
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    auth.GetAccessTokenExpirySeconds(),
	}, nil
}

func (s *userService) GetAuthenticatedUser(userID string) (*model.UserResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.FindByID(parsedID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &model.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
	}, nil
}

func (s *userService) RefreshToken(refreshTokenStr string) (*RefreshResponse, error) {
	if refreshTokenStr == "" {
		return nil, errors.New("refresh token is required")
	}

	// Find refresh token in database
	refreshToken, err := s.refreshRepo.FindByToken(refreshTokenStr)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if refresh token is expired
	if refreshToken.IsExpired() {
		// Delete expired token
		s.refreshRepo.DeleteByToken(refreshTokenStr)
		return nil, errors.New("refresh token has expired")
	}

	// Get user
	user, err := s.repo.FindByID(refreshToken.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new access token
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	return &RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   auth.GetAccessTokenExpirySeconds(),
	}, nil
}

func (s *userService) Logout(refreshTokenStr string) error {
	if refreshTokenStr == "" {
		return errors.New("refresh token is required")
	}

	// Delete refresh token from database
	err := s.refreshRepo.DeleteByToken(refreshTokenStr)
	if err != nil {
		return errors.New("failed to logout")
	}

	return nil
}
