添加认证信息的 API 调用
所有与 Elasticsearch 的 API 交互都需要添加认证信息。根据您的 filebeat.yml 配置，您使用的是用户名 elastic 和密码 changeme。
正确的 ILM 策略创建命令:

curl -X PUT "localhost:9200/_ilm/policy/filebeat" \
  -H 'Content-Type: application/json' \
  -u elastic:changeme \
  -d'
{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_age": "7d",
            "max_size": "10gb"
          },
          "set_priority": {
            "priority": 100
          }
        }
      },
      "warm": {
        "min_age": "30d",
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
        "min_age": "60d",
        "actions": {
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
'