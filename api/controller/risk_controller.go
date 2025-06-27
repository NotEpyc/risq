package controller

import (
	"risq_backend/internal/risk"
	"risq_backend/internal/startup"
	"risq_backend/pkg/logger"
	"risq_backend/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RiskController struct {
	riskService    risk.Service
	startupService startup.Service
}

func NewRiskController(riskService risk.Service, startupService startup.Service) *RiskController {
	return &RiskController{
		riskService:    riskService,
		startupService: startupService,
	}
}

func (c *RiskController) GetCurrentRisk(ctx *fiber.Ctx) error {
	// Get user ID from context
	userIDInterface := ctx.Locals("user_id")
	if userIDInterface == nil {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	userID, err := uuid.Parse(userIDInterface.(string))
	if err != nil {
		return response.BadRequest(ctx, "Invalid user ID", err)
	}

	// Get user's startup
	startup, err := c.startupService.GetByUserID(ctx.Context(), userID)
	if err != nil {
		logger.Errorf("Failed to get user's startup: %v", err)
		return response.BadRequest(ctx, "User has no startup. Please complete startup onboarding first.", err)
	}

	riskProfile, err := c.riskService.GetCurrentRisk(ctx.Context(), startup.ID)
	if err != nil {
		logger.Errorf("Failed to get current risk: %v", err)
		return response.NotFound(ctx, "Risk profile not found")
	}

	return response.Success(ctx, "Risk profile retrieved successfully", riskProfile)
}

func (c *RiskController) GetRiskEvolution(ctx *fiber.Ctx) error {
	// Get user ID from context
	userIDInterface := ctx.Locals("user_id")
	if userIDInterface == nil {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	userID, err := uuid.Parse(userIDInterface.(string))
	if err != nil {
		return response.BadRequest(ctx, "Invalid user ID", err)
	}

	// Get user's startup
	startup, err := c.startupService.GetByUserID(ctx.Context(), userID)
	if err != nil {
		logger.Errorf("Failed to get user's startup: %v", err)
		return response.BadRequest(ctx, "User has no startup. Please complete startup onboarding first.", err)
	}

	// Get limit from query parameter
	limit := ctx.QueryInt("limit", 10)

	history, err := c.riskService.GetEvolutionHistory(ctx.Context(), startup.ID, limit)
	if err != nil {
		logger.Errorf("Failed to get risk evolution: %v", err)
		return response.InternalError(ctx, "Failed to retrieve risk evolution", err)
	}

	return response.Success(ctx, "Risk evolution retrieved successfully", history)
}
