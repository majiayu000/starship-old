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

// ElasticsearchDB 代表 Elasticsearch 数据库连接
type ElasticsearchDB struct {
	// 暂时使用简单HTTP客户端，实际项目中可以使用Elasticsearch官方客户端
	client      *http.Client
	addresses   []string
	username    string
	password    string
	logger      *logger.Logger
	isConnected bool
}

// NewElasticsearchDB 创建新的Elasticsearch数据库连接
func NewElasticsearchDB(cfg *config.Config, logger *logger.Logger) (*ElasticsearchDB, error) {
	// 检查配置
	if !cfg.Databases.Elasticsearch.Enabled {
		return nil, errors.New("elasticsearch is not enabled in configuration")
	}

	if len(cfg.Databases.Elasticsearch.Addresses) == 0 {
		return nil, errors.New("no elasticsearch addresses provided")
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	es := &ElasticsearchDB{
		client:    client,
		addresses: cfg.Databases.Elasticsearch.Addresses,
		username:  cfg.Databases.Elasticsearch.Username,
		password:  cfg.Databases.Elasticsearch.Password,
		logger:    logger,
	}

	// 检查连接
	if err := es.HealthCheck(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}

	es.isConnected = true
	logger.Info("Connected to Elasticsearch")

	return es, nil
}

// Close 关闭Elasticsearch连接
func (e *ElasticsearchDB) Close() error {
	e.logger.Info("Closing Elasticsearch connection")
	e.isConnected = false
	return nil
}

// HealthCheck 检查Elasticsearch连接的健康状况
func (e *ElasticsearchDB) HealthCheck(ctx context.Context) error {
	// 尝试连接第一个地址
	if len(e.addresses) == 0 {
		return errors.New("no elasticsearch addresses available")
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", e.addresses[0]+"/_cluster/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	// 添加基本认证
	if e.username != "" && e.password != "" {
		req.SetBasicAuth(e.username, e.password)
	}

	// 发送请求
	resp, err := e.client.Do(req)
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
func (e *ElasticsearchDB) Type() string {
	return "elasticsearch"
}
