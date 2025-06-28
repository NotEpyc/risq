package app

import (
	"context"
	"fmt"
	"os"
	"time"

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
	"risq_backend/pkg/middlewares"
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

	// Initialize database (optional for health checks)
	if err := a.initDatabase(); err != nil {
		logger.Warnf("Database initialization failed: %v", err)
		logger.Info("Continuing without database (health checks will still work)")
	}

	// Initialize Redis (optional for health checks)
	if err := a.initRedis(); err != nil {
		logger.Warnf("Redis initialization failed: %v", err)
		logger.Info("Continuing without Redis (health checks will still work)")
	}

	// Initialize LLM client
	a.initLLMClient()

	// Initialize JWT service
	a.initJWTService()

	// Initialize event-driven services (optional for health checks)
	if err := a.initEventManager(); err != nil {
		logger.Warnf("Event manager initialization failed: %v", err)
		logger.Info("Continuing without event manager (health checks will still work)")
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
	logger.Info("Initializing database connection...")

	// Try DATABASE_URL first (Railway's preferred method)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		logger.Info("DATABASE_URL found - attempting connection...")
		// Log partial URL for debugging (hide sensitive parts)
		urlLen := len(databaseURL)
		if urlLen > 20 {
			logger.Infof("DATABASE_URL format: %s...%s", databaseURL[:8], databaseURL[urlLen-10:])
		}

		db, err := database.NewFromURL(databaseURL)
		if err != nil {
			logger.Errorf("DATABASE_URL connection failed: %v", err)
			// Don't return error - fall back to individual variables
			logger.Warn("Falling back to individual database environment variables...")
		} else {
			a.db = db
			logger.Info("✅ Database connected successfully via DATABASE_URL")
			return nil
		}
	} else {
		logger.Warn("DATABASE_URL not found, using individual environment variables")
	}

	// Fallback to individual environment variables
	logger.Info("Attempting connection with individual database configuration...")
	logger.Infof("DB Config: host=%s, port=%s, user=%s, dbname=%s, sslmode=%s",
		a.config.Database.Host, a.config.Database.Port, a.config.Database.User,
		a.config.Database.Name, a.config.Database.SSLMode)

	// Check if required fields are available
	if a.config.Database.Host == "" || a.config.Database.Host == "localhost" {
		logger.Error("Database host not configured properly")
		return fmt.Errorf("database host not configured - check DB_HOST or DATABASE_URL")
	}

	db, err := database.New(
		a.config.Database.Host,
		a.config.Database.Port,
		a.config.Database.User,
		a.config.Database.Password,
		a.config.Database.Name,
		a.config.Database.SSLMode,
	)
	if err != nil {
		logger.Errorf("Individual DB config connection failed: %v", err)
		return fmt.Errorf("database connection failed with both methods: %w", err)
	}

	a.db = db
	logger.Info("✅ Database connected successfully via individual config")
	return nil
}

func (a *App) initRedis() error {
	logger.Info("Connecting to Redis...")

	var redisURL string

	// Try REDIS_URL first (Railway's preferred method)
	if url := os.Getenv("REDIS_URL"); url != "" {
		redisURL = url
		logger.Info("Using REDIS_URL for connection")
	} else if privateURL := os.Getenv("REDIS_PRIVATE_URL"); privateURL != "" {
		redisURL = privateURL
		logger.Info("Using REDIS_PRIVATE_URL for connection")
	} else {
		// Fallback to individual environment variables
		redisURL = fmt.Sprintf("redis://%s:%s", a.config.Redis.Host, a.config.Redis.Port)
		if a.config.Redis.Password != "" {
			redisURL = fmt.Sprintf("redis://:%s@%s:%s", a.config.Redis.Password, a.config.Redis.Host, a.config.Redis.Port)
		}
		logger.Info("Using individual Redis config for connection")
	}

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

	// Always set up health route
	a.setupHealthRoute()

	// Check if database is available for full setup
	if a.db == nil {
		logger.Warn("Database not available - setting up limited routes")
		a.setupLimitedRoutes()
		return
	}

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

	// Setup full routes (this will override the fallback auth routes)
	api.SetupRoutes(a.fiberApp, controllers, a.jwtService)

	logger.Info("Full routes set up successfully")
}

func (a *App) setupHealthRoute() {
	logger.Info("Setting up health route...")

	// Add basic middleware
	a.fiberApp.Use(middlewares.CORS())
	a.fiberApp.Use(middlewares.Logger())

	// Health check with detailed status
	a.fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"message":   "Risk Assessment API is running",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
			"service":   "risq-backend",
		})
	})

	logger.Info("Health route setup complete")
}

func (a *App) initEventManager() error {
	logger.Info("Initializing Event Manager...")

	// Check if dependencies are available
	if a.db == nil {
		return fmt.Errorf("database not available")
	}
	if a.redis == nil {
		return fmt.Errorf("redis not available")
	}

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

	// Start event manager (if available)
	if a.eventManager != nil {
		if err := a.eventManager.Start(context.Background()); err != nil {
			logger.Warnf("Failed to start event manager: %v", err)
			logger.Info("Continuing without event manager")
		} else {
			logger.Info("Event manager started successfully")
		}
	} else {
		logger.Info("Event manager not available, skipping")
	}

	address := a.config.App.Host + ":" + a.config.App.Port
	logger.Infof("Server listening on: %s", address)
	logger.Info("Health check available at: /health")
	logger.Info("Application startup complete")
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

func (a *App) setupAuthRoutes() {
	logger.Info("Setting up auth routes...")

	// Add API v1 group
	v1 := a.fiberApp.Group("/api/v1")

	// Public auth routes (these work even without database for better error messages)
	public := v1.Group("/auth")

	// Simplified auth handlers that provide meaningful errors when DB is unavailable
	public.Post("/signup", func(c *fiber.Ctx) error {
		if a.db == nil {
			return c.Status(503).JSON(fiber.Map{
				"success": false,
				"message": "Service temporarily unavailable",
				"error":   "Database connection not available. Please try again later.",
			})
		}
		// This should not happen if we reach here, but just in case
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Auth service not properly initialized",
			"error":   "Please contact support",
		})
	})

	public.Post("/login", func(c *fiber.Ctx) error {
		if a.db == nil {
			return c.Status(503).JSON(fiber.Map{
				"success": false,
				"message": "Service temporarily unavailable",
				"error":   "Database connection not available. Please try again later.",
			})
		}
		// This should not happen if we reach here, but just in case
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Auth service not properly initialized",
			"error":   "Please contact support",
		})
	})

	logger.Info("Auth routes setup complete (fallback mode)")
}

func (a *App) setupLimitedRoutes() {
	logger.Info("Setting up limited routes (no database)...")

	// Add basic middleware
	a.fiberApp.Use(middlewares.CORS())
	a.fiberApp.Use(middlewares.Logger())

	// Simple auth endpoints that return proper error messages
	v1 := a.fiberApp.Group("/api/v1")
	auth := v1.Group("/auth")

	auth.Post("/signup", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"message": "Service temporarily unavailable",
			"error":   "Database connection required for user registration",
		})
	})

	auth.Post("/login", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"message": "Service temporarily unavailable",
			"error":   "Database connection required for user authentication",
		})
	})

	logger.Info("Limited routes setup complete")
}
