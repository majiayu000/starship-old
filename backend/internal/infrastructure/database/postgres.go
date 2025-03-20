package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// PostgresDB represents a PostgreSQL database connection
type PostgresDB struct {
	*sql.DB
	logger *logger.Logger
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg *config.Config, logger *logger.Logger) (*PostgresDB, error) {
	// Open database connection
	db, err := sql.Open("postgres", cfg.Databases.PostgreSQL.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.Databases.PostgreSQL.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Databases.PostgreSQL.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Databases.PostgreSQL.ConnMaxLifetime)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Connected to PostgreSQL database")

	return &PostgresDB{
		DB:     db,
		logger: logger,
	}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	p.logger.Info("Closing PostgreSQL database connection")
	return p.DB.Close()
}

// HealthCheck 检查数据库连接的健康状况
func (p *PostgresDB) HealthCheck(ctx context.Context) error {
	return p.PingContext(ctx)
}

// Type 返回数据库类型
func (p *PostgresDB) Type() string {
	return "postgresql"
}

// Transaction executes a function within a database transaction
func (p *PostgresDB) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := p.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
