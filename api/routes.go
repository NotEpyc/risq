package api

import (
	"time"

	"risq_backend/api/controller"
	"risq_backend/pkg/jwt"
	"risq_backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, controllers *Controllers, jwtService *jwt.Service) {
	// Add global middlewares
	app.Use(middlewares.CORS())
	app.Use(middlewares.Logger())

	// Health check with detailed status
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"message":   "Risk Assessment API is running",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
			"service":   "risq-backend",
		})
	})

	// API v1 routes
	v1 := app.Group("/api/v1")

	// Public routes (no authentication required)
	public := v1.Group("/auth")
	public.Post("/signup", controllers.UserController.Create)
	public.Post("/login", controllers.UserController.Login)

	// Protected routes (require JWT authentication)
	protected := v1.Group("/", middlewares.JWT(jwtService))

	// Startup routes (require authentication)
	startups := protected.Group("/startup")
	startups.Post("/onboard", controllers.StartupController.Submit)
	startups.Get("/profile", controllers.StartupController.GetByUser)

	// Decision routes (require startup context - user must have completed onboarding)
	decisions := protected.Group("/decisions", middlewares.StartupContext())
	decisions.Post("/speculate", controllers.DecisionController.SpeculateDecision)
	decisions.Post("/confirm", controllers.DecisionController.ConfirmDecision)
	decisions.Get("/", controllers.DecisionController.GetDecisions)
	decisions.Get("/:id", controllers.DecisionController.GetDecision)

	// Risk assessment routes (require startup context)
	risks := protected.Group("/risk", middlewares.StartupContext())
	risks.Get("/current", controllers.RiskController.GetCurrentRisk)
	risks.Get("/history", controllers.RiskController.GetRiskEvolution)
}

type Controllers struct {
	UserController     *controller.UserController
	StartupController  *controller.StartupController
	DecisionController *controller.DecisionController
	RiskController     *controller.RiskController
}
