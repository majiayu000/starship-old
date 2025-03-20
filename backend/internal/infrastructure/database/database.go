package database

import (
	"context"
	"fmt"

	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// DataAccessManager manages database connections
type DataAccessManager struct {
	PostgreSQL    *PostgresDB
	Elasticsearch *ElasticsearchDB
	VectorDB      *VectorDB
	Redis         *RedisDB
	Logger        *logger.Logger
}

// DatabaseProvider is the common interface for all database types
type DatabaseProvider interface {
	// 关闭数据库连接
	Close() error
	// 检查数据库健康状态
	HealthCheck(ctx context.Context) error
	// 数据库类型
	Type() string
}

// NewDataAccessManager creates a new database access manager
func NewDataAccessManager(cfg *config.Config, logger *logger.Logger) (*DataAccessManager, error) {
	manager := &DataAccessManager{
		Logger: logger,
	}

	// Initialize PostgreSQL if enabled
	if cfg.Databases.PostgreSQL.Enabled {
		postgreSQL, err := NewPostgresDB(cfg, logger)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to initialize PostgreSQL: %v", err))
		} else {
			manager.PostgreSQL = postgreSQL
		}
	}

	// Initialize Elasticsearch if enabled
	if cfg.Databases.Elasticsearch.Enabled {
		es, err := NewElasticsearchDB(cfg, logger)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to initialize Elasticsearch: %v", err))
		} else {
			manager.Elasticsearch = es
		}
	}

	// Initialize VectorDB if enabled
	if cfg.Databases.VectorDB.Enabled {
		vectorDB, err := NewVectorDB(cfg, logger)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to initialize VectorDB: %v", err))
		} else {
			manager.VectorDB = vectorDB
		}
	}

	// Initialize Redis if enabled
	if cfg.Databases.Redis.Enabled {
		redisDB, err := NewRedisDB(cfg, logger)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to initialize Redis: %v", err))
		} else {
			manager.Redis = redisDB
		}
	}

	// Log initialization status
	manager.logInitializationStatus()

	return manager, nil
}

// logInitializationStatus logs which databases were successfully initialized
func (m *DataAccessManager) logInitializationStatus() {
	var initializedDBs []string

	if m.PostgreSQL != nil {
		initializedDBs = append(initializedDBs, "PostgreSQL")
	}
	if m.Elasticsearch != nil {
		initializedDBs = append(initializedDBs, "Elasticsearch")
	}
	if m.VectorDB != nil {
		initializedDBs = append(initializedDBs, "VectorDB")
	}
	if m.Redis != nil {
		initializedDBs = append(initializedDBs, "Redis")
	}

	if len(initializedDBs) == 0 {
		m.Logger.Warn("No databases were initialized")
	} else {
		m.Logger.Info(fmt.Sprintf("Successfully initialized databases: %v", initializedDBs))
	}
}

// HasAnyDatabase checks if any database is initialized
func (m *DataAccessManager) HasAnyDatabase() bool {
	return m.PostgreSQL != nil || m.Elasticsearch != nil || m.VectorDB != nil || m.Redis != nil
}

// GetDatabase returns a specific database by type
func (m *DataAccessManager) GetDatabase(dbType string) DatabaseProvider {
	switch dbType {
	case "postgresql":
		return m.PostgreSQL
	case "elasticsearch":
		return m.Elasticsearch
	case "vectordb":
		return m.VectorDB
	case "redis":
		return m.Redis
	default:
		return nil
	}
}

// Close closes all database connections
func (m *DataAccessManager) Close() {
	if m.PostgreSQL != nil {
		if err := m.PostgreSQL.Close(); err != nil {
			m.Logger.Error(fmt.Sprintf("Failed to close PostgreSQL: %v", err))
		}
	}
	if m.Elasticsearch != nil {
		if err := m.Elasticsearch.Close(); err != nil {
			m.Logger.Error(fmt.Sprintf("Failed to close Elasticsearch: %v", err))
		}
	}
	if m.VectorDB != nil {
		if err := m.VectorDB.Close(); err != nil {
			m.Logger.Error(fmt.Sprintf("Failed to close VectorDB: %v", err))
		}
	}
	if m.Redis != nil {
		if err := m.Redis.Close(); err != nil {
			m.Logger.Error(fmt.Sprintf("Failed to close Redis: %v", err))
		}
	}
}
