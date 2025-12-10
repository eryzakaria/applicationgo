package service

import (
	"context"
	"errors"

	"suitemedia/internal/models"
	"suitemedia/internal/repository"
	"suitemedia/pkg/redis"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	List(ctx context.Context, params models.ListParams) ([]*models.UserResponse, int64, error)
	GetByID(ctx context.Context, id string) (*models.UserResponse, error)
	Create(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error)
	Update(ctx context.Context, id string, req models.UpdateUserRequest) (*models.UserResponse, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
	redis    *redis.Client
}

func NewUserService(userRepo repository.UserRepository, redis *redis.Client) UserService {
	return &userService{
		userRepo: userRepo,
		redis:    redis,
	}
}

func (s *userService) List(ctx context.Context, params models.ListParams) ([]*models.UserResponse, int64, error) {
	users, total, err := s.userRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		resp := user.ToResponse()
		responses[i] = &resp
	}

	return responses, total, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *userService) Create(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	// Check if email exists
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *userService) Update(ctx context.Context, id string, req models.UpdateUserRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(ctx, id)
}
