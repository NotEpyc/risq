package controller

import (
	"strings"
	"time"

	"risq_backend/internal/user"
	"risq_backend/pkg/jwt"
	"risq_backend/pkg/logger"
	"risq_backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService user.Service
	jwtService  *jwt.Service
}

func NewUserController(userService user.Service, jwtService *jwt.Service) *UserController {
	return &UserController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *UserController) Create(ctx *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Name     string `json:"name" validate:"required,min=2,max=100"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		return response.BadRequest(ctx, "Invalid input", err)
	}

	// Basic input validation
	if input.Email == "" {
		return response.BadRequest(ctx, "Email is required", nil)
	}
	if input.Name == "" {
		return response.BadRequest(ctx, "Name is required", nil)
	}
	if len(input.Password) < 8 {
		return response.BadRequest(ctx, "Password must be at least 8 characters", nil)
	}

	user, err := c.userService.Create(ctx.Context(), input.Email, input.Name, input.Password)
	if err != nil {
		logger.Errorf("Failed to create user: %v", err)
		if strings.Contains(err.Error(), "already exists") {
			return response.BadRequest(ctx, "User with this email already exists", err)
		}
		return response.InternalError(ctx, "Failed to create user", err)
	}

	// Generate JWT token for the new user
	token, err := c.jwtService.GenerateToken(user.ID, user.Email, user.Role, 24*time.Hour)
	if err != nil {
		logger.Errorf("Failed to generate token: %v", err)
		return response.InternalError(ctx, "Failed to generate authentication token", err)
	}

	return response.Success(ctx, "User created successfully", fiber.Map{
		"user":  user,
		"token": token,
	})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		return response.BadRequest(ctx, "Invalid input", err)
	}

	// Basic input validation
	if input.Email == "" {
		return response.BadRequest(ctx, "Email is required", nil)
	}
	if input.Password == "" {
		return response.BadRequest(ctx, "Password is required", nil)
	}

	user, err := c.userService.Login(ctx.Context(), input.Email, input.Password)
	if err != nil {
		logger.Errorf("Failed to authenticate user: %v", err)
		return response.Unauthorized(ctx, "Invalid credentials")
	}

	// Generate JWT token
	token, err := c.jwtService.GenerateToken(user.ID, user.Email, user.Role, 24*time.Hour)
	if err != nil {
		logger.Errorf("Failed to generate token: %v", err)
		return response.InternalError(ctx, "Failed to generate authentication token", err)
	}

	return response.Success(ctx, "Login successful", fiber.Map{
		"user":  user,
		"token": token,
	})
}
