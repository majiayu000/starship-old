# 日志实现指南：从应用到 ELK

## 目录

1. [日志格式](#1-日志格式)
2. [日志收集流程](#2-日志收集流程)
3. [配置指南](#3-配置指南)
4. [最佳实践](#4-最佳实践)
5. [监控和维护](#5-监控和维护)

## 1. 日志格式

我们的应用使用结构化的 JSON 日志格式，包含以下几种类型：

### 1.1 访问日志 (AccessLog)
```json
{
  "log_id": "uuid",
  "timestamp": "2024-01-01T12:00:00Z",
  "type": "access",
  "level": "INFO",
  "service": "cc-starship",
  "trace_id": "trace-id",
  "span_id": "span-id",
  "request": {
    "method": "GET",
    "path": "/api/v1/users",
    "query": "page=1&size=10",
    "headers": {
      "User-Agent": "...",
      "Content-Type": "application/json"
    },
    "client_ip": "192.168.1.1",
    "user_agent": "..."
  },
  "response": {
    "status": 200,
    "headers": {},
    "size": 1024,
    "time_ms": 45.2
  }
}
```

### 1.2 业务日志 (BusinessLog)
```json
{
  "log_id": "uuid",
  "timestamp": "2024-01-01T12:00:00Z",
  "type": "business",
  "level": "INFO",
  "service": "cc-starship",
  "trace_id": "trace-id",
  "span_id": "span-id",
  "action": "user.create",
  "user_id": "user-123",
  "resource_id": "resource-456",
  "details": {
    "custom_field": "value"
  }
}
```

### 1.3 系统日志 (SystemLog)
```json
{
  "log_id": "uuid",
  "timestamp": "2024-01-01T12:00:00Z",
  "type": "system",
  "level": "INFO",
  "service": "cc-starship",
  "component": "database",
  "event": "connection_pool_status",
  "metrics": {
    "active_connections": 10,
    "idle_connections": 5
  }
}
```

## 2. 日志收集流程

日志收集采用以下流程：

```
应用程序 → 日志文件 → Filebeat → Elasticsearch → Kibana
```

1. **应用程序生成日志**
   - 使用结构化 JSON 格式
   - 按类型分类（访问、业务、系统）
   - 包含必要的元数据（时间戳、追踪ID等）

2. **Filebeat 收集**
   - 监控日志文件
   - 解析 JSON 格式
   - 添加主机元数据
   - 添加 Docker 元数据（如果在容器环境中）

3. **Elasticsearch 存储**
   - 使用动态索引（按日期）
   - 应用适当的映射
   - 实现数据生命周期管理

4. **Kibana 可视化**
   - 实时监控和分析
   - 自定义仪表板
   - 告警配置

## 3. 配置指南

### 3.1 Filebeat 配置

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/*.log
  json.keys_under_root: true
  json.add_error_key: true
  json.message_key: message

processors:
  - add_host_metadata: ~
  - add_docker_metadata: ~

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  index: "logs-%{+yyyy.MM.dd}"
```

### 3.2 Docker 部署

```yaml
version: '3'
services:
  filebeat:
    image: docker.elastic.co/beats/filebeat:8.17.3
    user: root
    volumes:
      - ./filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - /var/log:/var/log:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      - strict.perms=false
    command: filebeat -e
    networks:
      - elastic
```

## 4. 最佳实践

### 4.1 日志生成
- 使用统一的日志格式
- 包含必要的上下文信息
- 合理的日志级别使用
- 敏感信息脱敏

### 4.2 日志收集
- 合理的文件轮转策略
- 适当的资源限制
- 错误处理和重试机制
- 监控和告警设置

### 4.3 存储和查询
- 合理的索引策略
- 数据生命周期管理
- 查询优化
- 备份策略

## 5. 监控和维护

### 5.1 监控指标
- Filebeat 状态
- 日志收集延迟
- 错误率
- 资源使用情况

### 5.2 常见问题处理
- 日志解析错误
- 连接问题
- 性能问题
- 存储空间管理

### 5.3 日常维护
- 日志文件清理
- 配置更新
- 性能优化
- 安全更新 