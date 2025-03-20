package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App       AppConfig
	Server    ServerConfig
	Databases DatabasesConfig
	Auth      AuthConfig
	Logger    LoggerConfig
	Features  FeaturesConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string
	Environment string
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabasesConfig holds configuration for all supported databases
type DatabasesConfig struct {
	Default       string              // Name of the default database to use
	PostgreSQL    PostgreSQLConfig    `mapstructure:"postgresql"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
	VectorDB      VectorDBConfig      `mapstructure:"vectordb"`
	Redis         RedisConfig         `mapstructure:"redis"`
}

// PostgreSQLConfig holds PostgreSQL-specific configuration
type PostgreSQLConfig struct {
	Enabled         bool
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// ElasticsearchConfig holds Elasticsearch-specific configuration
type ElasticsearchConfig struct {
	Enabled   bool
	Addresses []string
	Username  string
	Password  string
}

// VectorDBConfig holds vector database configuration
type VectorDBConfig struct {
	Enabled  bool
	Type     string // e.g., "pinecone", "milvus", etc.
	Endpoint string
	APIKey   string
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	Enabled      bool
	Address      string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Enabled        bool
	JWT            JWTConfig
	EnabledMethods []string // "jwt", "oauth", "basic", etc.
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	SecretKey      string
	Issuer         string
	ExpiryDuration time.Duration
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level      string           `mapstructure:"level"`
	Format     string           `mapstructure:"format"`
	Processors ProcessorsConfig `mapstructure:"processors"`
	Tracing    TracingConfig    `mapstructure:"tracing"`
}

// ProcessorsConfig holds configuration for log processors
type ProcessorsConfig struct {
	Console ConsoleConfig `mapstructure:"console"`
	File    FileConfig    `mapstructure:"file"`
	ELK     ELKConfig     `mapstructure:"elk"`
}

// ConsoleConfig holds console processor configuration
type ConsoleConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// FileConfig holds file processor configuration
type FileConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Path     string `mapstructure:"path"`
	Rotation struct {
		MaxSize    string `mapstructure:"max_size"`
		MaxAge     string `mapstructure:"max_age"`
		MaxBackups int    `mapstructure:"max_backups"`
	} `mapstructure:"rotation"`
}

// ELKConfig holds ELK processor configuration
type ELKConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Endpoint    string `mapstructure:"endpoint"`
	IndexPrefix string `mapstructure:"index_prefix"`
}

// TracingConfig holds tracing configuration
type TracingConfig struct {
	Enabled bool    `mapstructure:"enabled"`
	Sampler float64 `mapstructure:"sampler"`
}

// FeaturesConfig allows enabling/disabling specific features
type FeaturesConfig struct {
	EnableAuth         bool `mapstructure:"enable_auth"`
	EnableRateLimiting bool `mapstructure:"enable_rate_limiting"`
	EnableCORS         bool `mapstructure:"enable_cors"`
	EnableCaching      bool `mapstructure:"enable_caching"`
}

// Load loads configuration from config files and environment variables
func Load() (*Config, error) {
	v := viper.New()

	// Set default configurations
	setDefaults(v)

	// Read configuration file
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	if err := v.ReadInConfig(); err != nil {
		// If config file is not found, just use environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Override with environment variables
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "cc-starship")
	v.SetDefault("app.environment", "development")

	// Server defaults
	v.SetDefault("server.address", ":8080")
	v.SetDefault("server.readTimeout", "15s")
	v.SetDefault("server.writeTimeout", "15s")
	v.SetDefault("server.idleTimeout", "60s")

	// Database defaults
	v.SetDefault("databases.default", "") // 默认不使用任何数据库
	v.SetDefault("databases.postgresql.enabled", false)
	v.SetDefault("databases.postgresql.maxOpenConns", 25)
	v.SetDefault("databases.postgresql.maxIdleConns", 25)
	v.SetDefault("databases.postgresql.connMaxLifetime", "5m")

	// Database defaults - Elasticsearch
	v.SetDefault("databases.elasticsearch.enabled", false)

	// Database defaults - VectorDB
	v.SetDefault("databases.vectordb.enabled", false)
	v.SetDefault("databases.vectordb.type", "")

	// Database defaults - Redis
	v.SetDefault("databases.redis.enabled", false)
	v.SetDefault("databases.redis.address", "localhost:6379")
	v.SetDefault("databases.redis.password", "")
	v.SetDefault("databases.redis.db", 0)
	v.SetDefault("databases.redis.poolSize", 10)
	v.SetDefault("databases.redis.minIdleConns", 5)
	v.SetDefault("databases.redis.maxRetries", 3)
	v.SetDefault("databases.redis.dialTimeout", "5s")

	// Auth defaults
	v.SetDefault("auth.enabled", true)
	v.SetDefault("auth.enabledMethods", []string{"jwt"})
	v.SetDefault("auth.jwt.issuer", "api")
	v.SetDefault("auth.jwt.expiryDuration", "24h")

	// Features defaults
	v.SetDefault("features.enable_auth", true)
	v.SetDefault("features.enable_rate_limiting", false)
	v.SetDefault("features.enable_cors", true)
	v.SetDefault("features.enable_caching", false)

	// Logger defaults
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")
	v.SetDefault("logger.processors.console.enabled", true)
	v.SetDefault("logger.processors.file.enabled", true)
	v.SetDefault("logger.processors.file.path", "logs/app.log")
	v.SetDefault("logger.processors.file.rotation.max_size", "100MB")
	v.SetDefault("logger.processors.file.rotation.max_age", "7d")
	v.SetDefault("logger.processors.file.rotation.max_backups", 10)
	v.SetDefault("logger.processors.elk.enabled", false)
	v.SetDefault("logger.processors.elk.endpoint", "http://localhost:9200")
	v.SetDefault("logger.processors.elk.index_prefix", "cc-starship")
	v.SetDefault("logger.tracing.enabled", true)
	v.SetDefault("logger.tracing.sampler", 1.0)
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// If auth is enabled, validate JWT config
	if config.Auth.Enabled {
		// Validate required JWT secret key
		if config.Auth.JWT.SecretKey == "" {
			// For development, generate a random secret if not provided
			if config.App.Environment == "development" {
				config.Auth.JWT.SecretKey = "dev-secret-key"
			} else {
				return fmt.Errorf("jwt.secretKey is required when auth is enabled")
			}
		}
	}

	// Check if default database is enabled
	defaultDB := config.Databases.Default
	if defaultDB == "" {
		return nil // Allow running without any database
	}

	switch defaultDB {
	case "postgresql":
		if !config.Databases.PostgreSQL.Enabled {
			return fmt.Errorf("default database 'postgresql' is not enabled")
		}
		if config.Databases.PostgreSQL.URL == "" {
			return fmt.Errorf("postgresql URL is required when it's set as default")
		}
	case "elasticsearch":
		if !config.Databases.Elasticsearch.Enabled {
			return fmt.Errorf("default database 'elasticsearch' is not enabled")
		}
		if len(config.Databases.Elasticsearch.Addresses) == 0 {
			return fmt.Errorf("elasticsearch addresses are required when it's set as default")
		}
	case "vectordb":
		if !config.Databases.VectorDB.Enabled {
			return fmt.Errorf("default database 'vectordb' is not enabled")
		}
		if config.Databases.VectorDB.Endpoint == "" {
			return fmt.Errorf("vectordb endpoint is required when it's set as default")
		}
	case "redis":
		if !config.Databases.Redis.Enabled {
			return fmt.Errorf("default database 'redis' is not enabled")
		}
		if config.Databases.Redis.Address == "" {
			return fmt.Errorf("redis address is required when it's set as default")
		}
	default:
		return fmt.Errorf("unsupported default database type: %s", defaultDB)
	}

	return nil
}
