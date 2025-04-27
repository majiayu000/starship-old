#!/bin/bash

# 设置颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}更新前端依赖...${NC}"

# 进入前端目录
cd frontend

# 检查 node_modules 是否存在
if [ -d "node_modules" ]; then
  echo -e "${BLUE}删除旧的 node_modules...${NC}"
  rm -rf node_modules
fi

# 检查 package-lock.json 是否存在
if [ -f "package-lock.json" ]; then
  echo -e "${BLUE}删除旧的 package-lock.json...${NC}"
  rm package-lock.json
fi

# 安装依赖
echo -e "${BLUE}安装依赖...${NC}"
npm install

# 检查是否成功
if [ $? -eq 0 ]; then
  echo -e "${GREEN}前端依赖更新成功!${NC}"
else
  echo -e "${RED}前端依赖更新失败!${NC}"
  exit 1
fi

# 返回上级目录
cd ..

echo -e "${GREEN}完成!${NC}"
