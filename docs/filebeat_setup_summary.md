# Filebeat 安装配置与问题解决全记录

本文档详细记录了Filebeat的安装、配置和问题解决过程，作为团队内部技术参考。

## 目录

1. [概述](#概述)
2. [环境信息](#环境信息)
3. [初始安装](#初始安装)
4. [基础配置](#基础配置)
5. [Docker部署](#Docker部署)
6. [索引模板与数据流问题](#索引模板与数据流问题)
7. [日志查看方法](#日志查看方法)
8. [测试日志生成](#测试日志生成)
9. [常见问题与解决方案](#常见问题与解决方案)
10. [配置文件参考](#配置文件参考)

## 概述

Filebeat是Elastic Stack的一部分，作为轻量级日志收集器，用于从文件系统收集日志并发送到Elasticsearch或其他输出目标。本文档记录了为cc-starship项目设置Filebeat的完整过程，包括安装、配置、Docker部署以及遇到的问题和解决方案。

## 环境信息

- **操作系统**: macOS (darwin 24.1.0)
- **Docker版本**: Docker 版本20.10+
- **Filebeat版本**: 8.17.3
- **Elasticsearch版本**: 8.17.3
- **应用服务**: cc-starship（Go后端API服务）

## 初始安装

### 拉取Filebeat Docker镜像

首次尝试使用的命令存在问题：
```bash
docker pull elastic/filebeat
```

修正后的命令，使用正确的镜像路径：
```bash
docker pull docker.elastic.co/beats/filebeat:8.17.3
```

> **注意**: Elastic官方镜像存放在`docker.elastic.co`仓库，而非Docker Hub的`elastic`命名空间下。

## 基础配置

### 创建配置目录

```bash
mkdir -p filebeat && cd filebeat
```

### 创建基础配置文件

Filebeat需要一个配置文件来指定输入源、处理器和输出目标。创建`filebeat.yml`文件如下：

```yaml
# 输入配置
filebeat.inputs:
# 后端应用日志输入（JSON格式）
- type: log
  enabled: true
  paths:
    # 应用程序日志路径，基于应用配置文件中的设置
    - /app/logs/*.log
    - /usr/share/filebeat/project_logs/*.log
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
  tags: ["application", "backend", "api"]
  
  # 多行处理
  multiline:
    pattern: '^{'
    negate: true
    match: after
    max_lines: 500
    timeout: 5s

# 处理器配置
processors:
  # 添加元数据
  - add_host_metadata: ~
  
  # 字段处理和增强
  - add_fields:
      target: ''
      fields:
        environment: development
  
  # 时间戳处理（北京时间 UTC+8）
  - timestamp:
      field: timestamp
      layouts:
        - '2006-01-02T15:04:05Z'
        - '2006-01-02T15:04:05.999Z'
        - '2006-01-02T15:04:05'
        - '2006-01-02 15:04:05'
        - '2006/01/02 15:04:05'
        - 'UNIX'
        - 'UNIX_MS'
      target_field: "@timestamp"
      timezone: "Asia/Shanghai"
      ignore_missing: true
  
  # 添加时区信息
  - add_fields:
      target: ''
      fields:
        timezone: "Asia/Shanghai"

# 时间设置
setup.timezone: "Asia/Shanghai"

# 索引模板设置
setup.template.enabled: true
setup.template.name: "app-logs"
setup.template.pattern: "app-logs-*"
setup.template.settings:
  index.number_of_shards: 1
  index.number_of_replicas: 1
setup.template.overwrite: true

# 索引生命周期管理
setup.ilm.enabled: true
setup.ilm.rollover_alias: "app-logs"
setup.ilm.pattern: "{now/d}-000001"
setup.ilm.policy_name: "app-logs-policy"

# 输出配置
output.elasticsearch:
  hosts: ["localhost:9200"]
  username: "elastic"
  password: "changeme"
  index: "app-logs-%{+yyyy.MM.dd}"
  pipeline: "timezone-conversion"

# 启用监控
monitoring.enabled: true
monitoring.elasticsearch:
  hosts: ["localhost:9200"]
  username: "elastic"
  password: "changeme"

# 日志配置
logging:
  level: info
  to_stderr: true
  metrics.enabled: false
  
# 设置更大的缓冲区和批处理
queue.mem:
  events: 4096
  flush.min_events: 512
  flush.timeout: 5s
```

## Docker部署

### 基本Docker命令

可以使用以下命令运行单独的Filebeat容器：

```bash
docker run -d \
  --name=filebeat \
  --user=root \
  --volume="$(pwd)/filebeat.yml:/usr/share/filebeat/filebeat.yml:ro" \
  --volume="/var/log:/var/log:ro" \
  --volume="/var/lib/docker/containers:/var/lib/docker/containers:ro" \
  --network="elastic" \
  docker.elastic.co/beats/filebeat:8.17.3 \
  filebeat -e -strict.perms=false
```

### 创建Docker Compose配置

但更建议使用Docker Compose来管理Filebeat容器。创建`docker-compose.yml`文件：

```yaml
version: '3'

services:
  filebeat:
    image: docker.elastic.co/beats/filebeat:8.17.3
    user: root  # 需要 root 权限来读取日志文件
    network_mode: "host"  # 使用主机网络模式，可以直接访问本机服务
    volumes:
      - ./filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - ../logs:/app/logs:ro  # 从项目根目录挂载应用日志
      - ../../logs:/usr/share/filebeat/project_logs:ro  # 备用日志路径
    environment:
      - strict.perms=false
    command: filebeat -e
    restart: unless-stopped
```

### 创建Docker网络

如果需要与其他ELK容器通信，创建共享网络：

```bash
docker network create elastic
```

### 启动Filebeat容器

```bash
docker-compose up -d
```

## 索引模板与数据流问题

### 错误现象

在启动Filebeat后，我们遇到了以下错误：

```
Failed to connect to backoff(elasticsearch(http://localhost:9200)): Connection marked as failed because the onConnect callback failed: error loading template: failed to put data stream: could not put data stream: 400 Bad Request: {"error":{"root_cause":[{"type":"illegal_argument_exception","reason":"no matching index template found for data stream [app-logs]"}],"type":"illegal_argument_exception","reason":"no matching index template found for data stream [app-logs]"},"status":400}
```

这表明Filebeat试图创建一个数据流(data stream)，但没有找到匹配的索引模板。

### 问题分析

在Elasticsearch 8.x中：
1. 数据流需要特定格式的索引模板
2. Filebeat 8.x默认使用数据流（即使设置了传统索引）
3. 索引模板需要特定的命名规则才能与数据流关联

### 解决方案

我们通过彻底修改配置，采用传统索引而非数据流方式解决问题：

```yaml
# 使用传统方式配置索引模板
setup.template.enabled: true
setup.template.name: "filebeat-logs"  # 更改模板名称，避免与数据流命名冲突
setup.template.pattern: "filebeat-logs-*"  # 匹配索引名称的模式
setup.template.type: "index"  # 明确指定为索引模板，非数据流
setup.template.settings:
  index.number_of_shards: 1
  index.number_of_replicas: 1
setup.template.overwrite: true
# 禁用组件模板 (8.x feature)
setup.template.legacy.enabled: true  # 使用传统索引模板

# 完全禁用ILM及数据流
setup.ilm.enabled: false
setup.ilm.check_exists: false
setup.ilm.overwrite: false

# 输出配置 - 使用与模板匹配的索引命名模式
output.elasticsearch:
  hosts: ["localhost:9200"]
  username: "elastic"
  password: "changeme"
  index: "filebeat-logs-%{+yyyy.MM.dd}"  # 索引名与模板模式匹配
  pipeline: "timezone-conversion"
  # 禁用数据流
  allow_older_versions: true
```

### 时区处理

为了确保日志中的时间戳以北京时间(UTC+8)显示，我们使用了时间戳处理器和一个Elasticsearch管道：

#### 创建管道定义文件(pipeline.json)

```json
{
  "description": "Pipeline for timezone conversion to Asia/Shanghai",
  "processors": [
    {
      "date": {
        "field": "@timestamp",
        "timezone": "Asia/Shanghai",
        "formats": [
          "ISO8601",
          "UNIX",
          "UNIX_MS"
        ]
      }
    }
  ]
}
```

#### 创建设置管道的脚本(setup-pipeline.sh)

```bash
#!/bin/bash

# 创建 pipeline
curl -X PUT "localhost:9200/_ingest/pipeline/timezone-conversion" \
  -H "Content-Type: application/json" \
  -d @pipeline.json
```

## 日志查看方法

有多种方式查看Filebeat收集的日志：

### 1. 查看Filebeat容器日志

```bash
docker-compose logs filebeat | tail -n 50
```

### 2. 直接查看源日志文件

```bash
cat ../logs/app.log
```

### 3. 通过Filebeat容器查看日志

```bash
docker-compose exec filebeat cat /app/logs/app.log
```

### 4. 使用文本处理工具解析JSON日志

```bash
cat ../logs/app.log | jq '.'                    # 格式化显示所有日志
cat ../logs/app.log | jq 'select(.type=="access")'    # 只显示访问日志
cat ../logs/app.log | jq 'select(.level=="ERROR")'    # 只显示错误日志
```

### 5. 检查Elasticsearch中的索引

```bash
curl -u elastic:changeme -X GET "localhost:9200/_cat/indices?v"
```

### 6. 在Elasticsearch中搜索日志

```bash
curl -u elastic:changeme -X GET "localhost:9200/filebeat-logs-*/_search?pretty" -H 'Content-Type: application/json' -d'{"size": 1}'
```

## 测试日志生成

我们创建了一个脚本用于生成测试日志，用于验证Filebeat配置是否正确：

```bash
#!/bin/bash

# 创建日志目录
mkdir -p ../logs

# 生成访问日志
generate_access_log() {
    local timestamp=$(date +"%Y-%m-%dT%H:%M:%S.%3NZ")
    local log_id=$(uuidgen)
    cat << EOF >> ../logs/app.log
{
    "log_id": "$log_id",
    "timestamp": "$timestamp",
    "type": "access",
    "level": "INFO",
    "service": "cc-starship",
    "trace_id": "trace-$(uuidgen)",
    "span_id": "span-$(uuidgen | cut -c1-8)",
    "request": {
        "method": "GET",
        "path": "/api/v1/users",
        "query": "page=1&size=10",
        "headers": {
            "User-Agent": "Mozilla/5.0",
            "Content-Type": "application/json"
        },
        "client_ip": "192.168.1.1",
        "user_agent": "curl/7.64.1"
    },
    "response": {
        "status": 200,
        "headers": {},
        "size": 1024,
        "time_ms": 45.2
    }
}
EOF
}

# 生成业务日志
generate_business_log() {
    local timestamp=$(date +"%Y-%m-%dT%H:%M:%S.%3NZ")
    local log_id=$(uuidgen)
    cat << EOF >> ../logs/app.log
{
    "log_id": "$log_id",
    "timestamp": "$timestamp",
    "type": "business",
    "level": "INFO",
    "service": "cc-starship",
    "trace_id": "trace-$(uuidgen)",
    "span_id": "span-$(uuidgen | cut -c1-8)",
    "action": "user.create",
    "user_id": "user-123",
    "resource_id": "resource-456",
    "details": {
        "custom_field": "value"
    }
}
EOF
}

# 生成系统日志
generate_system_log() {
    local timestamp=$(date +"%Y-%m-%dT%H:%M:%S.%3NZ")
    local log_id=$(uuidgen)
    cat << EOF >> ../logs/app.log
{
    "log_id": "$log_id",
    "timestamp": "$timestamp",
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
EOF
}

# 生成一些测试日志
echo "Generating test logs..."
for i in {1..5}; do
    generate_access_log
    sleep 1
    generate_business_log
    sleep 1
    generate_system_log
    sleep 1
done

echo "Logs generated in ../logs/app.log"
```

使用方法：

```bash
chmod +x generate-logs.sh
./generate-logs.sh
```

## 常见问题与解决方案

### 1. Filebeat无法连接Elasticsearch

**症状**: 日志中出现连接错误，显示无法连接到Elasticsearch

**解决方案**:
- 确认Elasticsearch服务已启动并可访问
- 检查网络配置，确保Filebeat容器能够访问Elasticsearch（如使用host网络或共享网络）
- 验证Elasticsearch凭据是否正确

### 2. JSON解析错误

**症状**: 日志中出现"Error decoding JSON"错误

**解决方案**:
- 检查日志格式是否为有效的JSON
- 确认multiline设置正确，避免将多行日志拆分
- 调整JSON解析设置，如`message_key`和`keys_under_root`参数

### 3. 索引模板和数据流冲突

**症状**: 日志中出现"no matching index template found for data stream"错误

**解决方案**:
- 使用传统索引而非数据流（如本文档所述）
- 确保索引名称和模板名称一致
- 明确设置`setup.template.type: "index"`和`setup.ilm.enabled: false`

### 4. 时区不正确

**症状**: 日志中的时间戳不是期望的时区（如没有显示为北京时间）

**解决方案**:
- 在processors中添加timestamp处理器
- 设置`setup.timezone`参数
- 使用Elasticsearch管道进行时区转换
- 确保应用程序生成的日志带有正确的时区信息

### 5. 权限问题

**症状**: 日志中出现文件访问权限错误

**解决方案**:
- 在Docker中使用`--user: root`运行Filebeat
- 确保日志目录有正确的读取权限
- 使用`strict.perms=false`环境变量

## 配置文件参考

所有配置文件均已保存在`filebeat`目录下：

- `filebeat.yml` - Filebeat主要配置文件
- `docker-compose.yml` - Docker Compose配置
- `pipeline.json` - Elasticsearch管道配置
- `setup-pipeline.sh` - 设置管道的脚本
- `generate-logs.sh` - 生成测试日志的脚本

所有这些配置文件都已在本文档中详细说明，可以作为未来项目的参考。 