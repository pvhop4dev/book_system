package user_service

import (
	"book_system/internal/model"
	repo "book_system/internal/repository"
	"book_system/internal/service"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo     repo.IUserRepository
	tokenService service.ITokenService
}

// NewUserService creates a new instance of user service
func NewUserService(
	userRepo repo.IUserRepository,
	tokenService service.ITokenService,
) service.IUserService {
	return &userService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Check if user with email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user entity
	now := time.Now()
	user := &model.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		Role:      req.Role,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Convert to DTO and return
	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	users, total, err := s.userRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Convert users to DTOs
	userDTOs := make([]*model.User, len(users))
	for i, user := range users {
		userDTOs[i] = user
	}

	return userDTOs, total, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Get existing user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	return s.userRepo.Delete(ctx, userID)
}

func (s *userService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	tokenPair, err := s.tokenService.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	// Update last login
	user.LastLogin = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		// Log the error but don't fail the login
		// You might want to add proper logging here
	}

	return &model.LoginResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         user.ToDTO(),
	}, nil
}

func (s *userService) Register(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error) {
	// Create user
	createReq := &model.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     "user", // Default role for new users
	}

	user, err := s.CreateUser(ctx, createReq)
	if err != nil {
		return nil, err
	}

	// Generate token for immediate login
	loginReq := &model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	loginRes, err := s.Login(ctx, loginReq)
	if err != nil {
		// This shouldn't happen since we just created the user
		return nil, err
	}

	return &model.RegisterResponse{
		User:  user,
		Token: loginRes.Token,
	}, nil
}

func (s *userService) RefreshToken(ctx context.Context, token string) (*model.LoginResponse, error) {
	// Validate refresh token
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, uuid.MustParse(claims.UserID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	tokenPair, err := s.tokenService.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         user.ToDTO(),
	}, nil
}
