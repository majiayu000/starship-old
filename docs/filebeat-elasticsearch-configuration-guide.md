# Filebeat 数据流配置指南

*针对 Filebeat 8.17.3 与 Elasticsearch 8.17.3*

## 目录

- [概述](#概述)
- [架构](#架构)
- [环境要求](#环境要求)
- [安装步骤](#安装步骤)
- [配置详解](#配置详解)
- [验证部署](#验证部署)
- [故障排除](#故障排除)
- [性能优化](#性能优化)
- [运维建议](#运维建议)
- [参考资料](#参考资料)

## 概述

Filebeat 数据流是 Elastic Stack 中用于日志采集的高效方式，提供了简化的索引管理和优化的写入性能。本文档详细说明如何配置 Filebeat 8.17.3 与 Elasticsearch 8.17.3 一起使用数据流功能。

数据流的主要优势：

- **简化索引管理**：自动处理索引生命周期
- **优化写入性能**：新数据总是写入最新索引
- **统一数据访问**：通过单一数据流名称访问所有相关索引
- **无缝扩展**：随着数据量增长轻松扩展
- **自动化运维**：减少手动索引管理工作

## 架构

数据流的基本架构如下：

```
┌───────────────┐    ┌───────────────┐    ┌─────────────────────────────┐
│ 应用服务器    │    │    Filebeat    │    │       Elasticsearch         │
│ ┌───────────┐ │    │ ┌───────────┐ │    │ ┌─────────┐  ┌────────────┐ │
│ │ 应用日志  │─┼───▶│ │输入/处理器│─┼───▶│ │数据流   │─▶│后备索引    │ │
│ └───────────┘ │    │ └───────────┘ │    │ └─────────┘  └────────────┘ │
└───────────────┘    └───────────────┘    └─────────────────────────────┘
```

## 环境要求

### 软件版本
- Elasticsearch 8.17.3
- Filebeat 8.17.3
- Docker & Docker Compose (可选，用于容器化部署)

### 系统要求
- 最小 2GB 内存
- 2 CPU 核心以上
- 足够的磁盘空间用于日志存储

### 权限要求
- Elasticsearch 集群管理员权限（用于创建索引模板和ILM策略）
- 访问待监控服务器上日志文件的权限

## 安装步骤

### 1. 准备环境

创建工作目录并设置配置文件：

```bash
mkdir -p filebeat-config/
cd filebeat-config/
```

### 2. 创建 docker-compose.yml（可选，用于容器化部署）

```yaml
version: '3'
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.17.3
    environment:
      - discovery.type=single-node
      - ES_JAVA_OPTS=-Xms1g -Xmx1g
      - xpack.security.enabled=true
      - ELASTIC_PASSWORD=changeme
    ports:
      - 9200:9200
    volumes:
      - es_data:/usr/share/elasticsearch/data
    networks:
      - elastic

  filebeat:
    image: docker.elastic.co/beats/filebeat:8.17.3
    user: root
    volumes:
      - ./filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - ./project_logs:/usr/share/filebeat/project_logs:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    networks:
      - elastic
    depends_on:
      - elasticsearch
    command: filebeat -e -strict.perms=false

networks:
  elastic:

volumes:
  es_data:
```

### 3. 设置 ILM 策略（可选：如果尚未设置）

创建 `filebeat-ilm.json` 文件：

```json
{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_age": "1d",
            "max_size": "5gb"
          },
          "set_priority": {
            "priority": 100
          }
        }
      },
      "warm": {
        "min_age": "3d",
        "actions": {
          "shrink": {
            "number_of_shards": 1
          },
          "forcemerge": {
            "max_num_segments": 1
          },
          "set_priority": {
            "priority": 50
          }
        }
      },
      "cold": {
        "min_age": "30d",
        "actions": {
          "freeze": {},
          "set_priority": {
            "priority": 0
          }
        }
      },
      "delete": {
        "min_age": "90d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}
```

应用 ILM 策略：

```bash
curl -X PUT "localhost:9200/_ilm/policy/filebeat" \
  -H "Content-Type: application/json" \
  -u elastic:changeme \
  -d @filebeat-ilm.json
```

### 4. 创建索引模板（可选：如果需要自定义字段映射）

创建 `filebeat-template.json` 文件：

```json
{
  "index_patterns": ["logs-*-*"],
  "template": {
    "settings": {
      "index.lifecycle.name": "filebeat",
      "index.lifecycle.rollover_alias": "filebeat"
    },
    "mappings": {
      "properties": {
        "@timestamp": {
          "type": "date"
        },
        "message": {
          "type": "text"
        },
        "log_type": {
          "type": "keyword"
        },
        "service": {
          "type": "keyword"
        },
        "environment": {
          "type": "keyword"
        },
        "timezone": {
          "type": "keyword"
        }
      }
    }
  },
  "data_stream": {}
}
```

应用索引模板：

```bash
curl -X PUT "localhost:9200/_index_template/filebeat" \
  -H "Content-Type: application/json" \
  -u elastic:changeme \
  -d @filebeat-template.json
```

### 5. 创建处理管道（可选：用于时区转换和数据增强）

创建 `pipeline.json` 文件：

```json
{
  "description": "Pipeline for timezone conversion and data enrichment",
  "processors": [
    {
      "date": {
        "field": "timestamp",
        "formats": [
          "yyyy-MM-dd'T'HH:mm:ss.SSSZ",
          "yyyy-MM-dd'T'HH:mm:ssZ",
          "yyyy-MM-dd HH:mm:ss",
          "yyyy/MM/dd HH:mm:ss",
          "UNIX",
          "UNIX_MS"
        ],
        "timezone": "Asia/Shanghai",
        "target_field": "@timestamp",
        "ignore_failure": true
      }
    },
    {
      "set": {
        "field": "event.timezone",
        "value": "Asia/Shanghai",
        "ignore_failure": true
      }
    }
  ]
}
```

应用处理管道：

```bash
curl -X PUT "localhost:9200/_ingest/pipeline/timezone-conversion" \
  -H "Content-Type: application/json" \
  -u elastic:changeme \
  -d @pipeline.json
```

### 6. 配置 filebeat.yml

基于您提供的配置文件，创建 `filebeat.yml`：

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
    # 添加解析错误容忍
    ignore_decoding_error: true
  fields:
    log_type: application
    service: cc-starship
    timezone: Asia/Shanghai
  fields_under_root: true
  tags: ["application", "backend", "api"]
  
  # 简化多行处理配置
  multiline:
    pattern: '^{'
    negate: false
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

# 数据流配置
setup.data_stream.enabled: true
setup.data_stream.namespace: "cc-starship"  # 使用服务名作为命名空间

# 禁用内置模板管理，因为我们手动创建了模板
setup.template.enabled: false

# ILM配置保持不变
setup.ilm.enabled: true
setup.ilm.policy_name: "filebeat"
setup.ilm.rollover_alias: "filebeat" 
setup.ilm.pattern: "{now/d}-000001"
setup.ilm.check_exists: true
setup.ilm.overwrite: true

# 输出配置 - 使用与模板匹配的索引命名模式
output.elasticsearch:
  hosts: ["localhost:9200"]
  username: "elastic"
  password: "changeme"
  # 启用处理管道
  pipeline: "timezone-conversion"
  allow_older_versions: true
  
# 启用监控
monitoring.enabled: true
monitoring.elasticsearch:
  hosts: ["localhost:9200"]
  username: "elastic"
  password: "changeme"

# 增加日志级别，便于调试
logging:
  level: debug
  to_stderr: true
  metrics.enabled: false
  
# 设置更大的缓冲区和批处理
queue.mem:
  events: 4096
  flush.min_events: 512
  flush.timeout: 5s
```

### 7. 启动服务

如果使用 Docker Compose:

```bash
docker-compose up -d
```

如果使用本地安装:

```bash
# 启动 Elasticsearch
./bin/elasticsearch

# 在另一个终端启动 Filebeat
./filebeat -e -c filebeat.yml
```

## 配置详解

### filebeat.yml 配置详解

#### 输入配置

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /app/logs/*.log
    - /usr/share/filebeat/project_logs/*.log
```

此配置指定了 Filebeat 收集日志的路径，支持通配符。

#### JSON 解析

```yaml
json:
  keys_under_root: true
  add_error_key: true
  message_key: message
  overwrite_keys: true
  ignore_decoding_error: true
```

- `keys_under_root`: 将 JSON 解析后的字段放在根级别
- `add_error_key`: 在解析失败时添加错误信息
- `message_key`: 指定 JSON 中哪个字段作为消息内容
- `overwrite_keys`: 允许覆盖已存在的字段
- `ignore_decoding_error`: 忽略解析错误，继续处理

#### 字段和标签

```yaml
fields:
  log_type: application
  service: cc-starship
  timezone: Asia/Shanghai
fields_under_root: true
tags: ["application", "backend", "api"]
```

这些字段和标签使得在 Elasticsearch 中检索和过滤日志更加容易。

#### 多行处理

```yaml
multiline:
  pattern: '^{'
  negate: false
  match: after
  max_lines: 500
  timeout: 5s
```

处理多行日志条目，例如 JSON 格式的日志或异常堆栈。

#### 处理器配置

```yaml
processors:
  - add_host_metadata: ~
  - add_fields:
      target: ''
      fields:
        environment: development
```

处理器允许在发送到 Elasticsearch 之前修改和增强事件。

#### 时间戳处理

```yaml
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
```

配置时间戳解析，支持多种格式，并设置为上海时区。

#### 数据流配置

```yaml
setup.data_stream.enabled: true
setup.data_stream.namespace: "cc-starship"
```

启用数据流功能并设置命名空间。数据流格式为：`logs-{type}-{namespace}`，例如：`logs-application-cc-starship`。

#### ILM 配置

```yaml
setup.ilm.enabled: true
setup.ilm.policy_name: "filebeat"
setup.ilm.rollover_alias: "filebeat" 
setup.ilm.pattern: "{now/d}-000001"
setup.ilm.check_exists: true
setup.ilm.overwrite: true
```

配置索引生命周期管理，自动管理索引从热到冷再到删除的过程。

#### 输出配置

```yaml
output.elasticsearch:
  hosts: ["localhost:9200"]
  username: "elastic"
  password: "changeme"
  pipeline: "timezone-conversion"
  allow_older_versions: true
```

配置 Elasticsearch 输出，指定主机、认证信息和处理管道。

#### 监控和日志

```yaml
monitoring.enabled: true
monitoring.elasticsearch:
  hosts: ["localhost:9200"]
  username: "elastic"
  password: "changeme"

logging:
  level: debug
  to_stderr: true
  metrics.enabled: false
```

启用 Filebeat 自身的监控，并配置日志级别为 debug 以便排查问题。

#### 队列配置

```yaml
queue.mem:
  events: 4096
  flush.min_events: 512
  flush.timeout: 5s
```

配置内存队列大小和刷新策略，优化性能。

## 验证部署

### 1. 检查 Filebeat 状态

```bash
# 容器化部署
docker-compose ps

# 本地部署
filebeat test config -c filebeat.yml
filebeat test output -c filebeat.yml
```

### 2. 检查 Elasticsearch 数据流

```bash
# 列出所有数据流
curl -u elastic:changeme "localhost:9200/_data_stream?pretty"

# 查看特定数据流详情
curl -u elastic:changeme "localhost:9200/_data_stream/logs-*-cc-starship?pretty"
```

### 3. 验证日志数据

```bash
# 搜索日志数据
curl -u elastic:changeme "localhost:9200/logs-*-cc-starship/_search?pretty" -H "Content-Type: application/json" -d'
{
  "size": 10,
  "sort": [
    {
      "@timestamp": {
        "order": "desc"
      }
    }
  ]
}'
```

## 故障排除

### 常见问题及解决方案

#### 1. Filebeat 无法连接到 Elasticsearch

**症状**: Filebeat 日志中出现连接错误

**解决方案**:
- 检查 Elasticsearch 是否运行中: `curl localhost:9200`
- 验证用户名和密码是否正确
- 检查网络连接和防火墙设置

#### 2. 数据未出现在 Elasticsearch 中

**症状**: 没有日志数据出现在 Elasticsearch 中

**解决方案**:
- 检查 Filebeat 是否正确处理日志文件: `filebeat test config -c filebeat.yml`
- 确认日志文件路径是否正确
- 检查 JSON 解析设置是否与日志格式匹配
- 启用 Filebeat 调试日志: `logging.level: debug`

#### 3. 时间戳不正确

**症状**: 日志中的时间戳不符合预期

**解决方案**:
- 检查处理管道配置中的时区设置
- 确认时间戳格式与日志中的实际格式匹配
- 验证 timestamp 处理器配置是否正确

#### 4. 内存使用过高

**症状**: Filebeat 内存使用率高

**解决方案**:
- 优化队列设置
- 增加批处理大小
- 减少同时监控的文件数量

## 性能优化

### 调整参数

根据日志量和硬件资源，可调整以下参数以优化性能：

#### 1. 对于大量日志

```yaml
queue.mem:
  events: 8192
  flush.min_events: 1024
  flush.timeout: 3s

output.elasticsearch:
  bulk_max_size: 2048
  worker: 4
```

#### 2. 对于低延迟要求

```yaml
queue.mem:
  events: 2048
  flush.min_events: 256
  flush.timeout: 1s
```

#### 3. 资源受限环境

```yaml
queue.mem:
  events: 2048
  flush.min_events: 512
  flush.timeout: 10s

output.elasticsearch:
  bulk_max_size: 512
  worker: 2
```

### 硬件建议

- **CPU**: 每扫描 100 个文件建议分配 1 个 CPU 核心
- **内存**: 默认设置下每实例至少 256MB，高负载环境建议 1GB 以上
- **磁盘**: SSD 存储可显著提高性能，特别是对于大量日志文件

## 运维建议

### 监控指标

监控以下关键指标以确保 Filebeat 正常运行：

- **系统负载**: CPU 使用率和内存消耗
- **输入事件率**: 每秒处理的事件数
- **输出事件率**: 每秒发送到 Elasticsearch 的事件数
- **处理延迟**: 事件从生成到索引的时间
- **错误率**: 解析错误和发送失败的数量

您可以使用 Elasticsearch 的监控功能查看这些指标：

```
http://localhost:5601/app/monitoring#/beats
```

### 备份策略

- 定期备份 ILM 策略和索引模板
- 为关键日志配置快照备份
- 使用 Elasticsearch 快照功能进行数据备份

## 参考资料

- [Elastic 官方数据流文档](https://www.elastic.co/guide/en/elasticsearch/reference/8.17/data-streams.html)
- [Filebeat 官方配置指南](https://www.elastic.co/guide/en/beats/filebeat/8.17/configuring-howto-filebeat.html)
- [Elasticsearch ILM 文档](https://www.elastic.co/guide/en/elasticsearch/reference/8.17/index-lifecycle-management.html)
- [Filebeat 处理 JSON 日志](https://www.elastic.co/guide/en/beats/filebeat/8.17/filebeat-input-log.html#filebeat-input-log-config-json)
- [Filebeat 性能调优指南](https://www.elastic.co/guide/en/beats/filebeat/8.17/tune-filebeat-for-high-volumes.html) 