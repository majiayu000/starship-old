package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RotatingFileWriter 是一个支持日志轮转的文件写入器
type RotatingFileWriter struct {
	mu         sync.Mutex
	path       string
	file       *os.File
	maxSize    int64         // 最大尺寸（字节）
	maxAge     time.Duration // 最大保留时间
	maxBackups int           // 最大备份数量
	size       int64         // 当前尺寸
	lastRotate time.Time     // 上次轮转时间
}

// NewRotatingFileWriter 创建一个支持日志轮转的文件写入器
func NewRotatingFileWriter(path string, maxSize string, maxAge string, maxBackups int) (*RotatingFileWriter, error) {
	// 解析最大尺寸
	size, err := parseSize(maxSize)
	if err != nil {
		return nil, fmt.Errorf("invalid max size: %w", err)
	}

	// 解析最大保留时间
	duration, err := parseDuration(maxAge)
	if err != nil {
		return nil, fmt.Errorf("invalid max age: %w", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// 打开文件
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// 获取当前文件尺寸
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &RotatingFileWriter{
		path:       path,
		file:       file,
		maxSize:    size,
		maxAge:     duration,
		maxBackups: maxBackups,
		size:       info.Size(),
		lastRotate: time.Now(),
	}, nil
}

// Write 实现 io.Writer 接口
func (w *RotatingFileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 检查是否需要轮转
	if w.size+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	// 写入日志
	n, err = w.file.Write(p)
	w.size += int64(n)
	return n, err
}

// Close 关闭文件
func (w *RotatingFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// rotate 轮转日志文件
func (w *RotatingFileWriter) rotate() error {
	// 关闭当前文件
	if err := w.file.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %w", err)
	}

	// 获取当前时间
	now := time.Now()
	timestamp := now.Format("20060102-150405.000")

	// 新文件名
	dir := filepath.Dir(w.path)
	base := filepath.Base(w.path)
	ext := filepath.Ext(base)
	prefix := strings.TrimSuffix(base, ext)
	backupPath := filepath.Join(dir, fmt.Sprintf("%s.%s%s", prefix, timestamp, ext))

	// 重命名当前文件
	if err := os.Rename(w.path, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	// 打开新文件
	file, err := os.OpenFile(w.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	// 更新状态
	w.file = file
	w.size = 0
	w.lastRotate = now

	// 清理旧文件
	go w.cleanup()

	return nil
}

// cleanup 清理旧文件
func (w *RotatingFileWriter) cleanup() {
	dir := filepath.Dir(w.path)
	base := filepath.Base(w.path)
	ext := filepath.Ext(base)
	prefix := strings.TrimSuffix(base, ext)
	pattern := fmt.Sprintf("%s.*%s", prefix, ext)

	// 获取所有备份文件
	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		fmt.Printf("Failed to list backup files: %v\n", err)
		return
	}

	// 按修改时间排序
	type fileInfo struct {
		path    string
		modTime time.Time
	}
	backups := make([]fileInfo, 0, len(files))
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		backups = append(backups, fileInfo{file, info.ModTime()})
	}

	// 按时间排序
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].modTime.Before(backups[j].modTime)
	})

	// 删除超过最大数量的文件
	if len(backups) > w.maxBackups {
		for i := 0; i < len(backups)-w.maxBackups; i++ {
			os.Remove(backups[i].path)
		}
	}

	// 删除超过最大保留时间的文件
	now := time.Now()
	for _, backup := range backups {
		if now.Sub(backup.modTime) > w.maxAge {
			os.Remove(backup.path)
		}
	}
}

// parseSize 解析大小字符串（如 "100MB"）
func parseSize(size string) (int64, error) {
	size = strings.TrimSpace(size)
	if size == "" {
		return 100 * 1024 * 1024, nil // 默认 100MB
	}

	var multiplier int64 = 1
	if strings.HasSuffix(size, "KB") {
		multiplier = 1024
		size = strings.TrimSuffix(size, "KB")
	} else if strings.HasSuffix(size, "MB") {
		multiplier = 1024 * 1024
		size = strings.TrimSuffix(size, "MB")
	} else if strings.HasSuffix(size, "GB") {
		multiplier = 1024 * 1024 * 1024
		size = strings.TrimSuffix(size, "GB")
	}

	value, err := strconv.ParseInt(strings.TrimSpace(size), 10, 64)
	if err != nil {
		return 0, err
	}

	return value * multiplier, nil
}

// parseDuration 解析持续时间字符串（如 "7d"）
func parseDuration(duration string) (time.Duration, error) {
	duration = strings.TrimSpace(duration)
	if duration == "" {
		return 7 * 24 * time.Hour, nil // 默认 7 天
	}

	if strings.HasSuffix(duration, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(duration, "d"))
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	return time.ParseDuration(duration)
}
