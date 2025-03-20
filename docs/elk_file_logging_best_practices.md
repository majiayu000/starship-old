# ELK 日志收集最佳实践：基于文件的日志收集

## 目录

1. [介绍](#1-介绍)
2. [架构设计](#2-架构设计)
3. [应用程序配置](#3-应用程序配置)
4. [Filebeat 配置](#4-filebeat-配置)
5. [Logstash 处理](#5-logstash-处理)
6. [Elasticsearch 索引设计](#6-elasticsearch-索引设计)
7. [Kibana 可视化](#7-kibana-可视化)
8. [部署策略](#8-部署策略)
9. [故障排除](#9-故障排除)
10. [性能优化](#10-性能优化)
11. [安全性考量](#11-安全性考量)
12. [附录：配置示例](#12-附录配置示例)

## 1. 介绍

ELK 栈（Elasticsearch、Logstash、Kibana）是一套强大的开源日志管理解决方案。在大多数生产环境中，基于文件的日志收集是 ELK 实施的首选方法，它提供了以下优势：

- **解耦合**：应用只负责生成日志，不关心传输和存储
- **可靠性**：即使 ELK 组件临时不可用，日志仍然安全地写入文件
- **标准化**：符合"日志即流"的云原生设计理念
- **灵活性**：可以在不修改应用代码的情况下调整日志收集策略

### 1.1 日志收集策略比较

| 策略 | 优点 | 缺点 |
|------|------|------|
| **文件日志** | 解耦合、可靠、容错 | 需要额外的日志文件管理 |
| **直接 API 注入** | 实时性高、简单 | 高耦合、单点故障风险 |
| **标准输出** | 容器友好、简单 | 日志归档困难、不适合全部日志类型 |

## 2. 架构设计

基于文件的 ELK 日志收集架构包含以下组件：

```
应用程序 → 日志文件 → Filebeat → [Logstash] → Elasticsearch → Kibana
```

- **应用程序**：生成结构化（JSON）日志文件
- **Filebeat**：轻量级日志收集器，监控和传输日志文件
- **Logstash**（可选）：进行高级日志处理、转换和丰富化
- **Elasticsearch**：存储和索引日志数据
- **Kibana**：可视化和分析日志数据

### 2.1 组件职责

| 组件 | 主要职责 |
|------|----------|
| **应用程序** | 生成格式一致的日志文件，处理日志轮转 |
| **Filebeat** | 可靠地读取日志文件，记住读取位置，处理文件轮转 |
| **Logstash** | 解析、过滤、转换和丰富日志数据 |
| **Elasticsearch** | 存储、索引和搜索日志数据 |
| **Kibana** | 提供日志数据的可视化和查询界面 |

## 3. 应用程序配置

### 3.1 日志输出格式

对于与 ELK 的最佳集成，应用应该输出结构化的 JSON 日志：

```json
{
  "timestamp": "2023-11-19T08:45:12.345Z",
  "level": "INFO",
  "type": "access",
  "service": "api-server",
  "trace_id": "abc123",
  "message": "请求处理完成",
  "request": {
    "method": "GET",
    "path": "/api/v1/users",
    "ip": "192.168.1.1"
  },
  "response": {
    "status": 200,
    "time_ms": 45.2
  }
}
```

### 3.2 Go 应用程序配置示例

```yaml
# app.yaml 配置
logger:
  level: "info"
  format: "json"
  processors:
    console:
      enabled: true  # 开发环境保留控制台输出
    file:
      enabled: true  # 启用文件日志
      path: "logs/app.log"
      rotation:
        max_size: "100MB"
        max_age: "7d"
        max_backups: 10
    elk:
      enabled: false  # 禁用直接 ELK 发送
```

### 3.3 日志文件存储最佳实践

- **统一路径**：所有日志统一存放在一个目录
- **服务隔离**：每个服务使用单独的日志文件
- **命名规范**：使用 `服务名-日志类型.log` 的命名模式
- **权限设置**：确保日志目录有适当的读写权限
- **日志轮转**：配置合适的日志轮转策略，防止单个文件过大

## 4. Filebeat 配置

Filebeat 是 Elastic 开发的轻量级日志收集器，专为可靠地收集文件日志而设计。

### 4.1 基本配置示例

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /path/to/logs/*.log
  json.keys_under_root: true
  json.add_error_key: true
  json.message_key: message
  fields:
    environment: production
    application: cc-starship
  fields_under_root: true

output.logstash:
  hosts: ["logstash:5044"]

# 或直接输出到 Elasticsearch
# output.elasticsearch:
#   hosts: ["elasticsearch:9200"]
#   username: "elastic"
#   password: "${ELASTIC_PASSWORD}"
#   index: "logs-%{+yyyy.MM.dd}"
```

### 4.2 多日志类型处理

针对不同类型的日志文件配置不同的处理规则：

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /path/to/logs/access*.log
  tags: ["access"]
  json.keys_under_root: true
  
- type: log
  enabled: true
  paths:
    - /path/to/logs/error*.log
  tags: ["error"]
  json.keys_under_root: true
  
- type: log
  enabled: true
  paths:
    - /path/to/logs/business*.log
  tags: ["business"]
  json.keys_under_root: true
```

### 4.3 增强功能配置

```yaml
# 自动重载配置
filebeat.config.inputs:
  enabled: true
  path: ${path.config}/conf.d/*.yml
  reload.enabled: true
  reload.period: 10s

# 处理器配置（预处理）
processors:
  - add_host_metadata: ~
  - add_cloud_metadata: ~
  - add_docker_metadata: ~

# 日志记录
logging.level: info
logging.to_files: true
logging.files:
  path: /var/log/filebeat
  name: filebeat
  keepfiles: 7
  permissions: 0644
```

## 5. Logstash 处理

虽然 Filebeat 可以直接发送到 Elasticsearch，但 Logstash 提供了更强大的日志处理能力。

### 5.1 基本 Logstash 配置

```
input {
  beats {
    port => 5044
  }
}

filter {
  if [json_parsing_failure] {
    drop { }
  }
  
  date {
    match => [ "timestamp", "ISO8601" ]
    target => "@timestamp"
    remove_field => [ "timestamp" ]
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    user => "${ELASTIC_USERNAME}"
    password => "${ELASTIC_PASSWORD}"
    index => "%{[service]}-%{+YYYY.MM.dd}"
  }
}
```

### 5.2 高级日志处理

```
filter {
  # 标准化日志级别
  mutate {
    uppercase => [ "level" ]
  }
  
  # 分类处理不同类型的日志
  if [type] == "access" {
    # 解析请求路径信息
    grok {
      match => { "[request][path]" => "/api/(?<api_version>v\d+)/(?<resource>[^/]+)(/(?<resource_id>[^/]+))?" }
      tag_on_failure => ["_grokparsefailure_path"]
      break_on_match => false
    }
    
    # 设置响应时间阈值
    if [response][time_ms] > 500 {
      mutate {
        add_field => { "performance_alert" => "slow_response" }
      }
    }
  }
  
  if [type] == "business" {
    # 业务日志特殊处理
    # ...
  }
  
  # 添加地理位置信息（如果有IP）
  if [request][ip] {
    geoip {
      source => "[request][ip]"
      target => "[request][geo]"
    }
  }
  
  # 丰富跟踪信息
  if [trace_id] {
    mutate {
      add_field => { "has_trace" => "true" }
    }
  }
}
```

## 6. Elasticsearch 索引设计

为了高效存储和检索日志，合理设计 Elasticsearch 索引是关键。

### 6.1 索引模板

```json
PUT _index_template/logs-template
{
  "index_patterns": ["*-*"],
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 1,
      "refresh_interval": "5s",
      "index.lifecycle.name": "logs-policy"
    },
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "service": { "type": "keyword" },
        "level": { "type": "keyword" },
        "type": { "type": "keyword" },
        "message": { "type": "text" },
        "trace_id": { "type": "keyword" },
        "user_id": { "type": "keyword" },
        "request": {
          "properties": {
            "method": { "type": "keyword" },
            "path": { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
            "ip": { "type": "ip" },
            "user_agent": { "type": "text" }
          }
        },
        "response": {
          "properties": {
            "status": { "type": "integer" },
            "time_ms": { "type": "float" }
          }
        }
      }
    }
  }
}
```

### 6.2 索引生命周期管理

```json
PUT _ilm/policy/logs-policy
{
  "policy": {
    "phases": {
      "hot": {
        "actions": {
          "rollover": {
            "max_size": "50GB",
            "max_age": "1d"
          },
          "set_priority": {
            "priority": 100
          }
        }
      },
      "warm": {
        "min_age": "2d",
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
        "min_age": "7d",
        "actions": {
          "set_priority": {
            "priority": 0
          }
        }
      },
      "delete": {
        "min_age": "30d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}
```

## 7. Kibana 可视化

Kibana 提供了强大的日志可视化功能。以下是一些常用的仪表板配置。

### 7.1 索引模式创建

1. 导航至 **Stack Management > Kibana > Index Patterns**
2. 创建索引模式 `*-*`
3. 选择 `@timestamp` 作为时间字段

### 7.2 推荐的仪表板视图

#### 7.2.1 概览仪表板
- 按日志级别分布的饼图
- 按服务分布的饼图
- 随时间变化的日志数量曲线图
- 最近错误日志列表

#### 7.2.2 访问日志仪表板
- HTTP 状态码分布
- 平均响应时间趋势
- 最慢的 API 端点
- 按客户端 IP 的请求数量
- 请求量最大的资源

#### 7.2.3 业务日志仪表板
- 按操作类型分类的业务活动
- 按用户 ID 分组的活动数量
- 业务错误趋势

### 7.3 保存的搜索

创建常用搜索查询：

- 错误日志: `level:ERROR`
- 访问日志: `type:access`
- 业务日志: `type:business`
- 慢响应: `type:access AND response.time_ms:>500`
- 特定用户活动: `user_id:"123456"`

## 8. 部署策略

### 8.1 单机环境

适用于开发和小规模测试环境：

```bash
# 启动 Filebeat
filebeat -e -c filebeat.yml

# 或使用 systemd 服务
sudo systemctl enable filebeat
sudo systemctl start filebeat
```

### 8.2 Docker 环境

```yaml
# docker-compose.yml
version: '3'
services:
  app:
    image: your-app-image
    volumes:
      - ./logs:/app/logs

  filebeat:
    image: docker.elastic.co/beats/filebeat:8.7.0
    volumes:
      - ./filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - ./logs:/logs:ro
    command: filebeat -e -strict.perms=false
    depends_on:
      - elasticsearch

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.7.0
    environment:
      - discovery.type=single-node
      - ES_JAVA_OPTS=-Xms1g -Xmx1g
      - xpack.security.enabled=false
    ports:
      - "9200:9200"
    volumes:
      - es_data:/usr/share/elasticsearch/data

  kibana:
    image: docker.elastic.co/kibana/kibana:8.7.0
    ports:
      - "5601:5601"
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    depends_on:
      - elasticsearch

volumes:
  es_data:
```

### 8.3 Kubernetes 环境

推荐使用 Filebeat DaemonSet 模式：

```yaml
# filebeat-kubernetes.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: filebeat
  namespace: logging
spec:
  selector:
    matchLabels:
      app: filebeat
  template:
    metadata:
      labels:
        app: filebeat
    spec:
      serviceAccountName: filebeat
      containers:
      - name: filebeat
        image: docker.elastic.co/beats/filebeat:8.7.0
        args: ["-c", "/etc/filebeat.yml", "-e"]
        volumeMounts:
        - name: config
          mountPath: /etc/filebeat.yml
          subPath: filebeat.yml
        - name: logs
          mountPath: /var/log
          readOnly: true
      volumes:
      - name: config
        configMap:
          name: filebeat-config
      - name: logs
        hostPath:
          path: /var/log
```

## 9. 故障排除

### 9.1 常见问题与解决方案

| 问题 | 可能原因 | 解决方案 |
|------|----------|----------|
| Filebeat 不收集日志 | 文件权限问题 | 检查 Filebeat 用户对日志文件的读取权限 |
| | 路径配置错误 | 验证日志文件路径是否正确 |
| | 文件格式问题 | 检查日志文件是否符合预期格式 |
| JSON 解析错误 | 非法 JSON 格式 | 使用 jq 工具验证日志文件格式 |
| | 字段类型不匹配 | 在 Elasticsearch 映射中调整字段类型 |
| 日志丢失 | Logstash 缓冲区溢出 | 增加 Logstash 队列大小或内存 |
| | 网络连接问题 | 检查网络连接和防火墙规则 |

### 9.2 诊断工具

```bash
# 验证 Filebeat 配置
filebeat test config -c filebeat.yml

# 测试 Filebeat 输出
filebeat test output -c filebeat.yml

# 验证 JSON 格式
cat /path/to/logfile.log | jq

# 检查 Elasticsearch 索引
curl -X GET "localhost:9200/_cat/indices?v"

# 查看 Filebeat 状态
filebeat status
```

### 9.3 监控指标

- **Filebeat**: 监控收集的字节数、文件数和发送失败次数
- **Logstash**: 监控处理事件数量、处理时间和队列大小
- **Elasticsearch**: 监控索引大小、查询延迟和垃圾收集时间

## 10. 性能优化

### 10.1 Filebeat 优化

```yaml
# filebeat.yml 性能优化
filebeat.inputs:
  - type: log
    # ...
    harvester_buffer_size: 16384
    max_bytes: 10485760

filebeat.registry.flush: 5s

queue.mem:
  events: 4096
  flush.min_events: 512
  flush.timeout: 5s

output.logstash:
  hosts: ["logstash:5044"]
  worker: 4
  bulk_max_size: 2048
```

### 10.2 Logstash 优化

```
# logstash.yml
pipeline.workers: 4
pipeline.batch.size: 1000
pipeline.batch.delay: 50
```

### 10.3 Elasticsearch 优化

- 使用日期索引，定期轮转
- 实施索引生命周期管理
- 根据日志保留需求定期删除旧索引
- 使用适当的分片数量和副本数量

## 11. 安全性考量

### 11.1 加密传输

配置 SSL/TLS 进行安全传输：

```yaml
# filebeat.yml
output.logstash:
  hosts: ["logstash:5044"]
  ssl.enabled: true
  ssl.certificate_authorities: ["path/to/ca.pem"]
  ssl.certificate: "path/to/cert.pem"
  ssl.key: "path/to/key.pem"
```

### 11.2 认证配置

```yaml
# filebeat.yml
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  username: "${ELASTICSEARCH_USERNAME}"
  password: "${ELASTICSEARCH_PASSWORD}"
```

### 11.3 敏感数据处理

```yaml
# filebeat.yml
processors:
  - drop_fields:
      fields: ["credit_card", "password"]
  
  - rename:
      fields:
        - from: "user.email"
          to: "user.email_address"
          
  - dissect:
      tokenizer: "%{ip} %{port}"
      field: "host"
      target_prefix: "network"
```

## 12. 附录：配置示例

### 12.1 完整 Filebeat 配置

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /path/to/logs/*.log
  json.keys_under_root: true
  json.add_error_key: true
  json.message_key: message
  fields:
    environment: production
    application: cc-starship
  fields_under_root: true
  
processors:
  - add_host_metadata: ~
  - add_cloud_metadata: ~
  - add_docker_metadata: ~
  - drop_fields:
      fields: ["sensitive_data", "password"]

setup.ilm.enabled: true
setup.ilm.rollover_alias: "logs"
setup.ilm.pattern: "{now/d}-000001"
setup.ilm.policy_name: "logs-policy"

output.logstash:
  hosts: ["logstash:5044"]
  loadbalance: true
  worker: 2
  bulk_max_size: 2048

logging.level: info
logging.to_files: true
logging.files:
  path: /var/log/filebeat
  name: filebeat
  keepfiles: 7
  permissions: 0644
```

### 12.2 完整 Logstash 配置

```
input {
  beats {
    port => 5044
    ssl => true
    ssl_certificate => "/etc/pki/tls/certs/logstash.crt"
    ssl_key => "/etc/pki/tls/private/logstash.key"
  }
}

filter {
  if [json_parsing_failure] {
    drop { }
  }
  
  date {
    match => [ "timestamp", "ISO8601" ]
    target => "@timestamp"
    remove_field => [ "timestamp" ]
  }
  
  # 标准化日志级别
  mutate {
    uppercase => [ "level" ]
  }
  
  # 分类处理
  if [type] == "access" {
    grok {
      match => { "[request][path]" => "/api/(?<api_version>v\d+)/(?<resource>[^/]+)(/(?<resource_id>[^/]+))?" }
      tag_on_failure => ["_grokparsefailure_path"]
      break_on_match => false
    }
    
    if [response][time_ms] > 500 {
      mutate {
        add_field => { "performance_alert" => "slow_response" }
      }
    }
  }
  
  if [type] == "business" {
    # 业务日志特殊处理
  }
  
  # 添加环境标签
  mutate {
    add_field => { "environment" => "${ENV:production}" }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    user => "${ELASTIC_USERNAME}"
    password => "${ELASTIC_PASSWORD}"
    index => "%{[service]}-%{+YYYY.MM.dd}"
    pipeline => "logs-pipeline"
  }
}
```

### 12.3 索引生命周期策略

```json
PUT _ilm/policy/logs-policy
{
  "policy": {
    "phases": {
      "hot": {
        "actions": {
          "rollover": {
            "max_size": "50GB",
            "max_age": "1d"
          }
        }
      },
      "warm": {
        "min_age": "2d",
        "actions": {
          "shrink": {
            "number_of_shards": 1
          },
          "forcemerge": {
            "max_num_segments": 1
          }
        }
      },
      "cold": {
        "min_age": "7d",
        "actions": {
          "freeze": {}
        }
      },
      "delete": {
        "min_age": "30d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}
```

---

*本文档提供了关于使用 ELK 栈收集文件日志的最佳实践指南。根据您的具体环境和需求，可能需要调整配置。* 