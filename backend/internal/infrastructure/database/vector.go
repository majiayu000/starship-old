package database

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// VectorDB 代表向量数据库连接
type VectorDB struct {
	// 暂时使用简单HTTP客户端，实际项目中可以使用向量数据库的官方客户端
	client      *http.Client
	address     string
	apiKey      string
	logger      *logger.Logger
	isConnected bool
}

// NewVectorDB 创建新的向量数据库连接
func NewVectorDB(cfg *config.Config, logger *logger.Logger) (*VectorDB, error) {
	// 检查配置
	if !cfg.Databases.VectorDB.Enabled {
		return nil, errors.New("vector database is not enabled in configuration")
	}

	if cfg.Databases.VectorDB.Endpoint == "" {
		return nil, errors.New("vector database address not provided")
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	vdb := &VectorDB{
		client:  client,
		address: cfg.Databases.VectorDB.Endpoint,
		apiKey:  cfg.Databases.VectorDB.APIKey,
		logger:  logger,
	}

	// 检查连接
	if err := vdb.HealthCheck(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to vector database: %w", err)
	}

	vdb.isConnected = true
	logger.Info("Connected to Vector Database")

	return vdb, nil
}

// Close 关闭向量数据库连接
func (v *VectorDB) Close() error {
	v.logger.Info("Closing Vector Database connection")
	v.isConnected = false
	return nil
}

// HealthCheck 检查向量数据库连接的健康状况
func (v *VectorDB) HealthCheck(ctx context.Context) error {
	// 创建健康检查请求
	req, err := http.NewRequestWithContext(ctx, "GET", v.address+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	// 添加API密钥
	if v.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+v.apiKey)
	}

	// 发送请求
	resp, err := v.client.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Type 返回数据库类型
func (v *VectorDB) Type() string {
	return "vector"
}
