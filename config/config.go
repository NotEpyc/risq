package config

import (
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig      `json:"app"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	NATS     NATSConfig     `json:"nats"`
	LLM      LLMConfig      `json:"llm"`
	JWT      JWTConfig      `json:"jwt"`
	Log      LogConfig      `json:"log"`
	External ExternalConfig `json:"external"`
}

type AppConfig struct {
	Name string `json:"name"`
	Port string `json:"port"`
	Host string `json:"host"`
	Env  string `json:"env"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type NATSConfig struct {
	URL               string `json:"url"`
	ClusterID         string `json:"cluster_id"`
	ClientID          string `json:"client_id"`
	MaxReconnects     int    `json:"max_reconnects"`
	ReconnectWait     string `json:"reconnect_wait"`
	ConnectionTimeout string `json:"connection_timeout"`
}

type LLMConfig struct {
	OpenAIAPIKey string  `json:"openai_api_key"`
	Model        string  `json:"model"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
}

type JWTConfig struct {
	Secret    string `json:"secret"`
	ExpiresIn string `json:"expires_in"`
}

type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

type ExternalConfig struct {
	MarketDataAPIKey string `json:"market_data_api_key"`
	MarketDataURL    string `json:"market_data_url"`
	NewsAPIKey       string `json:"news_api_key"`
	NewsAPIURL       string `json:"news_api_url"`
}

func Load() *Config {
	// For Railway compatibility, check PORT first, then APP_PORT, then default
	port := getEnv("PORT", "")
	if port == "" {
		port = getEnv("APP_PORT", "8080")
	}

	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "Smart Risk Assessment API"),
			Port: port,
			Host: getEnv("APP_HOST", "0.0.0.0"), // Changed to 0.0.0.0 for Railway
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnvWithFallback([]string{"DB_HOST", "PGHOST"}, "localhost"),
			Port:     getEnvWithFallback([]string{"DB_PORT", "PGPORT"}, "5432"),
			User:     getEnvWithFallback([]string{"DB_USER", "PGUSER"}, "postgres"),
			Password: getEnvWithFallback([]string{"DB_PASSWORD", "PGPASSWORD"}, ""),
			Name:     getEnvWithFallback([]string{"DB_NAME", "PGDATABASE"}, "railway"),
			SSLMode:  getEnv("DB_SSL_MODE", "require"),
		},
		Redis: RedisConfig{
			Host:     getEnvWithFallback([]string{"REDIS_HOST", "REDIS_PRIVATE_URL"}, "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnvWithFallback([]string{"REDIS_PASSWORD", "REDIS_PASSWORD"}, ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		NATS: NATSConfig{
			URL:               getEnv("NATS_URL", "nats://localhost:4222"),
			ClusterID:         getEnv("NATS_CLUSTER_ID", "risq-cluster"),
			ClientID:          getEnv("NATS_CLIENT_ID", "risq-api"),
			MaxReconnects:     getEnvAsInt("NATS_MAX_RECONNECTS", 10),
			ReconnectWait:     getEnv("NATS_RECONNECT_WAIT", "2s"),
			ConnectionTimeout: getEnv("NATS_CONNECTION_TIMEOUT", "5s"),
		},
		LLM: LLMConfig{
			OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
			Model:        getEnv("OPENAI_MODEL", "gpt-4"),
			Temperature:  getEnvAsFloat("OPENAI_TEMPERATURE", 0.7),
			MaxTokens:    getEnvAsInt("OPENAI_MAX_TOKENS", 2048),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "default_secret_change_in_production"),
			ExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		External: ExternalConfig{
			MarketDataAPIKey: getEnv("MARKET_DATA_API_KEY", ""),
			MarketDataURL:    getEnv("MARKET_DATA_URL", "https://api.marketdata.com/v1"),
			NewsAPIKey:       getEnv("NEWS_API_KEY", ""),
			NewsAPIURL:       getEnv("NEWS_API_URL", "https://newsapi.org/v2"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvWithFallback(keys []string, defaultValue string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return defaultValue
}
