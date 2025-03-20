# ELK 日志管理系统设计方案

## 目录

- [1. 整体架构设计](#1-整体架构设计)
- [2. 日志收集策略](#2-日志收集策略)
- [3. Logstash 配置和处理管道](#3-logstash-配置和处理管道)
- [4. Elasticsearch 索引设计](#4-elasticsearch-索引设计)
- [5. Kibana 仪表盘设计](#5-kibana-仪表盘设计)
- [6. 部署方式](#6-部署方式)
- [7. 高级功能](#7-高级功能)
- [8. 最佳实践](#8-最佳实践)
- [9. 实施路线图](#9-实施路线图)
- [10. 与当前日志系统的集成](#10-与当前日志系统的集成)

## 1. 整体架构设计

基于ELK栈的完整日志管理解决方案架构如下：

```
应用服务器 ──► Filebeat ──┐
                          │
微服务实例 ──► Filebeat ──┼──► Logstash ──► Elasticsearch ◄──── Kibana
                          │                      │
其他日志源 ──► Filebeat ──┘                      │
                                                │
告警系统 ◄──── Elastalert ◄─────────────────────┘
```

### 组件说明:

1. **Filebeat**: 轻量级日志收集器，从各个服务器收集日志文件
2. **Logstash**: 日志处理和转换引擎，解析和结构化日志数据
3. **Elasticsearch**: 分布式搜索和分析引擎，存储和索引日志数据
4. **Kibana**: 可视化平台，提供搜索、查看和分析日志的界面
5. **Elastalert** (可选): 基于Elasticsearch的告警系统

## 2. 日志收集策略

### Filebeat 配置:

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /path/to/logs/app.log
  json.keys_under_root: true
  json.add_error_key: true
  json.message_key: message

output.logstash:
  hosts: ["logstash:5044"]
```

这个配置会:
- 监控日志文件
- 自动解析JSON格式
- 将JSON字段提升到顶层
- 将日志发送给Logstash进一步处理

## 3. Logstash 配置和处理管道

以下是适合当前日志格式的Logstash配置:

```
input {
  beats {
    port => 5044
  }
}

filter {
  if ![timestamp] {
    date {
      match => [ "@timestamp", "ISO8601" ]
      target => "timestamp"
    }
  } else {
    date {
      match => [ "timestamp", "ISO8601" ]
      target => "timestamp"
    }
  }
  
  # 标准化日志级别
  mutate {
    uppercase => [ "level" ]
  }
  
  # 添加环境标签
  mutate {
    add_field => { "environment" => "${ENVIRONMENT:development}" }
  }
  
  # 根据日志类型进行特定处理
  if [type] == "access" {
    # 处理访问日志
    mutate {
      add_field => { "request_full_path" => "%{[request][path]}?%{[request][query]}" }
    }
    
    # 提取 HTTP 状态码类别
    ruby {
      code => "
        status = event.get('[response][status]')
        if status
          event.set('status_category', status / 100)
        end
      "
    }
  }
  
  # 处理业务日志
  if [type] == "business" {
    # 可以添加特定于业务日志的增强
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "logs-%{+YYYY.MM.dd}"
    # 使用日志类型作为doctype (可选，在新版Elasticsearch中已弃用)
    document_type => "%{[type]}"
  }
}
```

## 4. Elasticsearch 索引设计

### 索引模板:

```json
{
  "index_patterns": ["logs-*"],
  "settings": {
    "number_of_shards": 3,
    "number_of_replicas": 1,
    "refresh_interval": "10s",
    "index.lifecycle.name": "logs_policy"
  },
  "mappings": {
    "properties": {
      "timestamp": { "type": "date" },
      "log_id": { "type": "keyword" },
      "type": { "type": "keyword" },
      "level": { "type": "keyword" },
      "service": { "type": "keyword" },
      "message": { "type": "text" },
      "trace_id": { "type": "keyword" },
      "span_id": { "type": "keyword" },
      "request": {
        "properties": {
          "method": { "type": "keyword" },
          "path": { "type": "keyword" },
          "query": { "type": "text" },
          "client_ip": { "type": "ip" },
          "user_agent": { "type": "text" }
        }
      },
      "response": {
        "properties": {
          "status": { "type": "integer" },
          "status_category": { "type": "byte" },
          "time_ms": { "type": "float" },
          "size": { "type": "long" }
        }
      },
      "details": { "type": "object", "dynamic": true },
      "context": { "type": "object", "dynamic": true }
    }
  }
}
```

### 索引生命周期策略:

```json
{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_age": "1d",
            "max_size": "50gb"
          }
        }
      },
      "warm": {
        "min_age": "3d",
        "actions": {
          "forcemerge": {
            "max_num_segments": 1
          },
          "shrink": {
            "number_of_shards": 1
          }
        }
      },
      "cold": {
        "min_age": "30d",
        "actions": {
          "allocate": {
            "require": {
              "data": "cold"
            }
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

## 5. Kibana 仪表盘设计

### 推荐的仪表盘:

1. **总览仪表盘**
   - 日志数量随时间变化趋势
   - 按日志级别分布的饼图
   - 按日志类型分布的饼图
   - 服务健康状况指标

2. **访问日志仪表板**
   - HTTP状态码分布
   - 响应时间分布
   - 最频繁访问的路径
   - 客户端IP地域分布
   - 错误率随时间变化

3. **业务日志仪表板**
   - 按操作类型分组的业务操作
   - 用户活动热图
   - 关键业务指标监控
   - 业务异常监控

4. **错误和警告仪表板**
   - 错误日志时间线
   - 按服务分组的错误计数
   - 常见错误类型分析
   - 错误追踪显示

### 保存的搜索:

为常见的日志查询创建保存的搜索:

- 错误日志: `level:ERROR`
- 访问日志: `type:access`
- 业务日志: `type:business`
- 慢响应: `response.time_ms:>500`
- 特定服务日志: `service:"cc-starship"`

## 6. 部署方式

### Docker Compose (开发环境):

```yaml
version: '3'
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.7.0
    environment:
      - discovery.type=single-node
      - ES_JAVA_OPTS=-Xms1g -Xmx1g
    ports:
      - "9200:9200"
    volumes:
      - es_data:/usr/share/elasticsearch/data

  logstash:
    image: docker.elastic.co/logstash/logstash:8.7.0
    depends_on:
      - elasticsearch
    ports:
      - "5044:5044"
    volumes:
      - ./logstash/pipeline:/usr/share/logstash/pipeline
      - ./logstash/config/logstash.yml:/usr/share/logstash/config/logstash.yml

  kibana:
    image: docker.elastic.co/kibana/kibana:8.7.0
    depends_on:
      - elasticsearch
    ports:
      - "5601:5601"
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200

  filebeat:
    image: docker.elastic.co/beats/filebeat:8.7.0
    user: root
    volumes:
      - ./filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - /var/log:/var/log:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    depends_on:
      - logstash

volumes:
  es_data:
```

### Kubernetes (生产环境):

推荐使用Elastic Cloud Kubernetes (ECK) 操作员进行部署，或使用Helm图表。

### 硬件需求 (生产环境):

- **Elasticsearch**: 至少3个节点，每个节点8GB内存、4核CPU、100GB存储
- **Logstash**: 至少2个节点，每个节点4GB内存、2核CPU
- **Kibana**: 1-2个节点，每个节点2GB内存、2核CPU
- **Filebeat**: 在每个日志产生的服务器上部署

## 7. 高级功能

### 日志审计和安全:

- 使用X-Pack Security启用身份验证和授权
- 对敏感信息启用字段级加密
- 配置日志的不可变性和审计追踪

### 告警配置:

使用Elastalert或内置的Elasticsearch Alerting设置告警:

```yaml
# Elastalert配置示例
name: Error Rate Alert
type: frequency
index: logs-*
num_events: 10
timeframe:
  minutes: 5
filter:
- term:
    level: "ERROR"
alert:
- "email"
email:
- "admin@example.com"
```

### APM 集成:

通过配置Elastic APM，将追踪数据与日志相关联，实现端到端可观测性:

```yaml
# APM配置
apm-server:
  host: apm:8200
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
```

## 8. 最佳实践

1. **日志保留策略**
   - 设置合理的日志保留期限（热数据、暖数据、冷数据和归档）
   - 使用ILM管理索引生命周期

2. **性能优化**
   - 监控和调整Elasticsearch堆大小
   - 优化查询和聚合
   - 定期执行索引优化
   - 对于高流量应用，考虑使用Logstash过滤器来减少不必要的字段

3. **故障恢复**
   - 配置Elasticsearch快照和恢复
   - 设置数据备份策略
   - 实现高可用性配置

4. **团队使用建议**
   - 创建团队特定的仪表板和访问控制
   - 举办ELK培训会议
   - 整理常用查询的知识库
   - 建立日志记录标准，确保日志内容一致

## 9. 实施路线图

1. **阶段1: 基础设置 (1-2周)**
   - 部署ELK栈
   - 配置Filebeat和基本日志收集
   - 创建基本索引和映射

2. **阶段2: 日志增强和仪表板 (2-3周)**
   - 实施高级Logstash过滤器
   - 创建初始Kibana仪表板
   - 设置基本告警

3. **阶段3: 优化和扩展 (2-4周)**
   - 优化性能和存储
   - 添加高级分析和可视化
   - 与现有监控系统集成

4. **阶段4: 培训和文档 (1-2周)**
   - 开发用户文档
   - 团队培训
   - 收集反馈并进行调整

## 10. 与当前日志系统的集成

以下是特定针对当前JSON日志格式的集成点:

1. **直接兼容性**
   - 当前JSON日志格式天然适合Elasticsearch
   - 不需要重大的解析转换，字段大部分可直接使用

2. **特定字段映射**
   - `log_id` → 用于唯一标识和追踪
   - `trace_id` & `span_id` → 用于分布式跟踪关联
   - `type` (access/business/system) → 用于日志类型过滤

3. **访问日志特殊处理**
   - 创建对`response.status`、`request.path`和`request.method`的专用可视化
   - 分析慢响应和错误模式

## 附录A: 常用Kibana查询

### KQL (Kibana查询语言) 示例:

```
# 查找所有错误日志
level:ERROR

# 查找特定服务的错误
level:ERROR and service:"cc-starship"

# 查找响应时间超过500ms的请求
type:access and response.time_ms > 500

# 查找特定时间范围内的日志
@timestamp >= "2023-03-01" and @timestamp <= "2023-03-31"

# 查找特定路径的错误请求
type:access and request.path:/api/v1/users and response.status >= 400

# 查找包含特定跟踪ID的所有日志
trace_id:"eb4f5d7c-19ee-4389-b75c-b1ab85e073b8"
```

## 附录B: ELK资源链接

- [Elasticsearch 官方文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [Logstash 官方文档](https://www.elastic.co/guide/en/logstash/current/index.html)
- [Kibana 官方文档](https://www.elastic.co/guide/en/kibana/current/index.html)
- [Filebeat 官方文档](https://www.elastic.co/guide/en/beats/filebeat/current/index.html)
- [Elastic Stack 最佳实践](https://www.elastic.co/guide/en/elasticsearch/reference/current/high-availability.html) 