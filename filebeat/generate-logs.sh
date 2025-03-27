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