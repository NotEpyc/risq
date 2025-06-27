package controller

import (
	"risq_backend/internal/decision"
	"risq_backend/internal/startup"
	"risq_backend/pkg/logger"
	"risq_backend/pkg/response"
	"risq_backend/types"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DecisionController struct {
	decisionService decision.Service
	startupService  startup.Service
}

func NewDecisionController(decisionService decision.Service, startupService startup.Service) *DecisionController {
	return &DecisionController{
		decisionService: decisionService,
		startupService:  startupService,
	}
}

func (c *DecisionController) SpeculateDecision(ctx *fiber.Ctx) error {
	logger.Info("Received decision speculation request")

	var input types.DecisionInput
	if err := ctx.BodyParser(&input); err != nil {
		return response.BadRequest(ctx, "Invalid input", err)
	}

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

	// Set startup ID in input
	input.StartupID = startup.ID

	result, err := c.decisionService.Speculate(ctx.Context(), input)
	if err != nil {
		logger.Errorf("Failed to speculate decision: %v", err)
		return response.InternalError(ctx, "Decision speculation failed", err)
	}

	return response.Success(ctx, "Decision speculation completed", result)
}

func (c *DecisionController) ConfirmDecision(ctx *fiber.Ctx) error {
	logger.Info("Received decision confirmation request")

	var input decision.DecisionConfirmInput
	if err := ctx.BodyParser(&input); err != nil {
		return response.BadRequest(ctx, "Invalid input", err)
	}

	// Get user ID from context
	userIDInterface := ctx.Locals("user_id")
	if userIDInterface == nil {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	userID, err := uuid.Parse(userIDInterface.(string))
	if err != nil {
		return response.BadRequest(ctx, "Invalid user ID", err)
	}

	// Get user's startup to validate ownership
	startup, err := c.startupService.GetByUserID(ctx.Context(), userID)
	if err != nil {
		logger.Errorf("Failed to get user's startup: %v", err)
		return response.BadRequest(ctx, "User has no startup. Please complete startup onboarding first.", err)
	}

	// Validate decision belongs to user's startup
	existingDecision, err := c.decisionService.GetByID(ctx.Context(), input.DecisionID)
	if err != nil {
		logger.Errorf("Failed to get decision: %v", err)
		return response.NotFound(ctx, "Decision not found")
	}

	if existingDecision.StartupID != startup.ID {
		return response.Unauthorized(ctx, "Decision does not belong to your startup")
	}

	decision, err := c.decisionService.Confirm(ctx.Context(), input)
	if err != nil {
		logger.Errorf("Failed to confirm decision: %v", err)
		return response.InternalError(ctx, "Decision confirmation failed", err)
	}

	return response.Success(ctx, "Decision confirmed successfully", decision)
}

func (c *DecisionController) GetDecisions(ctx *fiber.Ctx) error {
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

	decisions, err := c.decisionService.GetByStartupID(ctx.Context(), startup.ID)
	if err != nil {
		logger.Errorf("Failed to get decisions: %v", err)
		return response.InternalError(ctx, "Failed to retrieve decisions", err)
	}

	return response.Success(ctx, "Decisions retrieved successfully", decisions)
}

func (c *DecisionController) GetDecision(ctx *fiber.Ctx) error {
	decisionIDStr := ctx.Params("id")
	decisionID, err := uuid.Parse(decisionIDStr)
	if err != nil {
		return response.BadRequest(ctx, "Invalid decision ID", err)
	}

	// Get user ID from context
	userIDInterface := ctx.Locals("user_id")
	if userIDInterface == nil {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	userID, err := uuid.Parse(userIDInterface.(string))
	if err != nil {
		return response.BadRequest(ctx, "Invalid user ID", err)
	}

	// Get user's startup to validate ownership
	startup, err := c.startupService.GetByUserID(ctx.Context(), userID)
	if err != nil {
		logger.Errorf("Failed to get user's startup: %v", err)
		return response.BadRequest(ctx, "User has no startup. Please complete startup onboarding first.", err)
	}

	decision, err := c.decisionService.GetByID(ctx.Context(), decisionID)
	if err != nil {
		logger.Errorf("Failed to get decision: %v", err)
		return response.NotFound(ctx, "Decision not found")
	}

	// Validate decision belongs to user's startup
	if decision.StartupID != startup.ID {
		return response.Unauthorized(ctx, "Decision does not belong to your startup")
	}

	return response.Success(ctx, "Decision retrieved successfully", decision)
}
