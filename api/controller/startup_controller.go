package controller

import (
	"fmt"
	"risq_backend/internal/contextmem"
	"risq_backend/internal/risk"
	"risq_backend/internal/startup"
	"risq_backend/pkg/logger"
	"risq_backend/pkg/response"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StartupController struct {
	startupService    startup.Service
	riskService       risk.Service
	contextMemService contextmem.Service
	eventManager      EventManagerInterface
}

// EventManagerInterface defines the interface for event publishing
type EventManagerInterface interface {
	PublishStartupOnboarded(startupID, userID uuid.UUID, startupData, founderCV, businessPlan map[string]interface{}) error
}

func NewStartupController(startupService startup.Service, riskService risk.Service, contextMemService contextmem.Service) *StartupController {
	return &StartupController{
		startupService:    startupService,
		riskService:       riskService,
		contextMemService: contextMemService,
		eventManager:      nil, // Will be set by the app
	}
}

// SetEventManager sets the event manager for publishing events
func (c *StartupController) SetEventManager(eventManager EventManagerInterface) {
	c.eventManager = eventManager
}

func (c *StartupController) Submit(ctx *fiber.Ctx) error {
	logger.Info("=== STEP 1: Starting Submit function")

	// Try parsing directly into the struct
	var input startup.StartupOnboardingInput
	logger.Info("=== STEP 2: Starting BodyParser")
	if err := ctx.BodyParser(&input); err != nil {
		logger.Errorf("=== ERROR in BodyParser: %v", err)

		// Check if it's a JSON unmarshal error with specific field
		if strings.Contains(err.Error(), "implementation_plan") {
			logger.Error("=== FOUND THE ISSUE: implementation_plan field parsing error!")
			return response.BadRequest(ctx, "implementation_plan field has wrong type. Expected string, got object", err)
		}

		return response.BadRequest(ctx, "Invalid input", err)
	}

	logger.Info("=== STEP 3: BodyParser succeeded")
	logger.Infof("=== SUCCESS: Parsed startup: %s", input.Name)

	// Continue with basic validation
	if input.Name == "" {
		return response.BadRequest(ctx, "Startup name is required", nil)
	}
	if input.Industry == "" {
		return response.BadRequest(ctx, "Industry is required", nil)
	}
	if input.FundingStage == "" {
		return response.BadRequest(ctx, "Funding stage is required", nil)
	}
	if input.Location == "" {
		return response.BadRequest(ctx, "Location is required", nil)
	}

	logger.Info("=== STEP 4: Basic validation passed")

	// Get user ID from context
	userIDInterface := ctx.Locals("user_id")
	if userIDInterface == nil {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	userID, parseErr := uuid.Parse(userIDInterface.(string))
	if parseErr != nil {
		return response.BadRequest(ctx, "Invalid user ID", parseErr)
	}

	logger.Infof("=== STEP 5: Got user ID: %s", userID)

	// Check if user already has a startup
	existingStartup, err := c.startupService.GetByUserID(ctx.Context(), userID)
	if err == nil && existingStartup != nil {
		logger.Errorf("=== ERROR: User already has startup: %s", existingStartup.ID)
		return response.BadRequest(ctx, "User already has a startup profile", nil)
	}

	logger.Info("=== STEP 6: Creating startup in database")
	
	// Create startup with user linking
	startup, err := c.startupService.CreateWithUserLink(ctx.Context(), input, userID)
	if err != nil {
		logger.Errorf("=== ERROR creating startup: %v", err)
		return response.InternalError(ctx, "Failed to create startup", err)
	}

	logger.Infof("=== STEP 7: Startup created successfully with ID: %s", startup.ID)

	// Create initial risk profile
	startupInfo := formatStartupInfo(startup, input)
	riskProfile, err := c.riskService.CreateInitialProfile(ctx.Context(), startup.ID, startupInfo)
	if err != nil {
		logger.Warnf("=== WARNING: Failed to create initial risk profile: %v", err)
		// Continue without failing the whole request
	} else {
		logger.Infof("=== STEP 8: Initial risk profile created with score: %.2f", riskProfile.Score)
	}

	// Store startup context in memory for future decisions
	metadata := map[string]interface{}{
		"type":      "startup_onboarding",
		"industry":  startup.Industry,
		"stage":     startup.FundingStage,
		"team_size": startup.TeamSize,
		"location":  startup.Location,
	}

	if err := c.contextMemService.StoreStartupContext(ctx.Context(), startup.ID, startupInfo, metadata); err != nil {
		logger.Warnf("=== WARNING: Failed to store startup context: %v", err)
		// Continue without failing the whole request
	} else {
		logger.Info("=== STEP 9: Startup context stored successfully")
	}

	// Publish startup onboarded event to trigger the event-driven flow
	if c.eventManager != nil {
		logger.Info("=== STEP 10: Publishing startup onboarded event")
		
		startupData := map[string]interface{}{
			"id":             startup.ID.String(),
			"name":           startup.Name,
			"description":    startup.Description,
			"industry":       startup.Industry,
			"sector":         input.Sector,
			"funding_stage":  startup.FundingStage,
			"location":       startup.Location,
			"founded_date":   startup.FoundedDate,
			"team_size":      startup.TeamSize,
			"target_market":  input.TargetMarket,
			"business_model": input.BusinessModel,
		}

		founderCV := map[string]interface{}{
			"founders": input.FounderDetails,
		}

		businessPlan := map[string]interface{}{
			"implementation_plan":    input.ImplementationPlan,
			"development_timeline":   input.DevelopmentTimeline,
			"go_to_market_strategy":  input.GoToMarketStrategy,
			"initial_investment":     input.InitialInvestment,
			"monthly_burn_rate":      input.MonthlyBurnRate,
			"projected_revenue":      input.ProjectedRevenue,
			"funding_requirement":    input.FundingRequirement,
			"revenue_streams":        input.RevenueStreams,
			"technology_stack":       input.TechnologyStack,
			"competitor_analysis":    input.CompetitorAnalysis,
		}

		if err := c.eventManager.PublishStartupOnboarded(startup.ID, userID, startupData, founderCV, businessPlan); err != nil {
			logger.Warnf("=== WARNING: Failed to publish startup onboarded event: %v", err)
			// Continue without failing the whole request
		} else {
			logger.Infof("=== STEP 11: Published startup onboarded event for startup %s", startup.ID)
		}
	} else {
		logger.Warn("=== WARNING: Event manager not set, skipping event publication")
	}

	responseData := map[string]interface{}{
		"startup":      startup,
		"risk_profile": riskProfile,
		"message":      "Startup profile created successfully and event-driven analysis initiated",
	}

	logger.Info("=== SUCCESS: Startup onboarding completed")
	return response.Success(ctx, "Startup onboarded successfully", responseData)
}

func (c *StartupController) GetByID(ctx *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotImplemented, "Not implemented yet")
}

func (c *StartupController) GetByUser(ctx *fiber.Ctx) error {
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
		return response.NotFound(ctx, "No startup found for user")
	}

	return response.Success(ctx, "Startup profile retrieved successfully", startup)
}

// Helper function to format startup information for AI analysis
func formatStartupInfo(startup *startup.Startup, input startup.StartupOnboardingInput) string {
	return fmt.Sprintf(`
STARTUP PROFILE:
Name: %s
Industry: %s
Sector: %s
Funding Stage: %s
Location: %s
Team Size: %d
Business Model: %s
Target Market: %s

IMPLEMENTATION PLAN:
%s

DEVELOPMENT TIMELINE:
%s

GO-TO-MARKET STRATEGY:
%s

FINANCIAL INFORMATION:
- Initial Investment: $%.2f
- Monthly Burn Rate: $%.2f
- Projected Revenue: $%.2f
- Funding Requirement: $%.2f

REVENUE STREAMS:
%v

TECHNOLOGY STACK:
%v

COMPETITOR ANALYSIS:
%s

FOUNDER INFORMATION:
%v
`, 
		startup.Name,
		startup.Industry,
		input.Sector,
		startup.FundingStage,
		startup.Location,
		startup.TeamSize,
		input.BusinessModel,
		input.TargetMarket,
		input.ImplementationPlan,
		input.DevelopmentTimeline,
		input.GoToMarketStrategy,
		input.InitialInvestment,
		input.MonthlyBurnRate,
		input.ProjectedRevenue,
		input.FundingRequirement,
		input.RevenueStreams,
		input.TechnologyStack,
		input.CompetitorAnalysis,
		input.FounderDetails,
	)
}
