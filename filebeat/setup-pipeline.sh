#!/bin/bash

# 创建 pipeline
curl -X PUT "localhost:9200/_ingest/pipeline/timezone-conversion" \
  -H "Content-Type: application/json" \
  -u elastic:changeme \
  -d @pipeline.json 