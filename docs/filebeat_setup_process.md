# Filebeat 配置过程记录

## 目录

1. [问题背景](#1-问题背景)
2. [配置过程](#2-配置过程)
3. [遇到的问题和解决方案](#3-遇到的问题和解决方案)
4. [最终配置](#4-最终配置)
5. [验证和测试](#5-验证和测试)

## 1. 问题背景

在配置 Filebeat 过程中，我们遇到了以下需求和问题：

1. 需要收集应用程序的 JSON 格式日志
2. 需要确保时间戳显示为北京时间（UTC+8）
3. 需要正确处理 JSON 格式的日志解析
4. 需要确保与 Elasticsearch 的连接正常工作

## 2. 配置过程

### 2.1 初始配置尝试

最初配置遇到了几个问题：
- JSON 解析错误
- 时区显示不正确
- Elasticsearch 连接问题

### 2.2 配置优化过程

1. **日志路径配置**
   ```yaml
   paths:
     - /app/logs/*.log
     - /usr/share/filebeat/project_logs/*.log
   ```

2. **JSON 解析配置**
   ```yaml
   json:
     keys_under_root: true
     add_error_key: true
     message_key: message
     overwrite_keys: true
   ```

3. **时区配置**
   ```yaml
   timezone: "Asia/Shanghai"
   ```

4. **Docker 挂载配置**
   ```yaml
   volumes:
     - ./filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
     - ../logs:/app/logs:ro
   ```

## 3. 遇到的问题和解决方案

### 3.1 JSON 解析错误

**问题描述**：
```
Error decoding JSON: invalid character 'F' looking for beginning of value
```

**解决方案**：
- 调整 JSON 解析配置
- 添加多行处理支持
- 确保日志格式符合 JSON 规范

### 3.2 时区问题

**问题描述**：
- 日志显示的是 UTC 时间，需要显示北京时间

**解决方案**：
1. 添加时区处理器：
   ```yaml
   processors:
     - timestamp:
         field: timestamp
         layouts:
           - '2006-01-02T15:04:05Z'
           - '2006-01-02T15:04:05.999Z'
         target_field: "@timestamp"
         timezone: "Asia/Shanghai"
   ```

2. 在多个层级设置时区：
   - 全局时区设置
   - 输出配置时区设置
   - 处理器时区设置

### 3.3 Elasticsearch 连接问题

**问题描述**：
```
Failed to connect to Elasticsearch... Error: Get "http://elasticsearch:9200": EOF
```

**解决方案**：
1. 修改网络配置为 host 模式
2. 使用 localhost 替代容器名
3. 确保 Elasticsearch 服务可访问

## 4. 最终配置

### 4.1 Filebeat 配置

```yaml
# 输入配置
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
  fields:
    log_type: application
    service: cc-starship
    timezone: Asia/Shanghai
  fields_under_root: true

# 处理器配置
processors:
  - add_host_metadata: ~
  - timestamp:
      field: timestamp
      layouts:
        - '2006-01-02T15:04:05Z'
        - '2006-01-02T15:04:05.999Z'
      target_field: "@timestamp"
      timezone: "Asia/Shanghai"

# 输出配置
output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "app-logs-%{+yyyy.MM.dd}"
```

### 4.2 Docker Compose 配置

```yaml
version: '3'
services:
  filebeat:
    image: docker.elastic.co/beats/filebeat:8.17.3
    user: root
    network_mode: "host"
    volumes:
      - ./filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - ../logs:/app/logs:ro
    environment:
      - strict.perms=false
    command: filebeat -e
    restart: unless-stopped
```

## 5. 验证和测试

### 5.1 验证步骤

1. 检查 Filebeat 日志确认配置加载正常
2. 验证日志文件被正确监控
3. 确认时间戳显示正确（北京时间）
4. 验证 JSON 解析正确
5. 检查 Elasticsearch 索引创建情况

### 5.2 测试方法

1. 生成测试日志：
   ```bash
   echo '{"timestamp":"2024-01-01T12:00:00Z","message":"test"}' >> app.log
   ```

2. 检查 Elasticsearch 索引：
   ```bash
   curl -X GET "localhost:9200/app-logs-*/_search"
   ```

3. 验证时间转换：
   ```bash
   # 检查转换后的时间戳是否为北京时间
   curl -X GET "localhost:9200/app-logs-*/_search" -H 'Content-Type: application/json' -d'
   {
     "query": {
       "match_all": {}
     },
     "sort": [
       {
         "@timestamp": {
           "order": "desc"
         }
       }
     ]
   }'
   ```

## 总结

通过以上配置和调整，我们成功实现了：
1. JSON 日志的正确解析
2. 北京时间的显示
3. 与 Elasticsearch 的稳定连接
4. 日志的可靠收集和转发

后续可以考虑的优化点：
1. 添加更多的日志处理逻辑
2. 实现日志轮转策略
3. 添加监控告警机制
4. 优化性能配置 