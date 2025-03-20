package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/majiayu000/cc-starship/pkg/config"
)

// LogManager 日志管理器
type LogManager struct {
	config     *config.Config
	processors []LogProcessor
	logChan    chan interface{}
	wg         sync.WaitGroup
	done       chan struct{}
}

// NewLogManager 创建日志管理器
func NewLogManager(cfg *config.Config) (*LogManager, error) {
	manager := &LogManager{
		config:     cfg,
		processors: make([]LogProcessor, 0),
		logChan:    make(chan interface{}, 1000),
		done:       make(chan struct{}),
	}

	// 添加控制台处理器
	if cfg.Logger.Processors.Console.Enabled {
		manager.processors = append(manager.processors, NewConsoleProcessor())
	}

	// 添加文件处理器
	if cfg.Logger.Processors.File.Enabled {
		// 使用带轮转功能的文件处理器
		fileProcessor, err := NewRotatingFileProcessor(
			cfg.Logger.Processors.File.Path,
			cfg.Logger.Processors.File.Rotation.MaxSize,
			cfg.Logger.Processors.File.Rotation.MaxAge,
			cfg.Logger.Processors.File.Rotation.MaxBackups,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create file processor: %w", err)
		}
		manager.processors = append(manager.processors, fileProcessor)
	}

	// 添加远程处理器
	if cfg.Logger.Processors.ELK.Enabled {
		manager.processors = append(manager.processors, NewRemoteProcessor(cfg.Logger.Processors.ELK.Endpoint))
	}

	return manager, nil
}

// Start 启动日志管理器
func (m *LogManager) Start() {
	m.wg.Add(1)
	go m.processLogs()
}

// Stop 停止日志管理器
func (m *LogManager) Stop() {
	close(m.done)
	m.wg.Wait()

	// 关闭所有处理器
	for _, processor := range m.processors {
		if err := processor.Close(); err != nil {
			fmt.Printf("Error closing processor: %v\n", err)
		}
	}
}

// Send 发送日志
func (m *LogManager) Send(log interface{}) {
	select {
	case m.logChan <- log:
	default:
		// 如果通道已满，丢弃日志
		fmt.Printf("Log channel is full, dropping log\n")
	}
}

// processLogs 处理日志
func (m *LogManager) processLogs() {
	defer m.wg.Done()

	for {
		select {
		case log := <-m.logChan:
			for _, processor := range m.processors {
				if err := processor.Process(log); err != nil {
					fmt.Printf("Error processing log: %v\n", err)
				}
			}
		case <-m.done:
			return
		}
	}
}

// GenerateLogID 生成日志ID
func GenerateLogID() string {
	return uuid.New().String()
}

// GenerateSpanID 生成跨度ID
func GenerateSpanID() string {
	return uuid.New().String()[:8]
}

// GetHeaders 获取请求头
func GetHeaders(headers map[string][]string) map[string]string {
	result := make(map[string]string)
	for k, v := range headers {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}

// GetBody 获取请求体
func GetBody(r *http.Request) string {
	if r.Body == nil {
		return ""
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return string(body)
}
