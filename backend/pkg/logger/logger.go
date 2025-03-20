package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/majiayu000/cc-starship/pkg/config"
)

// Logger 是原始的日志器，为兼容性保留
type Logger struct {
	level      Level
	format     Format
	output     io.Writer
	logManager *LogManager
}

// Level represents log levels
type Level int

// Format represents log formats
type Format int

const (
	// Log levels
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel

	// Log formats
	TextFormat Format = iota
	JSONFormat
)

// 将旧的日志级别映射到新的日志级别
func mapLevelToLogLevel(level Level) LogLevel {
	switch level {
	case DebugLevel:
		return LogLevelDebug
	case InfoLevel:
		return LogLevelInfo
	case WarnLevel:
		return LogLevelWarn
	case ErrorLevel:
		return LogLevelError
	case FatalLevel:
		return LogLevelFatal
	default:
		return LogLevelInfo
	}
}

// New creates a new Logger instance with backward compatibility
func New(cfg config.LoggerConfig) *Logger {
	level := parseLevel(cfg.Level)
	format := parseFormat(cfg.Format)

	return &Logger{
		level:  level,
		format: format,
		output: os.Stdout,
	}
}

// SetLogManager sets the log manager for forwarding logs
func (l *Logger) SetLogManager(manager *LogManager) {
	l.logManager = manager
}

// parseLevel converts a string level to Level type
func parseLevel(level string) Level {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// parseFormat converts a string format to Format type
func parseFormat(format string) Format {
	switch format {
	case "json":
		return JSONFormat
	case "text":
		return TextFormat
	default:
		return TextFormat
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	if l.level <= DebugLevel {
		if l.logManager != nil {
			l.logSystemLog(LogLevelDebug, msg)
		} else {
			l.log("DEBUG", msg)
		}
	}
}

// Info logs an informational message
func (l *Logger) Info(msg string) {
	if l.level <= InfoLevel {
		if l.logManager != nil {
			l.logSystemLog(LogLevelInfo, msg)
		} else {
			l.log("INFO", msg)
		}
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	if l.level <= WarnLevel {
		if l.logManager != nil {
			l.logSystemLog(LogLevelWarn, msg)
		} else {
			l.log("WARN", msg)
		}
	}
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	if l.level <= ErrorLevel {
		if l.logManager != nil {
			l.logSystemLog(LogLevelError, msg)
		} else {
			l.log("ERROR", msg)
		}
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	if l.level <= FatalLevel {
		if l.logManager != nil {
			l.logSystemLog(LogLevelFatal, msg)
		} else {
			l.log("FATAL", msg)
		}
		os.Exit(1)
	}
}

// log handles the actual logging for legacy format
func (l *Logger) log(level, msg string) {
	timestamp := time.Now().Format(time.RFC3339)

	switch l.format {
	case JSONFormat:
		fmt.Fprintf(l.output, `{"timestamp":"%s","level":"%s","message":"%s"}`+"\n",
			timestamp, level, msg)
	default:
		fmt.Fprintf(l.output, "[%s] %s: %s\n",
			timestamp, level, msg)
	}
}

// logSystemLog sends a system log to the log manager
func (l *Logger) logSystemLog(level LogLevel, msg string) {
	systemLog := &SystemLog{
		BaseLog: BaseLog{
			LogID:     GenerateLogID(),
			Timestamp: time.Now(),
			Type:      LogTypeSystem,
			Level:     level,
			Service:   "cc-starship",
			Message:   msg,
		},
		Component: "system",
		Event:     "log",
	}

	l.logManager.Send(systemLog)
}

// SetOutput sets the output destination for the logger
func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

// SetLevel sets the logger level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// WithField creates a logger entry with a single field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	// This is a simplified implementation for compatibility
	return l
}

// WithFields creates a logger entry with multiple fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	// This is a simplified implementation for compatibility
	return l
}

// LogBusiness provides a convenient way to log business events
func (l *Logger) LogBusiness(level LogLevel, action, userID, resourceID string, details map[string]interface{}) {
	if l.logManager == nil {
		return
	}

	businessLog := &BusinessLog{
		BaseLog: BaseLog{
			LogID:     GenerateLogID(),
			Timestamp: time.Now(),
			Type:      LogTypeBusiness,
			Level:     level,
			Service:   "cc-starship",
		},
		Action:     action,
		UserID:     userID,
		ResourceID: resourceID,
		Details:    details,
	}

	l.logManager.Send(businessLog)
}

// LogSecurity provides a convenient way to log security events
func (l *Logger) LogSecurity(level LogLevel, action, userID string, details map[string]interface{}) {
	if l.logManager == nil {
		return
	}

	// For security logs, we use a business log with security type
	businessLog := &BusinessLog{
		BaseLog: BaseLog{
			LogID:     GenerateLogID(),
			Timestamp: time.Now(),
			Type:      LogTypeSecurity,
			Level:     level,
			Service:   "cc-starship",
		},
		Action:  action,
		UserID:  userID,
		Details: details,
	}

	l.logManager.Send(businessLog)
}
