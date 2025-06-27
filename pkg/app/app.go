package app

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"risq_backend/api"
	"risq_backend/api/controller"
	"risq_backend/config"
	"risq_backend/internal/contextmem"
	"risq_backend/internal/decision"
	"risq_backend/internal/llm"
	"risq_backend/internal/risk"
	"risq_backend/internal/startup"
	"risq_backend/internal/user"
	"risq_backend/pkg/cache"
	"risq_backend/pkg/database"
	"risq_backend/pkg/external"
	"risq_backend/pkg/jwt"
	pkgllm "risq_backend/pkg/llm"
	"risq_backend/pkg/logger"
)

type App struct {
	config       *config.Config
	db           *database.DB
	redis        *cache.Redis
	fiberApp     *fiber.App
	llmClient    *pkgllm.Client
	jwtService   *jwt.Service
	eventManager *EventManager
}

func New(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Initialize() error {
	// Initialize logger
	logger.Init(a.config.Log.Level)
	logger.Info("Initializing application...")

	// Initialize database
	if err := a.initDatabase(); err != nil {
		return err
	}

	// Initialize Redis
	if err := a.initRedis(); err != nil {
		return err
	}

	// Initialize LLM client
	a.initLLMClient()

	// Initialize JWT service
	a.initJWTService()

	// Initialize event-driven services
	if err := a.initEventManager(); err != nil {
		return err
	}

	// Initialize Fiber app
	a.initFiberApp()

	// Setup routes
	a.setupRoutes()

	logger.Info("Application initialized successfully")
	return nil
}

func (a *App) initJWTService() {
	logger.Info("Initializing JWT service...")
	a.jwtService = jwt.NewService(a.config.JWT.Secret, "risq-api")
	logger.Info("JWT service initialized")
}

func (a *App) initDatabase() error {
	logger.Info("Connecting to database...")
	logger.Infof("Database configuration: host=%s, port=%s, user=%s, name=%s, sslmode=%s",
		a.config.Database.Host, a.config.Database.Port, a.config.Database.User,
		a.config.Database.Name, a.config.Database.SSLMode)

	db, err := database.New(
		a.config.Database.Host,
		a.config.Database.Port,
		a.config.Database.User,
		a.config.Database.Password,
		a.config.Database.Name,
		a.config.Database.SSLMode,
	)
	if err != nil {
		return err
	}

	a.db = db
	logger.Info("Database connected and migrated successfully")
	return nil
}

func (a *App) initRedis() error {
	logger.Info("Connecting to Redis...")

	redisURL := fmt.Sprintf("redis://%s:%s", a.config.Redis.Host, a.config.Redis.Port)
	redis, err := cache.NewRedisConnection(redisURL)
	if err != nil {
		return err
	}

	a.redis = redis
	logger.Info("Redis connected successfully")
	return nil
}

func (a *App) initLLMClient() {
	logger.Info("Initializing LLM client...")
	a.llmClient = pkgllm.NewClient(a.config.LLM.OpenAIAPIKey)
	logger.Info("LLM client initialized")
}

func (a *App) initFiberApp() {
	logger.Info("Initializing Fiber app...")

	a.fiberApp = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Errorf("Fiber error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		},
	})

	logger.Info("Fiber app initialized")
}

func (a *App) setupRoutes() {
	logger.Info("Setting up routes...")

	// Create repositories
	userRepo := user.NewRepository(a.db.GetConn())
	startupRepo := startup.NewRepository(a.db.GetConn())
	decisionRepo := decision.NewRepository(a.db.GetConn())
	riskRepo := risk.NewRepository(a.db.GetConn())

	// Create services
	llmService := llm.NewService(a.llmClient)
	contextMemService := contextmem.NewService(a.redis, llmService)

	userService := user.NewService(userRepo)
	startupService := startup.NewService(startupRepo, userRepo)
	riskService := risk.NewService(riskRepo, llmService)
	decisionService := decision.NewService(decisionRepo, llmService, contextMemService, riskService)

	// Create controllers
	startupController := controller.NewStartupController(startupService, riskService, contextMemService)
	startupController.SetEventManager(a.eventManager) // Set event manager for publishing events

	controllers := &api.Controllers{
		UserController:     controller.NewUserController(userService, a.jwtService),
		StartupController:  startupController,
		DecisionController: controller.NewDecisionController(decisionService, startupService),
		RiskController:     controller.NewRiskController(riskService, startupService),
	}

	// Setup routes
	api.SetupRoutes(a.fiberApp, controllers, a.jwtService)

	logger.Info("Routes set up successfully")
}

func (a *App) initEventManager() error {
	logger.Info("Initializing Event Manager...")

	// Create external services
	marketDataService := external.NewMarketDataService(
		a.config.External.MarketDataAPIKey,
		a.config.External.MarketDataURL,
		a.config.External.NewsAPIKey,
		a.config.External.NewsAPIURL,
	)

	// Create services that will be used by event handlers
	riskRepo := risk.NewRepository(a.db.GetConn())

	llmService := llm.NewService(a.llmClient)
	contextMemService := contextmem.NewService(a.redis, llmService)
	riskService := risk.NewService(riskRepo, llmService)

	// Create event manager
	a.eventManager = NewEventManager(
		a.config,
		riskService,
		llmService,
		contextMemService,
		marketDataService,
	)

	logger.Info("Event Manager initialized")
	return nil
}

func (a *App) Start() error {
	logger.Infof("Starting server on %s:%s", a.config.App.Host, a.config.App.Port)

	// Start event manager
	if err := a.eventManager.Start(context.Background()); err != nil {
		return fmt.Errorf("failed to start event manager: %w", err)
	}

	address := a.config.App.Host + ":" + a.config.App.Port
	return a.fiberApp.Listen(address)
}

func (a *App) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down application...")

	// Stop event manager
	if a.eventManager != nil {
		if err := a.eventManager.Stop(); err != nil {
			logger.Errorf("Error stopping event manager: %v", err)
		}
	}

	// Close database connection
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			logger.Errorf("Error closing database: %v", err)
		}
	}

	// Close Redis connection
	if a.redis != nil {
		if err := a.redis.Close(); err != nil {
			logger.Errorf("Error closing Redis: %v", err)
		}
	}

	// Shutdown Fiber app
	if a.fiberApp != nil {
		if err := a.fiberApp.Shutdown(); err != nil {
			logger.Errorf("Error shutting down Fiber app: %v", err)
		}
	}

	logger.Info("Application shut down successfully")
	return nil
}
