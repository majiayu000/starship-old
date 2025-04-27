#!/bin/bash

# 设置颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 设置默认端口
BACKEND_PORT=8080
FRONTEND_PORT=5555

# 显示帮助信息
function show_help {
  echo -e "${BLUE}全栈应用启动脚本${NC}"
  echo -e "用法: $0 [选项]"
  echo -e "选项:"
  echo -e "  -h, --help         显示帮助信息"
  echo -e "  -b, --backend-only 仅启动后端服务"
  echo -e "  -f, --frontend-only 仅启动前端服务"
  echo -e "  --backend-port PORT 设置后端端口 (默认: ${BACKEND_PORT})"
  echo -e "  --frontend-port PORT 设置前端端口 (默认: ${FRONTEND_PORT})"
  echo -e "  --no-browser       不自动打开浏览器"
  exit 0
}

# 解析命令行参数
BACKEND_ONLY=false
FRONTEND_ONLY=false
OPEN_BROWSER=true

while [[ $# -gt 0 ]]; do
  case $1 in
    -h|--help)
      show_help
      ;;
    -b|--backend-only)
      BACKEND_ONLY=true
      shift
      ;;
    -f|--frontend-only)
      FRONTEND_ONLY=true
      shift
      ;;
    --backend-port)
      BACKEND_PORT="$2"
      shift 2
      ;;
    --frontend-port)
      FRONTEND_PORT="$2"
      shift 2
      ;;
    --no-browser)
      OPEN_BROWSER=false
      shift
      ;;
    *)
      echo -e "${RED}未知选项: $1${NC}"
      show_help
      ;;
  esac
done

# 检查是否同时设置了仅后端和仅前端
if [[ "$BACKEND_ONLY" == true && "$FRONTEND_ONLY" == true ]]; then
  echo -e "${RED}错误: 不能同时设置 --backend-only 和 --frontend-only${NC}"
  exit 1
fi

echo -e "${GREEN}启动全栈应用...${NC}"

# 启动后端
if [[ "$FRONTEND_ONLY" == false ]]; then
  echo -e "${BLUE}启动后端服务 (端口: ${BACKEND_PORT})...${NC}"

  # 检查Go是否安装
  if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: Go 未安装，请先安装 Go${NC}"
    exit 1
  fi

  cd backend
  echo -e "${BLUE}编译后端...${NC}"
  go build -o api_binary ./cmd/api || {
    echo -e "${RED}后端编译失败${NC}"
    exit 1
  }

  # 设置环境变量
  export SERVER_PORT=$BACKEND_PORT

  echo -e "${BLUE}运行后端服务...${NC}"
  ./api_binary &
  BACKEND_PID=$!
  cd ..

  # 检查后端是否成功启动
  echo -e "${BLUE}等待后端服务启动...${NC}"
  sleep 2
  for i in {1..5}; do
    if curl -s http://localhost:${BACKEND_PORT}/health &> /dev/null; then
      echo -e "${GREEN}后端服务已成功启动 (端口: ${BACKEND_PORT})${NC}"
      break
    fi

    if [[ $i -eq 5 ]]; then
      echo -e "${YELLOW}警告: 无法确认后端服务是否成功启动，但将继续执行${NC}"
    else
      echo -e "${BLUE}等待后端服务启动 (尝试 $i/5)...${NC}"
      sleep 2
    fi
  done
fi

# 启动前端
if [[ "$BACKEND_ONLY" == false ]]; then
  echo -e "${BLUE}启动前端服务 (端口: ${FRONTEND_PORT})...${NC}"

  # 检查Node.js是否安装
  if ! command -v node &> /dev/null; then
    echo -e "${RED}错误: Node.js 未安装，请先安装 Node.js${NC}"
    exit 1
  fi

  cd frontend

  # 检查依赖是否已安装
  if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}前端依赖未安装，正在安装...${NC}"
    npm install || {
      echo -e "${RED}依赖安装失败${NC}"
      # 如果后端已启动，则关闭
      if [[ "$FRONTEND_ONLY" == false && "$BACKEND_PID" != "" ]]; then
        kill $BACKEND_PID
      fi
      exit 1
    }
  fi

  # 设置环境变量
  export PORT=$FRONTEND_PORT
  export NEXT_PUBLIC_API_URL=http://localhost:${BACKEND_PORT}/api/v1

  echo -e "${BLUE}运行前端服务...${NC}"
  npm run dev &
  FRONTEND_PID=$!
  cd ..

  # 检查前端是否成功启动
  echo -e "${BLUE}等待前端服务启动...${NC}"
  sleep 3
  for i in {1..5}; do
    if curl -s http://localhost:${FRONTEND_PORT} &> /dev/null; then
      echo -e "${GREEN}前端服务已成功启动 (端口: ${FRONTEND_PORT})${NC}"
      break
    fi

    if [[ $i -eq 5 ]]; then
      echo -e "${YELLOW}警告: 无法确认前端服务是否成功启动，但将继续执行${NC}"
    else
      echo -e "${BLUE}等待前端服务启动 (尝试 $i/5)...${NC}"
      sleep 2
    fi
  done

  # 打开浏览器
  if [[ "$OPEN_BROWSER" == true ]]; then
    echo -e "${GREEN}在浏览器中打开应用...${NC}"
    if [[ "$OSTYPE" == "darwin"* ]]; then
      open http://localhost:${FRONTEND_PORT}
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
      xdg-open http://localhost:${FRONTEND_PORT} &> /dev/null || {
        echo -e "${YELLOW}无法自动打开浏览器，请手动访问: http://localhost:${FRONTEND_PORT}${NC}"
      }
    elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
      start http://localhost:${FRONTEND_PORT} || {
        echo -e "${YELLOW}无法自动打开浏览器，请手动访问: http://localhost:${FRONTEND_PORT}${NC}"
      }
    else
      echo -e "${YELLOW}无法自动打开浏览器，请手动访问: http://localhost:${FRONTEND_PORT}${NC}"
    fi
  else
    echo -e "${BLUE}应用已启动，请访问: http://localhost:${FRONTEND_PORT}${NC}"
  fi
fi

# 设置关闭处理
function cleanup {
  echo -e "\n${GREEN}正在关闭服务...${NC}"

  if [[ "$FRONTEND_ONLY" == false && -n "$BACKEND_PID" ]]; then
    echo -e "${BLUE}关闭后端服务...${NC}"
    kill $BACKEND_PID &> /dev/null || echo -e "${YELLOW}后端服务已经关闭${NC}"
  fi

  if [[ "$BACKEND_ONLY" == false && -n "$FRONTEND_PID" ]]; then
    echo -e "${BLUE}关闭前端服务...${NC}"
    kill $FRONTEND_PID &> /dev/null || echo -e "${YELLOW}前端服务已经关闭${NC}"
  fi

  echo -e "${GREEN}所有服务已关闭${NC}"
  exit 0
}

# 捕获中断信号
trap cleanup SIGINT SIGTERM

# 显示服务状态
if [[ "$BACKEND_ONLY" == true ]]; then
  echo -e "${GREEN}后端服务已启动! 按 Ctrl+C 停止服务${NC}"
  echo -e "${BLUE}API 地址: http://localhost:${BACKEND_PORT}${NC}"
elif [[ "$FRONTEND_ONLY" == true ]]; then
  echo -e "${GREEN}前端服务已启动! 按 Ctrl+C 停止服务${NC}"
  echo -e "${BLUE}前端地址: http://localhost:${FRONTEND_PORT}${NC}"
else
  echo -e "${GREEN}全栈应用已启动! 按 Ctrl+C 停止所有服务${NC}"
  echo -e "${BLUE}前端地址: http://localhost:${FRONTEND_PORT}${NC}"
  echo -e "${BLUE}API 地址: http://localhost:${BACKEND_PORT}${NC}"
fi

# 保持脚本运行
wait
