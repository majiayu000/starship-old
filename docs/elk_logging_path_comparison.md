# ELK 日志收集路径选择：直连 vs. Logstash 处理

## 目录

1. [概述](#1-概述)
2. [两种日志收集路径](#2-两种日志收集路径)
3. [方法比较](#3-方法比较)
4. [何时选择直连方式](#4-何时选择直连方式)
5. [何时选择 Logstash 处理](#5-何时选择-logstash-处理)
6. [实施直连方式](#6-实施直连方式)
7. [实施 Logstash 处理方式](#7-实施-logstash-处理方式)
8. [混合策略](#8-混合策略)
9. [迁移路径](#9-迁移路径)
10. [性能考量](#10-性能考量)
11. [参考配置](#11-参考配置)

## 1. 概述

在 ELK 日志收集架构中，有两种主要的日志收集路径：

1. **直连方式**: Filebeat → Elasticsearch
2. **处理方式**: Filebeat → Logstash → Elasticsearch

每种方法都有其优缺点和适用场景。本文档将帮助您了解这些差异，并指导您为项目选择合适的方法。

## 2. 两种日志收集路径

### 2.1 直连方式

```
应用程序 → 日志文件 → Filebeat → Elasticsearch → Kibana
```

在这种方法中，Filebeat 直接将日志发送到 Elasticsearch 进行索引。Filebeat 可以执行一些基本的处理（如 JSON 解析、添加字段等），但处理能力有限。

### 2.2 Logstash 处理方式

```
应用程序 → 日志文件 → Filebeat → Logstash → Elasticsearch → Kibana
```

在这种方法中，Filebeat 将日志发送到 Logstash，Logstash 执行更复杂的处理和转换，然后将处理后的日志发送到 Elasticsearch。

## 3. 方法比较

| 特性 | 直连方式 | Logstash 处理方式 |
|-----|---------|-----------------|
| **架构复杂性** | 简单 | 复杂 |
| **资源消耗** | 低 | 高 |
| **日志处理能力** | 基本 | 高级 |
| **延迟** | 较低 | 较高 |
| **数据转换** | 有限 | 丰富 |
| **缓冲机制** | 有限 | 强大 |
| **多源聚合** | 不支持 | 支持 |
| **容错性** | 基本 | 高级 |
| **可扩展性** | 良好 | 优秀 |
| **维护成本** | 低 | 高 |

## 4. 何时选择直连方式

在以下情况下，Filebeat 直接发送到 Elasticsearch 可能是更好的选择：

- **日志结构已经标准化**: 日志已经是结构良好的 JSON 格式，无需复杂处理
- **资源有限**: 在低规格环境或边缘计算场景中
- **简单场景**: 小型项目或开发/测试环境
- **低延迟要求**: 需要实时查看日志
- **单一日志类型**: 只处理一种格式的日志

## 5. 何时选择 Logstash 处理

在以下情况下，通过 Logstash 处理是更好的选择：

- **复杂数据处理**: 需要复杂的数据转换、过滤或丰富
- **多源日志汇聚**: 整合来自多个不同系统、不同格式的日志
- **高日志量**: 需要强大的缓冲机制来处理日志峰值
- **数据路由**: 需要根据内容将日志发送到不同的目的地
- **生产环境**: 需要高可靠性和可扩展性
- **复杂格式**: 日志格式不一致或需要特殊解析
- **数据丰富化**: 需要添加额外信息（GeoIP、查询外部服务等）

## 6. 实施直连方式

### 6.1 配置 Filebeat

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /path/to/logs/*.log
  json.keys_under_root: true
  json.add_error_key: true
  json.message_key: message
  
  # 基本处理能力
  processors:
    - add_host_metadata: ~
    - add_cloud_metadata: ~
    - add_docker_metadata: ~
    - add_fields:
        fields:
          environment: production
          service: api-server

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  username: "${ES_USERNAME}"
  password: "${ES_PASSWORD}"
  index: "logs-%{+yyyy.MM.dd}"
  
  # 提供一定的缓冲能力
  bulk_max_size: 2048
  worker: 4
```

### 6.2 索引模板

为了确保字段正确映射，应该在 Elasticsearch 中创建索引模板：

```json
PUT _index_template/logs-template
{
  "index_patterns": ["logs-*"],
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 1
    },
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "level": { "type": "keyword" },
        "message": { "type": "text" },
        "service": { "type": "keyword" },
        "environment": { "type": "keyword" }
      }
    }
  }
}
```

## 7. 实施 Logstash 处理方式

### 7.1 配置 Filebeat

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /path/to/logs/*.log
  json.keys_under_root: true
  json.add_error_key: true
  json.message_key: message
  
  # 只做基本处理，主要处理交给 Logstash
  processors:
    - add_host_metadata: ~

# 输出到 Logstash
output.logstash:
  hosts: ["logstash:5044"]
  loadbalance: true
```

### 7.2 配置 Logstash

```
input {
  beats {
    port => 5044
  }
}

filter {
  # 处理 JSON 解析错误
  if [json_parsing_failure] {
    drop { }
  }
  
  # 标准化时间戳
  date {
    match => [ "timestamp", "ISO8601" ]
    target => "@timestamp"
    remove_field => [ "timestamp" ]
  }
  
  # 标准化日志级别
  mutate {
    uppercase => [ "level" ]
  }
  
  # 根据日志类型分流处理
  if [type] == "access" {
    # 解析 API 路径
    grok {
      match => { "[request][path]" => "/api/(?<api_version>v\d+)/(?<resource>[^/]+)(/(?<resource_id>[^/]+))?" }
      tag_on_failure => ["_grokparsefailure_path"]
    }
    
    # 对慢响应标记
    if [response][time_ms] > 500 {
      mutate {
        add_field => { "performance_alert" => "slow_response" }
      }
    }
  }
  
  # 添加地理位置信息
  if [request][ip] {
    geoip {
      source => "[request][ip]"
      target => "[request][geo]"
    }
  }
  
  # 添加环境标签
  mutate {
    add_field => { "environment" => "${ENV:production}" }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    user => "${ES_USERNAME}"
    password => "${ES_PASSWORD}"
    index => "logs-%{+YYYY.MM.dd}"
  }
}
```

## 8. 混合策略

在某些情况下，混合两种方法可能是最佳选择：

- 对关键日志：使用 Logstash 处理路径，确保高可靠性和丰富的处理
- 对常规日志：使用直连方式，节省资源
- 对调试日志：可以考虑直接写入文件，仅在需要时手动导入

配置混合策略的 Filebeat 示例：

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /path/to/logs/critical*.log
  tags: ["critical"]
  json.keys_under_root: true

- type: log
  enabled: true
  paths:
    - /path/to/logs/regular*.log
  tags: ["regular"]
  json.keys_under_root: true

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  username: "${ES_USERNAME}"
  password: "${ES_PASSWORD}"
  indices:
    - index: "regular-logs-%{+yyyy.MM.dd}"
      when.contains:
        tags: "regular"

output.logstash:
  hosts: ["logstash:5044"]
  indices:
    - index: "critical-logs-%{+yyyy.MM.dd}"
      when.contains:
        tags: "critical"
```

## 9. 迁移路径

如果您从直连方式开始，但后来需要 Logstash 的功能，以下是平滑迁移的步骤：

1. **部署 Logstash**：设置 Logstash 实例，配置为仅转发日志（无复杂处理）
2. **更新 Filebeat 配置**：修改输出从 Elasticsearch 到 Logstash
3. **验证流程**：确保日志正常流入
4. **逐步添加 Logstash 处理**：一次添加一个处理规则，并监控影响
5. **根据需要扩展**：随着日志量增长，扩展 Logstash 节点

## 10. 性能考量

### 10.1 直连方式性能优化

- 增加 Filebeat 工作线程数量
- 优化批处理大小
- 确保足够的 Elasticsearch 节点处理索引请求

```yaml
output.elasticsearch:
  worker: 4
  bulk_max_size: 2048
  flush_interval: "1s"
```

### 10.2 Logstash 方式性能优化

- 调整 Logstash 管道工作线程数
- 优化批处理大小和延迟
- 使用持久队列处理流量峰值

```
# logstash.yml
pipeline.workers: 4
pipeline.batch.size: 1000
pipeline.batch.delay: 50

# 持久队列配置
queue.type: persisted
queue.max_bytes: 4gb
```

## 11. 参考配置

### 11.1 完整的直连方式配置

参见[6.1 配置 Filebeat](#61-配置-filebeat)部分。

### 11.2 完整的 Logstash 处理方式配置

参见[7.1 配置 Filebeat](#71-配置-filebeat)和[7.2 配置 Logstash](#72-配置-logstash)部分。

---

*无论您选择哪种方法，关键是确保它符合您的日志处理需求、可用资源和团队能力。最佳实践是从简单开始，随着需求的增长逐步扩展复杂性。* 