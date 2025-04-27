#!/bin/bash

# 设置颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}启动全栈应用...${NC}"

# 启动后端
echo -e "${BLUE}启动后端服务...${NC}"
cd backend
go build -o api_binary ./cmd/api
./api_binary &
BACKEND_PID=$!
cd ..

# 等待后端启动
echo -e "${BLUE}等待后端服务启动...${NC}"
sleep 3

# 启动前端
echo -e "${BLUE}启动前端服务...${NC}"
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

# 等待前端启动
echo -e "${BLUE}等待前端服务启动...${NC}"
sleep 5

# 打开浏览器
echo -e "${GREEN}在浏览器中打开应用...${NC}"
open http://localhost:3000

# 设置关闭处理
function cleanup {
  echo -e "${GREEN}关闭服务...${NC}"
  kill $BACKEND_PID
  kill $FRONTEND_PID
  exit 0
}

# 捕获中断信号
trap cleanup SIGINT

# 保持脚本运行
echo -e "${GREEN}服务已启动! 按 Ctrl+C 停止所有服务${NC}"
wait
