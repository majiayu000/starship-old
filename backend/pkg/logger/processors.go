package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ConsoleProcessor 控制台日志处理器
type ConsoleProcessor struct {
	writer io.Writer
}

// NewConsoleProcessor 创建控制台日志处理器
func NewConsoleProcessor() *ConsoleProcessor {
	return &ConsoleProcessor{
		writer: os.Stdout,
	}
}

// Process 处理日志
func (p *ConsoleProcessor) Process(log interface{}) error {
	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	_, err = fmt.Fprintf(p.writer, "%s\n", data)
	return err
}

// Close 关闭处理器
func (p *ConsoleProcessor) Close() error {
	return nil
}

// FileProcessor 文件日志处理器
type FileProcessor struct {
	writer  io.Writer
	file    *os.File
	path    string
	rotator *RotatingFileWriter
}

// NewFileProcessor 创建文件日志处理器
func NewFileProcessor(path string) (*FileProcessor, error) {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	processor := &FileProcessor{
		path: path,
	}

	return processor, nil
}

// NewRotatingFileProcessor 创建带轮转功能的文件日志处理器
func NewRotatingFileProcessor(path string, maxSize string, maxAge string, maxBackups int) (*FileProcessor, error) {
	// 创建轮转写入器
	rotator, err := NewRotatingFileWriter(path, maxSize, maxAge, maxBackups)
	if err != nil {
		return nil, err
	}

	processor := &FileProcessor{
		path:    path,
		writer:  rotator,
		rotator: rotator,
	}

	return processor, nil
}

// Process 处理日志
func (p *FileProcessor) Process(log interface{}) error {
	// 如果还没有初始化写入器，则创建普通文件
	if p.writer == nil {
		file, err := os.OpenFile(p.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		p.writer = file
		p.file = file
	}

	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	_, err = fmt.Fprintf(p.writer, "%s\n", data)
	return err
}

// Close 关闭处理器
func (p *FileProcessor) Close() error {
	if p.rotator != nil {
		return p.rotator.Close()
	}
	if p.file != nil {
		return p.file.Close()
	}
	return nil
}

// RemoteProcessor 远程日志处理器
type RemoteProcessor struct {
	client *http.Client
	url    string
}

// NewRemoteProcessor 创建远程日志处理器
func NewRemoteProcessor(url string) *RemoteProcessor {
	return &RemoteProcessor{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		url: url,
	}
}

// Process 处理日志
func (p *RemoteProcessor) Process(log interface{}) error {
	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	resp, err := p.client.Post(p.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send log: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Close 关闭处理器
func (p *RemoteProcessor) Close() error {
	return nil
}
