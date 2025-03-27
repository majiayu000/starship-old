# Filebeat 配置与应用指南

## 目录

1. [介绍](#1-介绍)
2. [安装与部署](#2-安装与部署)
3. [基本配置](#3-基本配置)
4. [输入配置详解](#4-输入配置详解)
5. [处理器配置](#5-处理器配置)
6. [输出配置](#6-输出配置)
7. [模块系统](#7-模块系统)
8. [多日志源管理](#8-多日志源管理)
9. [性能优化](#9-性能优化)
10. [监控与管理](#10-监控与管理)
11. [安全性配置](#11-安全性配置)
12. [常见问题排查](#12-常见问题排查)
13. [配置示例集](#13-配置示例集)

## 1. 介绍

Filebeat 是 Elastic Stack 的轻量级日志收集器，专为高效、可靠地收集和发送日志数据而设计。它可以监控指定的日志文件或位置，收集日志事件，并将它们转发到 Elasticsearch 或 Logstash 进行索引。

### 1.1 Filebeat 的主要优势

- **轻量级**：占用资源少，对系统影响小
- **可靠性**：内置的数据发送重试机制和记录功能
- **零维护**：一旦设置好，无需额外维护
- **多源支持**：可以从多种来源收集日志
- **自动发现**：能够自动发现日志文件
- **内置解析器**：支持多种常见日志格式的解析

### 1.2 Filebeat 的工作原理

Filebeat 使用 harvester（收割机）和 prospector（勘察器）的概念：

- **Harvester**：负责读取单个文件的内容，它逐行读取每个文件，并将内容发送到输出
- **Prospector**：负责管理 harvester，并查找所有匹配的文件源

收集过程包括：

1. 定位日志文件并打开
2. 逐行读取文件内容
3. 将内容发送到配置的输出
4. 记住上次读取的位置，以便在下次启动时继续

## 2. 安装与部署

### 2.1 在不同平台上安装 Filebeat

#### MacOS (Homebrew)

```bash
brew tap elastic/tap
brew install elastic/tap/filebeat
```

#### Debian/Ubuntu

```bash
# 添加 Elastic 的 GPG 密钥
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo apt-key add -

# 添加 Elastic 的 APT 仓库
echo "deb https://artifacts.elastic.co/packages/8.x/apt stable main" | sudo tee -a /etc/apt/sources.list.d/elastic-8.x.list

# 更新包信息并安装
sudo apt-get update && sudo apt-get install filebeat
```

#### CentOS/RHEL

```bash
# 添加 Elastic 的 YUM 仓库
sudo rpm --import https://artifacts.elastic.co/GPG-KEY-elasticsearch
sudo tee /etc/yum.repos.d/elastic.repo > /dev/null <<EOT
[elastic-8.x]
name=Elastic repository for 8.x packages
baseurl=https://artifacts.elastic.co/packages/8.x/yum
gpgcheck=1
gpgkey=https://artifacts.elastic.co/GPG-KEY-elasticsearch
enabled=1
autorefresh=1
type=rpm-md
EOT

# 安装
sudo yum install filebeat
```

#### Docker

```bash
docker pull docker.elastic.co/beats/filebeat:8.17.3
```

### 2.2 目录结构

安装后的 Filebeat 主要目录结构：

```
/etc/filebeat/           # 配置文件目录
  filebeat.yml           # 主配置文件
  modules.d/             # 模块配置文件
/var/lib/filebeat/       # 数据目录，存储注册表等
/var/log/filebeat/       # 日志目录
```

### 2.3 基本管理命令

```bash
# 启动 Filebeat
sudo systemctl start filebeat    # 使用 systemd
sudo service filebeat start      # 使用 init.d

# 停止 Filebeat
sudo systemctl stop filebeat
sudo service filebeat stop

# 查看状态
sudo systemctl status filebeat
sudo service filebeat status

# 启用开机自启
sudo systemctl enable filebeat
sudo chkconfig filebeat on

# 检查配置是否有效
filebeat test config -c /etc/filebeat/filebeat.yml

# 测试输出连接
filebeat test output -c /etc/filebeat/filebeat.yml
```

### 2.4 Docker 部署示例

```bash
docker run -d \
  --name=filebeat \
  --user=root \
  --volume="$(pwd)/filebeat.yml:/usr/share/filebeat/filebeat.yml:ro" \
  --volume="/var/log:/var/log:ro" \
  --volume="/var/lib/docker/containers:/var/lib/docker/containers:ro" \
  docker.elastic.co/beats/filebeat:8.17.3 \
  filebeat -e -strict.perms=false
```

## 3. 基本配置

### 3.1 配置文件结构

Filebeat 使用 YAML 格式的配置文件。主配置文件 `filebeat.yml` 有以下主要部分：

```yaml
filebeat.inputs:        # 输入配置
filebeat.config:        # 外部配置
filebeat.modules:       # 模块配置
processors:             # 处理器配置
output:                 # 输出配置
logging:                # 日志配置
```

### 3.2 最小配置示例

最简单的可运行配置示例：

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/*.log

output.elasticsearch:
  hosts: ["localhost:9200"]
```

### 3.3 配置验证

在应用新配置前，使用以下命令验证配置是否有效：

```bash
filebeat test config -c /etc/filebeat/filebeat.yml
```

### 3.4 变量和环境变量

Filebeat 配置中可以使用两种变量：

- **环境变量**：使用 `${ENV_VAR}` 语法
- **内部变量**：如 `%{[field.name]}`、`${hostname}`

```yaml
# 使用环境变量
output.elasticsearch:
  hosts: ["${ES_HOST:localhost}:${ES_PORT:9200}"]
  username: "${ES_USER}"
  password: "${ES_PASSWORD}"

# 使用内部变量
output.elasticsearch:
  index: "logs-%{[agent.version]}-%{+yyyy.MM.dd}"
```

## 4. 输入配置详解

Filebeat 支持多种输入类型，最常用的是 `log` 类型。

### 4.1 基本日志输入配置

```yaml
filebeat.inputs:
- type: log                # 输入类型
  enabled: true           # 是否启用
  paths:                  # 文件路径，支持通配符
    - /var/log/*.log
    - /var/log/app/*.log
  exclude_files: ['\.gz$'] # 排除的文件
  ignore_older: 24h       # 忽略 24 小时前的文件
  scan_frequency: 10s     # 扫描新文件的频率
  tail_files: false       # 是否从文件末尾开始收集
  symlinks: true          # 是否跟踪符号链接
  harvester_buffer_size: 16384  # 收割机缓冲区大小
  max_bytes: 10485760     # 单个消息最大字节数
```

### 4.2 多行处理

对于像 Java 堆栈跟踪这样的多行日志，Filebeat 提供了多行处理功能：

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/app.log
  multiline:
    pattern: '^[0-9]{4}-[0-9]{2}-[0-9]{2}'  # 匹配日期开头的行
    negate: true                           # 如果不匹配，则与上一行组合
    match: after                           # 不匹配的行添加到上一个匹配行之后
    max_lines: 500                         # 一个事件最多包含的行数
    timeout: 5s                            # 等待更多行的超时时间
```

### 4.3 JSON 日志处理

对于 JSON 格式的日志，Filebeat 提供了自动解析功能：

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/json-logs/*.log
  json:
    keys_under_root: true    # 将 JSON 字段置于根级别
    message_key: message     # 指定消息字段
    add_error_key: true      # 添加解析错误字段
    expand_keys: true        # 展开嵌套的 JSON 对象
    overwrite_keys: false    # 是否覆盖现有字段
```

### 4.4 容器日志收集

从 Docker 容器收集日志：

```yaml
filebeat.inputs:
- type: container
  enabled: true
  paths:
    - /var/lib/docker/containers/*/*.log
  exclude_lines: ["^\\s+[\\-`('.|_]"]  # 排除分隔线
  processors:
    - add_docker_metadata:
        host: "unix:///var/run/docker.sock"
```

### 4.5 其他输入类型

Filebeat 还支持其他多种输入类型：

```yaml
# Syslog 输入
- type: syslog
  enabled: true
  protocol.tcp:
    host: "localhost:9000"

# TCP 输入
- type: tcp
  enabled: true
  host: "localhost:9000"
  max_message_size: 10MB

# Kafka 输入
- type: kafka
  enabled: true
  hosts: ["kafka:9092"]
  topics: ["logs"]
  group_id: "filebeat"
```

## 5. 处理器配置

处理器可以处理、过滤和增强事件，在发送到输出前进行处理。

### 5.1 添加元数据

```yaml
processors:
  - add_host_metadata: ~     # 添加主机元数据
  - add_cloud_metadata: ~    # 添加云服务提供商元数据
  - add_docker_metadata: ~   # 添加 Docker 元数据
  - add_kubernetes_metadata: # 添加 Kubernetes 元数据
      host: ${NODE_NAME}
      matchers:
      - logs_path:
          logs_path: "/var/log/containers/"
```

### 5.2 字段处理

```yaml
processors:
  # 添加字段
  - add_fields:
      target: ''
      fields:
        environment: production
        service: backend
  
  # 删除字段
  - drop_fields:
      fields: ["agent.ephemeral_id", "ecs.version"]
  
  # 重命名字段
  - rename:
      fields:
        - from: "old_field"
          to: "new_field"
      fail_on_error: false
      ignore_missing: true
```

### 5.3 事件过滤

```yaml
processors:
  # 根据条件丢弃事件
  - drop_event:
      when:
        regexp:
          message: "^DEBUG"
  
  # 根据条件包含事件
  - include_fields:
      fields: ["message", "timestamp", "level"]
  
  # 截断字段
  - truncate_fields:
      fields:
        - message
      max_bytes: 1024
```

### 5.4 复杂处理示例

```yaml
processors:
  - drop_event:
      when:
        or:
          - equals:
              status: debug
          - regexp:
              message: "^DEBUG"
  
  - add_tags:
      tags: [json]
      target: "tags"
      when:
        contains:
          message: "{"
  
  # 使用 Grok 提取字段
  - dissect:
      tokenizer: "%{timestamp} %{level} %{message}"
      field: "log"
      target_prefix: "parsed"
```

## 6. 输出配置

Filebeat 支持多种输出目标，最常用的是 Elasticsearch 和 Logstash。

### 6.1 Elasticsearch 输出

```yaml
output.elasticsearch:
  hosts: ["https://elasticsearch:9200"]
  protocol: "https"
  username: "${ES_USERNAME}"
  password: "${ES_PASSWORD}"
  
  # 索引配置
  index: "filebeat-%{[agent.version]}-%{+yyyy.MM.dd}"
  
  # SSL/TLS 配置
  ssl.certificate_authorities: ["/etc/pki/root/ca.pem"]
  ssl.certificate: "/etc/pki/client/cert.pem"
  ssl.key: "/etc/pki/client/cert.key"
  
  # 缓冲和重试
  bulk_max_size: 50
  worker: 3
  max_retries: 3
  retry.max_duration: "30s"
  
  # 高级设置
  parameters:
    pipeline: "my-pipeline"
  compression_level: 5
```

### 6.2 Logstash 输出

```yaml
output.logstash:
  hosts: ["logstash:5044"]
  loadbalance: true
  
  # SSL/TLS 配置
  ssl.enabled: true
  ssl.certificate_authorities: ["/etc/pki/root/ca.pem"]
  ssl.certificate: "/etc/pki/client/cert.pem"
  ssl.key: "/etc/pki/client/cert.key"
  
  # 缓冲和重试
  bulk_max_size: 2048
  worker: 4
  pipelining: 2
  slow_start: true
```

### 6.3 Kafka 输出

```yaml
output.kafka:
  hosts: ["kafka1:9092", "kafka2:9092"]
  topic: "logs"
  partition.round_robin:
    reachable_only: true
  required_acks: 1
  compression: gzip
  max_message_bytes: 1000000
```

### 6.4 文件输出 (调试用)

```yaml
output.file:
  path: "/tmp/filebeat"
  filename: filebeat
  rotate_every_kb: 10000
  number_of_files: 7
```

### 6.5 控制台输出 (调试用)

```yaml
output.console:
  pretty: true
```

### 6.6 多输出配置

使用条件输出到不同目标：

```yaml
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  indices:
    - index: "apache-%{+yyyy.MM.dd}"
      when.contains:
        source: "/var/log/apache/"
    - index: "nginx-%{+yyyy.MM.dd}"
      when.contains:
        source: "/var/log/nginx/"
```

## 7. 模块系统

Filebeat 模块提供了针对常见日志类型的快速配置。

### 7.1 启用模块

```bash
# 列出可用模块
filebeat modules list

# 启用模块
filebeat modules enable system nginx mysql
```

或在配置文件中：

```yaml
filebeat.modules:
- module: system
  enabled: true
- module: nginx
  enabled: true
```

### 7.2 常用模块参数

```yaml
filebeat.modules:
- module: system
  enabled: true
  var.paths: ["/var/log/syslog*", "/var/log/auth.log*"]
  
- module: nginx
  enabled: true
  var.paths: ["/var/log/nginx/access.log*"]
  var.pipeline: "nginx-access-pipeline"
```

### 7.3 自定义模块配置

```yaml
filebeat.modules:
- module: apache
  access:
    enabled: true
    var.paths: ["/path/to/apache/access.log*"]
  error:
    enabled: true
    var.paths: ["/path/to/apache/error.log*"]
  
- module: mysql
  error:
    enabled: true
    var.paths: ["/var/log/mysql/error.log*"]
  slowlog:
    enabled: false
```

## 8. 多日志源管理

### 8.1 使用标签区分日志源

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/app1/*.log
  tags: ["app1"]
  fields:
    application: app1
    environment: production
  fields_under_root: false

- type: log
  enabled: true
  paths:
    - /var/log/app2/*.log
  tags: ["app2"]
  fields:
    application: app2
    environment: production
  fields_under_root: false
```

### 8.2 根据条件路由不同输出

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/important/*.log
  tags: ["critical"]

- type: log
  enabled: true
  paths:
    - /var/log/regular/*.log
  tags: ["regular"]

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  indices:
    - index: "critical-%{+yyyy.MM.dd}"
      when.contains:
        tags: "critical"
    - index: "regular-%{+yyyy.MM.dd}"
      when.contains:
        tags: "regular"
```

### 8.3 使用独立配置文件

```yaml
# filebeat.yml
filebeat.config.inputs:
  enabled: true
  path: configs/*.yml
  reload.enabled: true
  reload.period: 10s

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
```

然后在 `configs/` 目录中创建单独的配置文件：

```yaml
# configs/app1.yml
- type: log
  enabled: true
  paths:
    - /var/log/app1/*.log
  fields:
    application: app1
```

## 9. 性能优化

### 9.1 资源使用调优

```yaml
# 收集器设置
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/*.log
  harvester_buffer_size: 16384
  max_bytes: 10485760
  close_inactive: 5m
  scan_frequency: 10s

# 注册表设置
filebeat.registry.flush: 5s

# 事件计数设置
filebeat.spool_size: 2048
filebeat.idle_timeout: 5s

# 内存队列设置
queue.mem:
  events: 4096
  flush.min_events: 512
  flush.timeout: 5s
```

### 9.2 输出性能调优

```yaml
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  worker: 4
  bulk_max_size: 2048
  compression_level: 4

# 或 Logstash 输出
output.logstash:
  hosts: ["logstash:5044"]
  worker: 4
  bulk_max_size: 2048
  pipelining: 2
  slow_start: true
```

### 9.3 降低 CPU 使用率

```yaml
# 降低扫描频率
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/*.log
  scan_frequency: 30s  # 默认 10s

# 使用 multiline.flush_pattern 减少内存消耗
filebeat.inputs:
- type: log
  multiline:
    pattern: '^\['
    negate: false
    match: after
    flush_pattern: '^\[FLUSH\]'
```

### 9.4 限制内存使用

```yaml
# 限制单个收集器的内存使用
filebeat.inputs:
- type: log
  harvester_buffer_size: 8192  # 默认 16384
  max_bytes: 5242880           # 默认 10MB

# 限制内存队列大小
queue.mem:
  events: 2048                 # 默认 4096
```

## 10. 监控与管理

### 10.1 启用内部监控

```yaml
# 启用监控
monitoring.enabled: true

# 发送到 Elasticsearch
monitoring.elasticsearch:
  hosts: ["elasticsearch-monitoring:9200"]
  username: "beats_monitor"
  password: "your-password"
```

### 10.2 HTTP 端点监控

```yaml
# 启用 HTTP 端点
http.enabled: true
http.host: localhost
http.port: 5066

# 设置基本认证
http.username: admin
http.password: changeme
```

可以通过以下端点访问指标：

- `http://localhost:5066/` - 状态页
- `http://localhost:5066/stats` - 统计信息
- `http://localhost:5066/metrics` - Prometheus 格式的指标

### 10.3 日志配置

```yaml
logging.level: info
logging.to_files: true
logging.files:
  path: /var/log/filebeat
  name: filebeat
  keepfiles: 7
  permissions: 0600
```

### 10.4 使用 API 的常见操作

```bash
# 查看索引的文档数量
curl -X GET "localhost:9200/filebeat-*/_count"

# 查看索引的映射
curl -X GET "localhost:9200/filebeat-*/_mapping"

# 查看索引模板
curl -X GET "localhost:9200/_template/filebeat-*"
```

## 11. 安全性配置

### 11.1 传输加密 (TLS/SSL)

```yaml
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  protocol: "https"
  
  # CA 证书验证
  ssl.certificate_authorities: ["/etc/filebeat/ca.pem"]
  
  # 客户端证书
  ssl.certificate: "/etc/filebeat/client.pem"
  ssl.key: "/etc/filebeat/client.key"
  
  # 其他 SSL 选项
  ssl.verification_mode: full
  ssl.supported_protocols: [TLSv1.2, TLSv1.3]
```

### 11.2 认证配置

```yaml
# 基本认证
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  username: "${ES_USERNAME}"
  password: "${ES_PASSWORD}"

# API 密钥认证
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  api_key: "${ES_API_KEY}"
```

### 11.3 敏感信息处理

```yaml
# 使用环境变量
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  username: "${ES_USERNAME}"
  password: "${ES_PASSWORD}"

# 或使用 keystore
# 首先添加密钥到 keystore
# filebeat keystore add ES_PASSWORD
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  username: "${ES_USERNAME}"
  password: "${ES_PASSWORD}"
```

### 11.4 数据脱敏

```yaml
processors:
  - drop_fields:
      fields: ["sensitive_field", "password", "credit_card"]
  
  # 部分字段脱敏
  - script:
      lang: javascript
      source: >
        function process(event) {
          var email = event.Get("user.email");
          if (email) {
            event.Put("user.email", email.replace(/(.{2})(.*)(@.*)/, "$1***$3"));
          }
          return event;
        }
```

## 12. 常见问题排查

### 12.1 没有收集到日志

检查点：

1. 文件路径是否正确
   ```bash
   ls -la /path/to/logs/*.log
   ```

2. Filebeat 是否有权限读取文件
   ```bash
   sudo usermod -a -G loggroup filebeat  # 添加 filebeat 用户到日志组
   ```

3. 检查注册表状态
   ```bash
   filebeat registry export
   ```

4. 检查是否符合 `ignore_older` 设置
   ```yaml
   filebeat.inputs:
   - type: log
     ignore_older: 0  # 临时设为 0 以测试
   ```

### 12.2 日志解析问题

对于 JSON 解析问题：

1. 验证日志是否为有效的 JSON
   ```bash
   cat /path/to/log.json | jq
   ```

2. 启用详细调试
   ```yaml
   logging.level: debug
   logging.selectors: ["json"]
   ```

3. 使用控制台输出测试
   ```yaml
   output.console:
     pretty: true
   ```

### 12.3 连接问题

1. 测试网络连接
   ```bash
   telnet elasticsearch 9200
   ```

2. 测试输出配置
   ```bash
   filebeat test output
   ```

3. 检查 TLS/SSL 配置
   ```bash
   openssl s_client -connect elasticsearch:9200
   ```

### 12.4 性能问题

1. 检查 CPU 和内存使用
   ```bash
   top -p $(pgrep -f filebeat)
   ```

2. 检查收集的文件数
   ```bash
   filebeat --once -d "publish"
   ```

3. 调整扫描频率和关闭逻辑
   ```yaml
   filebeat.inputs:
   - type: log
     scan_frequency: 30s
     close_inactive: 5m
   ```

### 12.5 常见错误消息与解决方案

- **"Failed to connect to backoff"**: 检查网络连接和输出目标是否可达
- **"Permission denied"**: 检查文件权限，确保 Filebeat 可以读取日志文件
- **"Harvester could not be started"**: 检查文件存在性和访问权限
- **"Error loading pipeline"**: 检查 Elasticsearch Ingest Pipeline 配置
- **"X.509 certificate"**: 检查 SSL/TLS 证书配置

## 13. 配置示例集

### 13.1 基础收集与转发

最简单的配置，收集系统日志并发送到 Elasticsearch：

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/syslog
    - /var/log/auth.log

output.elasticsearch:
  hosts: ["localhost:9200"]
```

### 13.2 多种日志格式处理

处理不同格式的日志：

```yaml
filebeat.inputs:
# JSON 日志
- type: log
  enabled: true
  paths:
    - /var/log/app/*.json
  json.keys_under_root: true
  json.add_error_key: true
  json.message_key: message
  fields:
    type: json-logs

# 多行日志（如 Java 堆栈）
- type: log
  enabled: true
  paths:
    - /var/log/java-app/*.log
  multiline:
    pattern: '^[0-9]{4}-[0-9]{2}-[0-9]{2}'
    negate: true
    match: after
  fields:
    type: java-logs

# Apache 日志
- type: log
  enabled: true
  paths:
    - /var/log/apache2/access.log
  fields:
    type: apache-access

output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "logs-%{[fields.type]}-%{+yyyy.MM.dd}"
```

### 13.3 生产环境完整配置

适用于生产环境的完整配置：

```yaml
# ===== Filebeat 全局配置 =====
filebeat.registry.flush: 5s
filebeat.shutdown_timeout: 30s
filebeat.config.modules:
  path: ${path.config}/modules.d/*.yml
  reload.enabled: true
  reload.period: 10s

# ===== 输入配置 =====
filebeat.inputs:
# 应用程序日志
- type: log
  enabled: true
  paths:
    - /var/log/app/*.log
  json.keys_under_root: true
  json.add_error_key: true
  fields:
    service: my-service
    environment: production
  fields_under_root: false
  tags: ["app", "json"]

# 系统日志
- type: log
  enabled: true
  paths:
    - /var/log/syslog
    - /var/log/auth.log
  tags: ["system"]

# ===== 处理器配置 =====
processors:
  - add_host_metadata: ~
  - add_cloud_metadata: ~
  - add_docker_metadata: ~
  - add_kubernetes_metadata:
      host: ${NODE_NAME}
      matchers:
      - logs_path:
          logs_path: "/var/log/containers/"
          
  # 添加字段
  - add_fields:
      target: ''
      fields:
        datacenter: east-1
        
  # 仅包含指定字段
  - include_fields:
      fields: ["@timestamp", "message", "log.*", "service", "host", "agent", "cloud"]
      
  # 敏感数据处理
  - drop_fields:
      fields: ["password", "token", "authorization"]

# ===== 输出配置 =====
output.elasticsearch:
  hosts: ["https://elasticsearch-1:9200", "https://elasticsearch-2:9200"]
  username: "${ES_USERNAME}"
  password: "${ES_PASSWORD}"
  index: "%{[service]}-%{+yyyy.MM.dd}"
  ssl.certificate_authorities: ["/etc/filebeat/ca.pem"]
  bulk_max_size: 1024
  worker: 4
  compression_level: 5

# 失败时的备用输出
output.file:
  enabled: false
  path: "/tmp/filebeat-failure"
  filename: filebeat-errors
  rotate_every_kb: 10000
  number_of_files: 7

# ===== 监控配置 =====
monitoring:
  enabled: true
  elasticsearch:
    hosts: ["https://monitoring-es:9200"]
    username: "${MONITORING_USERNAME}"
    password: "${MONITORING_PASSWORD}"

# ===== 日志配置 =====
logging.level: info
logging.to_files: true
logging.files:
  path: /var/log/filebeat
  name: filebeat
  keepfiles: 7
  permissions: 0600

# ===== HTTP 端点 =====
http.enabled: true
http.host: localhost
http.port: 5066
```

### 13.4 Kubernetes 环境配置

在 Kubernetes 中收集容器日志：

```yaml
filebeat.autodiscover:
  providers:
    - type: kubernetes
      node: ${NODE_NAME}
      hints.enabled: true
      hints.default_config:
        type: container
        paths:
          - /var/log/containers/*-${data.kubernetes.container.id}.log
      scope: node

processors:
  - add_cloud_metadata: ~
  - add_host_metadata: ~
  - add_kubernetes_metadata:
      host: ${NODE_NAME}
      matchers:
      - logs_path:
          logs_path: "/var/log/containers/"
          
  - drop_event:
      when:
        or:
          - equals:
              kubernetes.container.name: "filebeat"
          - equals:
              kubernetes.namespace: "kube-system"

output.elasticsearch:
  hosts: ['${ELASTICSEARCH_HOST:elasticsearch}:${ELASTICSEARCH_PORT:9200}']
  username: ${ELASTICSEARCH_USERNAME}
  password: ${ELASTICSEARCH_PASSWORD}
  index: "kubernetes-%{[agent.version]}-%{+yyyy.MM.dd}"
```

### 13.5 Docker 容器环境配置

适用于 Docker 环境的配置：

```yaml
filebeat.inputs:
- type: container
  enabled: true
  paths:
    - '/var/lib/docker/containers/*/*.log'
  json.keys_under_root: true
  json.message_key: log
  json.add_error_key: true
  processors:
    - add_docker_metadata:
        host: "unix:///var/run/docker.sock"
    # 提取容器标签作为字段
    - add_fields:
        target: container
        fields:
          app: '${data.docker.container.labels.app}'
          env: '${data.docker.container.labels.environment}'

# 降低内存使用量
filebeat.spool_size: 1024
filebeat.idle_timeout: 2s
queue.mem:
  events: 2048
  flush.min_events: 256
  flush.timeout: 2s

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  index: "docker-logs-%{+yyyy.MM.dd}"
```

---

*本文档提供了 Filebeat 配置的全面指南，覆盖了从基本安装到高级配置的各个方面。根据您的具体环境和需求，您可能需要调整某些配置项。* 