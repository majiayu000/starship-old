#!/bin/bash

# 更新导入路径
find backend -type f -name "*.go" -exec sed -i '' 's|github.com/majiayu000/cc-starship|github.com/yourusername/starter-template|g' {} \;

echo "导入路径已更新"
