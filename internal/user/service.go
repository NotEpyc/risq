package user

import (
	"context"
	"fmt"

	"risq_backend/pkg/auth"
	"risq_backend/pkg/logger"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, email, name, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, email, name, password string) (*User, error) {
	logger.Infof("Creating user with email: %s", email)

	// Check if user already exists
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		return nil, fmt.Errorf("failed to process password")
	}

	user := &User{
		ID:           uuid.New(),
		Email:        email,
		Name:         name,
		PasswordHash: hashedPassword,
		Role:         "founder",
	}

	if err := s.repo.Create(ctx, user); err != nil {
		logger.Errorf("Failed to create user: %v", err)
		return nil, err
	}

	logger.Infof("Successfully created user: %s", user.ID)
	return user, nil
}

func (s *service) Login(ctx context.Context, email, password string) (*User, error) {
	logger.Infof("Attempting login for email: %s", email)

	// Get user by email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		logger.Errorf("Failed to get user by email: %v", err)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Compare password
	if err := auth.ComparePassword(user.PasswordHash, password); err != nil {
		logger.Errorf("Invalid password for user: %s", email)
		return nil, fmt.Errorf("invalid credentials")
	}

	logger.Infof("Successfully authenticated user: %s", user.ID)
	return user, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	logger.Debugf("Getting user by ID: %s", id)

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("Failed to get user by ID: %v", err)
		return nil, err
	}

	return user, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
	logger.Debugf("Getting user by email: %s", email)

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		logger.Errorf("Failed to get user by email: %v", err)
		return nil, err
	}

	return user, nil
}

func (s *service) Update(ctx context.Context, user *User) error {
	logger.Infof("Updating user: %s", user.ID)

	if err := s.repo.Update(ctx, user); err != nil {
		logger.Errorf("Failed to update user: %v", err)
		return err
	}

	logger.Infof("Successfully updated user: %s", user.ID)
	return nil
}
