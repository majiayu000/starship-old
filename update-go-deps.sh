#!/bin/bash

# 设置颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}更新 Go 依赖...${NC}"

# 进入后端目录
cd backend

# 运行 go mod tidy
echo -e "${BLUE}运行 go mod tidy...${NC}"
go mod tidy

# 检查是否成功
if [ $? -eq 0 ]; then
  echo -e "${GREEN}Go 依赖更新成功!${NC}"
else
  echo -e "${RED}Go 依赖更新失败!${NC}"
  exit 1
fi

# 返回上级目录
cd ..

echo -e "${GREEN}完成!${NC}"
