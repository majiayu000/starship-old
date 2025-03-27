# cc-starship 应用日志架构设计

本文档详细说明了cc-starship项目的日志架构设计，包括日志模型、处理器、格式规范和与Filebeat的集成方案。

## 目录

1. [整体架构](#整体架构)
2. [日志类型定义](#日志类型定义)
3. [日志模型设计](#日志模型设计)
4. [日志处理器](#日志处理器)
5. [配置结构](#配置结构)
6. [应用集成方式](#应用集成方式)
7. [与Filebeat集成](#与Filebeat集成)
8. [日志查询最佳实践](#日志查询最佳实践)

## 整体架构

cc-starship项目的日志系统采用多层架构设计，包括：

1. **日志产生层**：由应用代码通过Logger组件产生结构化日志
2. **日志处理层**：LogManager和Processors负责对日志进行处理和分发
3. **日志存储层**：包括控制台输出、文件存储和Elasticsearch存储
4. **日志收集层**：使用Filebeat收集文件日志并发送到Elasticsearch
5. **日志查询层**：通过Elasticsearch API或Kibana界面查询和分析日志

整体数据流向：
```
应用代码 -> Logger -> LogManager -> Processors -> 日志文件 -> Filebeat -> Elasticsearch -> Kibana
```

## 日志类型定义

我们定义了四种基本日志类型：

1. **访问日志(Access Log)**：记录API请求和响应信息，用于监控和审计
2. **业务日志(Business Log)**：记录业务操作和状态变更
3. **系统日志(System Log)**：记录系统组件状态和资源使用情况
4. **安全日志(Security Log)**：记录与安全相关的事件和操作

日志级别划分：
- **DEBUG**：开发调试信息
- **INFO**：常规操作信息
- **WARN**：潜在问题警告
- **ERROR**：错误信息
- **FATAL**：致命错误，可能导致应用崩溃

## 日志模型设计

所有日志都基于结构化JSON格式，并继承自BaseLog基类：

### 基础日志结构(BaseLog)

```go
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
```

### 访问日志(AccessLog)

```go
type AccessLog struct {
    BaseLog
    Request  RequestInfo  `json:"request"`  // 请求信息
    Response ResponseInfo `json:"response"` // 响应信息
}

type RequestInfo struct {
    Method    string            `json:"method"`     // HTTP方法
    Path      string            `json:"path"`       // 请求路径
    Query     string            `json:"query"`      // 查询参数
    Headers   map[string]string `json:"headers"`    // 请求头
    ClientIP  string            `json:"client_ip"`  // 客户端IP
    UserAgent string            `json:"user_agent"` // 用户代理
    Body      string            `json:"body"`       // 请求体
}

type ResponseInfo struct {
    Status   int               `json:"status"`     // 状态码
    Headers  map[string]string `json:"headers"`    // 响应头
    Body     string            `json:"body"`       // 响应体
    Size     int64             `json:"size"`       // 响应大小
    TimeMS   float64           `json:"time_ms"`    // 处理时间
    DBTimeMS float64           `json:"db_time_ms"` // 数据库时间
    CacheHit bool              `json:"cache_hit"`  // 缓存命中
}
```

### 业务日志(BusinessLog)

```go
type BusinessLog struct {
    BaseLog
    Action     string                 `json:"action"`      // 业务动作
    UserID     string                 `json:"user_id"`     // 用户ID
    ResourceID string                 `json:"resource_id"` // 资源ID
    Details    map[string]interface{} `json:"details"`     // 详细信息
}
```

### 系统日志(SystemLog)

```go
type SystemLog struct {
    BaseLog
    Component string                 `json:"component"` // 组件名称
    Event     string                 `json:"event"`     // 事件类型
    Metrics   map[string]float64     `json:"metrics"`   // 指标数据
    Resources map[string]interface{} `json:"resources"` // 资源信息
}
```

## 日志处理器

我们实现了三种日志处理器：

### 控制台处理器(ConsoleProcessor)

将日志输出到标准输出或标准错误，主要用于开发环境和实时监控。

```go
type ConsoleProcessor struct {
    writer io.Writer
}
```

### 文件处理器(FileProcessor)

将日志写入文件，支持日志轮转，可以根据大小或时间进行日志文件切割。

```go
type FileProcessor struct {
    writer  io.Writer
    file    *os.File
    path    string
    rotator *RotatingFileWriter
}

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
```

### 远程处理器(RemoteProcessor)

将日志直接发送到远程Elasticsearch服务器。

```go
type RemoteProcessor struct {
    client *http.Client
    url    string
}
```

## 配置结构

日志配置通过应用主配置文件`app.yaml`进行定义：

```yaml
logger:
  level: "info"
  format: "json"
  processors:
    console:
      enabled: true
    file:
      enabled: true
      path: "logs/app.log"
      rotation:
        max_size: "100MB"
        max_age: "7d"
        max_backups: 10
    elk:
      enabled: false
      endpoint: "http://localhost:9200"
      index_prefix: "cc-starship"
  tracing:
    enabled: true
    sampler: 1.0
```

对应的Go配置模型：

```go
type LoggerConfig struct {
    Level      string           `mapstructure:"level"`
    Format     string           `mapstructure:"format"`
    Processors ProcessorsConfig `mapstructure:"processors"`
    Tracing    TracingConfig    `mapstructure:"tracing"`
}

type ProcessorsConfig struct {
    Console ConsoleConfig `mapstructure:"console"`
    File    FileConfig    `mapstructure:"file"`
    ELK     ELKConfig     `mapstructure:"elk"`
}
```

## 应用集成方式

### 中间件集成

访问日志通过Gin中间件自动记录：

```go
// Logging creates a middleware that logs request information
func Logging(logManager *logger.LogManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // 创建访问日志
        accessLog := &logger.AccessLog{
            BaseLog: logger.BaseLog{
                LogID:     logger.GenerateLogID(),
                Timestamp: time.Now(),
                Type:      logger.LogTypeAccess,
                Level:     logger.LogLevelInfo,
                Service:   "cc-starship",
                TraceID:   c.GetHeader("X-Trace-ID"),
                SpanID:    logger.GenerateSpanID(),
            },
            Request: logger.RequestInfo{
                Method:    c.Request.Method,
                Path:      path,
                Query:     raw,
                Headers:   logger.GetHeaders(c.Request.Header),
                ClientIP:  c.ClientIP(),
                UserAgent: c.Request.UserAgent(),
                Body:      logger.GetBody(c.Request),
            },
        }

        // 处理请求
        c.Next()

        // 更新响应信息
        accessLog.Response = logger.ResponseInfo{
            Status:  c.Writer.Status(),
            Headers: logger.GetHeaders(c.Writer.Header()),
            Size:    int64(c.Writer.Size()),
            TimeMS:  float64(time.Since(start).Milliseconds()),
        }

        // 发送日志
        logManager.Send(accessLog)
    }
}
```

### 业务日志记录

服务层可以使用Logger记录业务日志：

```go
// 示例：记录用户创建事件
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
    // ... 业务逻辑 ...
    
    // 记录业务日志
    s.logger.LogBusiness(logger.LogLevelInfo, "user.create", user.ID, "", map[string]interface{}{
        "email": user.Email,
        "role": user.Role,
    })
    
    return nil
}
```

### 系统日志记录

基础设施组件可以记录系统日志：

```go
// 示例：记录数据库连接池状态
func (p *PostgresDB) logConnectionStatus() {
    stats := p.DB.Stats()
    
    p.logger.LogSystem(logger.LogLevelInfo, "database", "connection_pool_status", map[string]float64{
        "max_open_connections": float64(stats.MaxOpenConnections),
        "open_connections": float64(stats.OpenConnections),
        "in_use": float64(stats.InUse),
        "idle": float64(stats.Idle),
    })
}
```

## 与Filebeat集成

整合方式主要通过文件日志处理器与Filebeat进行集成：

1. 应用将结构化JSON日志写入文件(`logs/app.log`)
2. Filebeat监控并读取日志文件
3. Filebeat解析JSON并发送到Elasticsearch

关键配置点：

### 应用侧配置

```yaml
logger:
  processors:
    file:
      enabled: true
      path: "logs/app.log"
      rotation:
        max_size: "100MB"
        max_age: "7d"
        max_backups: 10
```

### Filebeat侧配置

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /app/logs/*.log
  json:
    keys_under_root: true
    add_error_key: true
    message_key: message
    overwrite_keys: true
```

## 日志查询最佳实践

### 基于日志类型查询

```
GET filebeat-logs-*/_search
{
  "query": {
    "term": {
      "type": "access"
    }
  }
}
```

### 按级别查询错误日志

```
GET filebeat-logs-*/_search
{
  "query": {
    "term": {
      "level": "ERROR"
    }
  }
}
```

### 查询特定API路径的访问日志

```
GET filebeat-logs-*/_search
{
  "query": {
    "bool": {
      "must": [
        { "term": { "type": "access" } },
        { "wildcard": { "request.path": "/api/v1/users*" } }
      ]
    }
  }
}
```

### 跟踪请求链路

```
GET filebeat-logs-*/_search
{
  "query": {
    "term": {
      "trace_id": "trace-xxxx-xxxx-xxxx"
    }
  },
  "sort": [
    { "timestamp": { "order": "asc" } }
  ]
}
```

### 统计API性能

```
GET filebeat-logs-*/_search
{
  "size": 0,
  "query": {
    "term": { "type": "access" }
  },
  "aggs": {
    "api_paths": {
      "terms": { "field": "request.path.keyword" },
      "aggs": {
        "avg_response_time": { "avg": { "field": "response.time_ms" } },
        "max_response_time": { "max": { "field": "response.time_ms" } },
        "p95_response_time": {
          "percentiles": {
            "field": "response.time_ms",
            "percents": [ 95 ]
          }
        }
      }
    }
  }
} 