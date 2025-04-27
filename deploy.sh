#!/bin/bash

# 设置颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 显示帮助信息
function show_help {
  echo -e "${BLUE}全栈应用部署脚本${NC}"
  echo -e "用法: $0 [选项]"
  echo -e "选项:"
  echo -e "  -h, --help         显示帮助信息"
  echo -e "  -b, --backend-only 仅部署后端"
  echo -e "  -f, --frontend-only 仅部署前端"
  exit 0
}

# 解析命令行参数
BACKEND_ONLY=false
FRONTEND_ONLY=false

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

echo -e "${GREEN}开始部署全栈应用...${NC}"

# 部署后端
if [[ "$FRONTEND_ONLY" == false ]]; then
  echo -e "${BLUE}部署后端...${NC}"
  
  # 检查Go是否安装
  if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: Go 未安装，请先安装 Go${NC}"
    exit 1
  fi
  
  cd backend
  
  # 更新依赖
  echo -e "${BLUE}更新后端依赖...${NC}"
  go mod tidy
  
  # 编译后端
  echo -e "${BLUE}编译后端...${NC}"
  go build -o api_binary ./cmd/api || {
    echo -e "${RED}后端编译失败${NC}"
    exit 1
  }
  
  echo -e "${GREEN}后端编译成功!${NC}"
  
  # 创建 Docker 镜像
  if command -v docker &> /dev/null; then
    echo -e "${BLUE}创建 Docker 镜像...${NC}"
    docker build -t myapp-backend . || {
      echo -e "${YELLOW}警告: Docker 镜像创建失败，但将继续执行${NC}"
    }
  else
    echo -e "${YELLOW}警告: Docker 未安装，跳过创建 Docker 镜像${NC}"
  fi
  
  cd ..
fi

# 部署前端
if [[ "$BACKEND_ONLY" == false ]]; then
  echo -e "${BLUE}部署前端...${NC}"
  
  # 检查Node.js是否安装
  if ! command -v node &> /dev/null; then
    echo -e "${RED}错误: Node.js 未安装，请先安装 Node.js${NC}"
    exit 1
  fi
  
  cd frontend
  
  # 安装依赖
  echo -e "${BLUE}安装前端依赖...${NC}"
  npm install || {
    echo -e "${RED}前端依赖安装失败${NC}"
    exit 1
  }
  
  # 构建前端
  echo -e "${BLUE}构建前端...${NC}"
  npm run build || {
    echo -e "${RED}前端构建失败${NC}"
    exit 1
  }
  
  echo -e "${GREEN}前端构建成功!${NC}"
  
  # 创建 Docker 镜像
  if command -v docker &> /dev/null; then
    echo -e "${BLUE}创建 Docker 镜像...${NC}"
    # 检查是否存在 Dockerfile
    if [ -f "Dockerfile" ]; then
      docker build -t myapp-frontend . || {
        echo -e "${YELLOW}警告: Docker 镜像创建失败，但将继续执行${NC}"
      }
    else
      echo -e "${YELLOW}警告: 前端目录中没有 Dockerfile，跳过创建 Docker 镜像${NC}"
    fi
  else
    echo -e "${YELLOW}警告: Docker 未安装，跳过创建 Docker 镜像${NC}"
  fi
  
  cd ..
fi

echo -e "${GREEN}部署完成!${NC}"

# 显示部署信息
if [[ "$BACKEND_ONLY" == true ]]; then
  echo -e "${BLUE}后端已部署:${NC}"
  echo -e "  - 二进制文件: ${GREEN}backend/api_binary${NC}"
  if command -v docker &> /dev/null; then
    echo -e "  - Docker 镜像: ${GREEN}myapp-backend${NC}"
  fi
elif [[ "$FRONTEND_ONLY" == true ]]; then
  echo -e "${BLUE}前端已部署:${NC}"
  echo -e "  - 构建目录: ${GREEN}frontend/.next${NC}"
  if command -v docker &> /dev/null && [ -f "frontend/Dockerfile" ]; then
    echo -e "  - Docker 镜像: ${GREEN}myapp-frontend${NC}"
  fi
else
  echo -e "${BLUE}全栈应用已部署:${NC}"
  echo -e "  - 后端二进制文件: ${GREEN}backend/api_binary${NC}"
  echo -e "  - 前端构建目录: ${GREEN}frontend/.next${NC}"
  if command -v docker &> /dev/null; then
    echo -e "  - 后端 Docker 镜像: ${GREEN}myapp-backend${NC}"
    if [ -f "frontend/Dockerfile" ]; then
      echo -e "  - 前端 Docker 镜像: ${GREEN}myapp-frontend${NC}"
    fi
  fi
fi

echo -e "\n${BLUE}启动应用:${NC}"
echo -e "  - 后端: ${GREEN}cd backend && ./api_binary${NC}"
echo -e "  - 前端: ${GREEN}cd frontend && npm start${NC}"

if command -v docker &> /dev/null; then
  echo -e "\n${BLUE}使用 Docker 启动:${NC}"
  echo -e "  - 后端: ${GREEN}docker run -p 8080:8080 myapp-backend${NC}"
  if [ -f "frontend/Dockerfile" ]; then
    echo -e "  - 前端: ${GREEN}docker run -p 5555:5555 myapp-frontend${NC}"
  fi
fi
