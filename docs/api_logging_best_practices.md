# API开发与日志系统集成最佳实践

本文档提供了在cc-starship项目中进行API开发时，如何正确集成日志系统的最佳实践指南，帮助开发人员规范日志使用并充分利用日志系统的功能。

## 目录

1. [日志记录原则](#日志记录原则)
2. [日志级别使用指南](#日志级别使用指南)
3. [API层日志实践](#API层日志实践)
4. [服务层日志实践](#服务层日志实践)
5. [仓储层日志实践](#仓储层日志实践)
6. [分布式追踪集成](#分布式追踪集成)
7. [敏感信息处理](#敏感信息处理)
8. [错误处理与日志](#错误处理与日志)
9. [性能考量](#性能考量)
10. [最佳实践示例](#最佳实践示例)

## 日志记录原则

在API开发中记录日志时，请遵循以下基本原则：

1. **一致性**：使用统一的日志格式和结构
2. **完整性**：确保包含足够的上下文信息
3. **适度性**：避免过度记录或记录无用信息
4. **分类明确**：正确使用日志类型和级别
5. **可追踪性**：使用追踪ID关联相关日志
6. **安全性**：不记录敏感信息如密码、令牌等
7. **结构化**：使用结构化JSON格式便于查询分析

## 日志级别使用指南

### DEBUG

- 用于详细的开发调试信息
- **示例场景**：函数输入参数、中间计算结果、详细执行流程
- **使用建议**：仅在开发环境启用，生产环境禁用

```go
s.logger.Debug("Calculating user permission...")
```

### INFO

- 用于记录正常操作和状态变更
- **示例场景**：用户登录、资源创建、系统启动、定时任务执行
- **使用建议**：记录业务流程的关键节点

```go
s.logger.Info("User registered successfully")
```

### WARN

- 用于潜在问题或异常情况但不影响正常运行
- **示例场景**：配置错误但有默认值、性能下降、重试操作
- **使用建议**：应当关注并最终解决，但不需要立即处理

```go
s.logger.Warn("Database connection pool nearing capacity")
```

### ERROR

- 用于记录错误事件，但服务仍可继续运行
- **示例场景**：API调用失败、数据库操作异常、外部服务不可用
- **使用建议**：需要立即关注并处理

```go
s.logger.Error("Failed to update user profile")
```

### FATAL

- 用于导致应用终止的严重错误
- **示例场景**：关键配置缺失、无法恢复的系统故障
- **使用建议**：记录后通常伴随着应用退出

```go
s.logger.Fatal("Database connection failed after maximum retries")
```

## API层日志实践

API层主要关注请求处理和响应生成，日志应当反映HTTP交互细节。

### 中间件自动记录

使用中间件自动记录所有API请求和响应，无需在每个处理器中手动添加：

```go
// api/middleware/logging.go
func Logging(logManager *logger.LogManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 记录请求开始
        start := time.Now()
        
        // 生成请求ID如果不存在
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
            c.Header("X-Request-ID", requestID)
        }
        
        // 创建访问日志基本结构
        accessLog := createAccessLog(c, requestID)
        
        // 处理请求
        c.Next()
        
        // 更新响应信息
        accessLog.Response = createResponseInfo(c, time.Since(start))
        
        // 发送日志
        logManager.Send(accessLog)
    }
}
```

### API处理器中的关键操作日志

在特定API处理器中记录重要业务事件：

```go
// api/handlers/user_handler.go
func (h *UserHandler) CreateUser(c *gin.Context) {
    // 解析请求
    var req UserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Error("Invalid create user request")
        utils.ErrorResponse(c, errors.NewBadRequest("Invalid request body", err))
        return
    }
    
    // 记录操作开始
    requestID := c.GetHeader("X-Request-ID")
    h.logger.WithField("request_id", requestID).Info("Creating new user")
    
    // 处理业务逻辑...
    
    // 记录操作结果
    if err != nil {
        h.logger.WithField("request_id", requestID).Error("Failed to create user")
        // 返回错误响应...
        return
    }
    
    h.logger.WithField("request_id", requestID).
        WithField("user_id", user.ID).
        Info("User created successfully")
    
    // 返回成功响应...
}
```

## 服务层日志实践

服务层主要关注业务逻辑和领域操作，日志应当反映业务流程和状态变更。

### 业务操作日志

使用`LogBusiness`记录业务操作：

```go
// core/services/user_service.go
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
    traceID := extractTraceID(ctx)
    
    // 检查用户是否已存在
    existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
    if err == nil && existingUser != nil {
        s.logger.LogBusiness(
            logger.LogLevelWarn,
            "user.create.duplicate",
            "", // 匿名操作，无用户ID
            "",
            map[string]interface{}{
                "email": user.Email,
                "trace_id": traceID,
            },
        )
        return errors.NewConflict("User with this email already exists", nil)
    }
    
    // 创建用户
    if err := s.userRepo.Create(ctx, user); err != nil {
        s.logger.LogBusiness(
            logger.LogLevelError,
            "user.create.failed",
            "", 
            "",
            map[string]interface{}{
                "email": user.Email,
                "error": err.Error(),
                "trace_id": traceID,
            },
        )
        return errors.NewInternal("Failed to create user", err)
    }
    
    // 记录创建成功
    s.logger.LogBusiness(
        logger.LogLevelInfo,
        "user.create.success",
        "", 
        user.ID,
        map[string]interface{}{
            "email": user.Email,
            "role": user.Role,
            "trace_id": traceID,
        },
    )
    
    return nil
}
```

### 事务操作日志

记录事务开始、提交和回滚：

```go
// core/services/order_service.go
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
    traceID := extractTraceID(ctx)
    
    // 开始事务
    s.logger.Info("Starting transaction for order creation")
    
    err := s.db.Transaction(ctx, func(tx *sql.Tx) error {
        // 创建订单
        if err := s.orderRepo.CreateWithTx(ctx, tx, order); err != nil {
            s.logger.WithField("trace_id", traceID).Error("Failed to create order")
            return err
        }
        
        // 更新库存
        if err := s.inventoryRepo.UpdateWithTx(ctx, tx, order.Items); err != nil {
            s.logger.WithField("trace_id", traceID).Error("Failed to update inventory")
            return err
        }
        
        s.logger.WithField("trace_id", traceID).Info("Order operations completed successfully")
        return nil
    })
    
    if err != nil {
        s.logger.WithField("trace_id", traceID).Error("Transaction rolled back: " + err.Error())
        return err
    }
    
    s.logger.WithField("trace_id", traceID).Info("Transaction committed successfully")
    return nil
}
```

## 仓储层日志实践

仓储层主要关注数据操作和持久化，日志应当反映数据访问细节。

### 记录关键数据操作

```go
// repositories/postgres/user_repository.go
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    traceID := extractTraceID(ctx)
    
    r.logger.WithField("trace_id", traceID).Debug("Inserting user into database")
    
    query := `INSERT INTO users (...) VALUES (...)`
    _, err := r.db.ExecContext(ctx, query, ...)
    
    if err != nil {
        r.logger.LogSystem(
            logger.LogLevelError,
            "database",
            "query_execution_failed",
            map[string]float64{
                "execution_time_ms": float64(time.Since(start).Milliseconds()),
            },
            map[string]interface{}{
                "operation": "user_insert",
                "error": err.Error(),
                "trace_id": traceID,
            },
        )
        return err
    }
    
    r.logger.LogSystem(
        logger.LogLevelInfo,
        "database",
        "query_execution_success",
        map[string]float64{
            "execution_time_ms": float64(time.Since(start).Milliseconds()),
        },
        map[string]interface{}{
            "operation": "user_insert",
            "user_id": user.ID,
            "trace_id": traceID,
        },
    )
    
    return nil
}
```

### 记录查询性能

```go
// repositories/postgres/user_repository.go
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    traceID := extractTraceID(ctx)
    start := time.Now()
    
    query := `SELECT * FROM users WHERE email = $1`
    row := r.db.QueryRowContext(ctx, query, email)
    
    // 解析结果...
    
    elapsed := time.Since(start)
    if elapsed > slowQueryThreshold {
        r.logger.LogSystem(
            logger.LogLevelWarn,
            "database",
            "slow_query",
            map[string]float64{
                "execution_time_ms": float64(elapsed.Milliseconds()),
                "threshold_ms": float64(slowQueryThreshold.Milliseconds()),
            },
            map[string]interface{}{
                "operation": "find_user_by_email",
                "trace_id": traceID,
            },
        )
    } else {
        r.logger.LogSystem(
            logger.LogLevelInfo,
            "database",
            "query_execution",
            map[string]float64{
                "execution_time_ms": float64(elapsed.Milliseconds()),
            },
            map[string]interface{}{
                "operation": "find_user_by_email",
                "trace_id": traceID,
            },
        )
    }
    
    return user, nil
}
```

## 分布式追踪集成

确保跨服务和组件的请求可追踪。

### 创建和传递追踪ID

```go
// pkg/utils/tracing.go
func GenerateTraceID() string {
    return uuid.New().String()
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
    return context.WithValue(ctx, traceIDKey, traceID)
}

func ExtractTraceID(ctx context.Context) string {
    if traceID, ok := ctx.Value(traceIDKey).(string); ok {
        return traceID
    }
    return ""
}
```

### 在API层注入追踪ID

```go
// api/middleware/tracing.go
func Tracing() gin.HandlerFunc {
    return func(c *gin.Context) {
        traceID := c.GetHeader("X-Trace-ID")
        if traceID == "" {
            traceID = utils.GenerateTraceID()
            c.Header("X-Trace-ID", traceID)
        }
        
        ctx := utils.WithTraceID(c.Request.Context(), traceID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}
```

### 使用追踪ID关联日志

```go
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
    traceID := utils.ExtractTraceID(ctx)
    
    s.logger.WithField("trace_id", traceID).Debug("Getting user by ID")
    
    // 业务逻辑...
    
    s.logger.WithField("trace_id", traceID).
        WithField("user_id", id).
        Info("User retrieved successfully")
    
    return user, nil
}
```

## 敏感信息处理

避免记录敏感信息或使用适当的脱敏处理。

### 定义敏感字段脱敏规则

```go
// pkg/logger/sanitizer.go
var sensitiveFields = map[string]bool{
    "password": true,
    "token": true,
    "secret": true,
    "credit_card": true,
    "ssn": true,
}

func SanitizeData(data map[string]interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    
    for k, v := range data {
        if sensitiveFields[strings.ToLower(k)] {
            result[k] = "********"
        } else if nestedMap, ok := v.(map[string]interface{}); ok {
            result[k] = SanitizeData(nestedMap)
        } else {
            result[k] = v
        }
    }
    
    return result
}
```

### 在日志记录时应用脱敏

```go
// api/handlers/auth_handler.go
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 错误处理...
        return
    }
    
    // 注意我们只记录了email字段，没有记录password
    h.logger.WithField("email", req.Email).Info("Login attempt")
    
    // 业务逻辑...
}
```

## 错误处理与日志

结合错误处理与日志记录，确保错误上下文完整。

### 定义带上下文的错误

```go
// pkg/errors/errors.go
type AppError struct {
    Code    int    `json:"-"`       // HTTP状态码
    Message string `json:"message"` // 用户友好的错误消息
    Details string `json:"details"` // 详细错误信息
    Err     error  `json:"-"`       // 原始错误
}
```

### 记录错误并保留上下文

```go
// core/services/user_service.go
func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
    traceID := utils.ExtractTraceID(ctx)
    logger := s.logger.WithField("trace_id", traceID).WithField("user_id", user.ID)
    
    // 检查用户是否存在
    existingUser, err := s.userRepo.FindByID(ctx, user.ID)
    if err != nil {
        logger.WithError(err).Error("Failed to find user for update")
        return errors.NewNotFound("User not found", err)
    }
    
    // 更新用户
    if err := s.userRepo.Update(ctx, user); err != nil {
        logger.WithError(err).Error("Failed to update user in database")
        return errors.NewInternal("Failed to update user", err)
    }
    
    logger.Info("User updated successfully")
    return nil
}
```

### 在API层处理错误并记录

```go
// api/utils/response.go
func ErrorResponse(c *gin.Context, err error) {
    var status int
    var errorResponse interface{}
    
    // 检查是否为应用错误
    if appErr, ok := errors.IsAppError(err); ok {
        status = appErr.StatusCode()
        errorResponse = map[string]string{
            "message": appErr.Message,
        }
        
        // 包含详情（如果有）
        if appErr.Details != "" {
            errorResponse = map[string]string{
                "message": appErr.Message,
                "details": appErr.Details,
            }
        }
        
        // 记录详细错误
        logger := GetRequestLogger(c)
        logger.WithField("error_code", status).
               WithField("error_message", appErr.Message).
               WithError(appErr.Err).
               Error("API error response")
    } else {
        // 对于通用错误，使用内部服务器错误
        status = http.StatusInternalServerError
        errorResponse = map[string]string{
            "message": "An internal server error occurred",
        }
        
        // 记录未处理的错误
        logger := GetRequestLogger(c)
        logger.WithError(err).Error("Unhandled internal error")
    }
    
    // 创建错误响应
    response := Response{
        Success: false,
        Error:   errorResponse,
    }
    
    // 发送响应
    c.JSON(status, response)
}
```

## 性能考量

避免日志系统成为性能瓶颈。

### 使用异步日志处理

我们的LogManager已经实现了异步处理机制：

```go
// pkg/logger/manager.go
func (m *LogManager) Send(log interface{}) {
    select {
    case m.logChan <- log:
    default:
        // 如果通道已满，丢弃日志
        fmt.Printf("Log channel is full, dropping log\n")
    }
}

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
```

### 避免高频调试日志

在性能敏感路径上避免大量DEBUG级别日志：

```go
// 只在DEBUG启用时生成复杂日志信息
if logger.IsLevelEnabled(logger.LogLevelDebug) {
    detailedInfo := generateDetailedDebugInfo() // 昂贵的操作
    logger.WithField("details", detailedInfo).Debug("Detailed operation state")
}
```

### 批量处理日志

对于高流量路径，考虑批量记录：

```go
// infrastructure/metrics/request_metrics.go
type RequestMetricsCollector struct {
    mu          sync.Mutex
    metrics     map[string]*PathMetrics
    flushTicker *time.Ticker
    logger      *logger.Logger
}

func (c *RequestMetricsCollector) RecordRequest(path string, duration time.Duration, statusCode int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // 更新内存中的指标...
}

func (c *RequestMetricsCollector) flush() {
    c.mu.Lock()
    metricsSnapshot := c.metrics
    c.metrics = make(map[string]*PathMetrics)
    c.mu.Unlock()
    
    // 批量记录所有路径的聚合指标
    systemLog := &logger.SystemLog{
        BaseLog: logger.BaseLog{
            LogID:     logger.GenerateLogID(),
            Timestamp: time.Now(),
            Type:      logger.LogTypeSystem,
            Level:     logger.LogLevelInfo,
            Service:   "cc-starship",
        },
        Component: "api",
        Event:     "request_metrics",
    }
    
    // 记录聚合指标...
    c.logger.Send(systemLog)
}
```

## 最佳实践示例

以下是一个完整的API端点实现，展示如何将日志系统集成到API开发中。

### 用户注册API实现示例

```go
// api/handlers/auth_handler.go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/majiayu000/cc-starship/internal/core/ports"
    "github.com/majiayu000/cc-starship/pkg/errors"
    "github.com/majiayu000/cc-starship/pkg/logger"
    "github.com/majiayu000/cc-starship/pkg/utils"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=6"`
    FirstName string `json:"firstName" binding:"required"`
    LastName  string `json:"lastName" binding:"required"`
}

// Register 处理用户注册
func (h *AuthHandler) Register(c *gin.Context) {
    traceID := utils.ExtractTraceID(c.Request.Context())
    requestID := c.GetHeader("X-Request-ID")
    logger := h.logger.WithField("trace_id", traceID).WithField("request_id", requestID)
    
    logger.Info("Processing user registration request")
    
    // 解析请求体
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.WithError(err).Warn("Invalid registration request body")
        utils.ErrorResponse(c, errors.NewBadRequest("Invalid request body", err))
        return
    }
    
    // 记录注册尝试（注意不记录密码）
    logger.WithField("email", req.Email).
           WithField("first_name", req.FirstName).
           WithField("last_name", req.LastName).
           Info("User registration attempt")
    
    // 注册用户
    user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName)
    if err != nil {
        logger.WithError(err).Error("User registration failed")
        utils.ErrorResponse(c, err)
        return
    }
    
    // 生成令牌
    token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        logger.WithError(err).Error("Token generation failed after registration")
        utils.ErrorResponse(c, errors.NewInternal("Failed to generate token", err))
        return
    }
    
    // 记录成功注册
    logger.WithField("user_id", user.ID).Info("User registered successfully")
    
    // 构造响应体
    response := AuthResponse{
        Token: token,
        User: map[string]interface{}{
            "id":        user.ID,
            "email":     user.Email,
            "firstName": user.FirstName,
            "lastName":  user.LastName,
            "role":      user.Role,
        },
    }
    
    // 发送响应
    utils.JSONResponse(c, http.StatusCreated, response)
    
    // 记录业务事件日志
    h.logManager.Send(&logger.BusinessLog{
        BaseLog: logger.BaseLog{
            LogID:     logger.GenerateLogID(),
            Timestamp: utils.Now(),
            Type:      logger.LogTypeBusiness,
            Level:     logger.LogLevelInfo,
            Service:   "cc-starship",
            TraceID:   traceID,
            SpanID:    logger.GenerateSpanID(),
            Message:   "User registration completed",
        },
        Action:     "user.register",
        UserID:     user.ID,
        ResourceID: user.ID,
        Details: map[string]interface{}{
            "email": user.Email,
            "role":  user.Role,
        },
    })
}
```

通过在API开发中遵循这些最佳实践，您将确保日志系统提供全面的系统可观测性，并支持问题诊断和性能优化。 