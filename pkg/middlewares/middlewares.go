package middlewares

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"risq_backend/pkg/jwt"
	"risq_backend/pkg/response"
)

func JWT(jwtService *jwt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Missing authorization header")
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Unauthorized(c, "Invalid authorization header format")
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return response.Unauthorized(c, "Missing token")
		}

		// Validate the token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			return response.Unauthorized(c, "Invalid or expired token")
		}

		// Add user info to context
		c.Locals("user_id", claims.UserID.String())
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// StartupContext middleware ensures user has completed startup onboarding
// For now, we'll delegate this validation to the controllers to avoid circular dependencies
// Each controller that requires startup context will validate startup ownership directly
func StartupContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from previous JWT middleware
		userIDInterface := c.Locals("user_id")
		if userIDInterface == nil {
			return response.Unauthorized(c, "User not authenticated")
		}

		// Validate user ID format
		_, err := uuid.Parse(userIDInterface.(string))
		if err != nil {
			return response.BadRequest(c, "Invalid user ID", err)
		}

		// Note: Actual startup validation happens in controllers
		// This middleware just ensures JWT authentication is present
		return c.Next()
	}
}

func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}

		return c.Next()
	}
}

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Simple request logging
		method := c.Method()
		path := c.Path()
		ip := c.IP()

		// Process request
		err := c.Next()

		// Log after processing
		status := c.Response().StatusCode()
		fmt.Printf("[%s] %s %s - %d\n", ip, method, path, status)

		return err
	}
}
