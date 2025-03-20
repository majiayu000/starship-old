package logger

import "time"

// LogType 定义日志类型
type LogType string

const (
	LogTypeAccess   LogType = "access"   // 访问日志
	LogTypeBusiness LogType = "business" // 业务日志
	LogTypeSystem   LogType = "system"   // 系统日志
	LogTypeSecurity LogType = "security" // 安全日志
)

// LogLevel 定义日志级别
type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelFatal LogLevel = "FATAL"
)

// BaseLog 基础日志结构
type BaseLog struct {
	LogID     string    `json:"log_id"`    // 日志唯一ID
	Timestamp time.Time `json:"timestamp"` // 时间戳
	Type      LogType   `json:"type"`      // 日志类型
	Level     LogLevel  `json:"level"`     // 日志级别
	Service   string    `json:"service"`   // 服务名称
	Instance  string    `json:"instance"`  // 实例ID
	TraceID   string    `json:"trace_id"`  // 追踪ID
	SpanID    string    `json:"span_id"`   // 跨度ID
	Message   string    `json:"message"`   // 日志消息
	Context   any       `json:"context"`   // 上下文信息
}

// RequestInfo 请求信息
type RequestInfo struct {
	Method    string            `json:"method"`     // HTTP方法
	Path      string            `json:"path"`       // 请求路径
	Query     string            `json:"query"`      // 查询参数
	Headers   map[string]string `json:"headers"`    // 请求头
	ClientIP  string            `json:"client_ip"`  // 客户端IP
	UserAgent string            `json:"user_agent"` // 用户代理
	Body      string            `json:"body"`       // 请求体
}

// ResponseInfo 响应信息
type ResponseInfo struct {
	Status   int               `json:"status"`     // 状态码
	Headers  map[string]string `json:"headers"`    // 响应头
	Body     string            `json:"body"`       // 响应体
	Size     int64             `json:"size"`       // 响应大小
	TimeMS   float64           `json:"time_ms"`    // 处理时间
	DBTimeMS float64           `json:"db_time_ms"` // 数据库时间
	CacheHit bool              `json:"cache_hit"`  // 缓存命中
}

// AccessLog 访问日志结构
type AccessLog struct {
	BaseLog
	Request  RequestInfo  `json:"request"`  // 请求信息
	Response ResponseInfo `json:"response"` // 响应信息
}

// BusinessLog 业务日志结构
type BusinessLog struct {
	BaseLog
	Action     string                 `json:"action"`      // 业务动作
	UserID     string                 `json:"user_id"`     // 用户ID
	ResourceID string                 `json:"resource_id"` // 资源ID
	Details    map[string]interface{} `json:"details"`     // 详细信息
}

// SystemLog 系统日志结构
type SystemLog struct {
	BaseLog
	Component string                 `json:"component"` // 组件名称
	Event     string                 `json:"event"`     // 事件类型
	Metrics   map[string]float64     `json:"metrics"`   // 指标数据
	Resources map[string]interface{} `json:"resources"` // 资源信息
}

// LogProcessor 日志处理器接口
type LogProcessor interface {
	Process(log interface{}) error
	Close() error
}
