# Full-Stack Starter Template

A modern, production-ready starter template for building full-stack web applications with Go (backend) and Next.js (frontend).

## Features

### Backend (Go)

- **Clean Architecture**: Follows the ports and adapters pattern for maintainable code
- **API Layer**: RESTful API with Gin framework
- **Database**: PostgreSQL with connection pooling
- **Authentication**: JWT-based authentication
- **Logging**: Structured logging with multiple outputs
- **Configuration**: Environment-based configuration with Viper
- **Error Handling**: Consistent error handling and responses
- **Middleware**: Common middleware for authentication, logging, etc.

### Frontend (Next.js)

- **React 18**: Modern React with hooks and functional components
- **TypeScript**: Type-safe code
- **UI Components**: Shadcn UI components
- **API Client**: Fetch-based API client
- **State Management**: React hooks for state management
- **Form Handling**: Form validation and submission
- **Notifications**: Toast notifications
- **Responsive Design**: Mobile-first responsive design

## Project Structure

```
.
├── backend/                 # Go backend
│   ├── cmd/                 # Application entry points
│   ├── configs/             # Configuration files
│   ├── internal/            # Internal packages
│   │   ├── api/             # API layer
│   │   ├── core/            # Core business logic
│   │   ├── infrastructure/  # Infrastructure concerns
│   │   └── repositories/    # Data access
│   ├── pkg/                 # Shared packages
│   └── scripts/             # Utility scripts
└── frontend/                # Next.js frontend
    ├── public/              # Static assets
    └── src/                 # Source code
        ├── api/             # API clients
        ├── app/             # Next.js app router
        ├── components/      # React components
        ├── hooks/           # Custom hooks
        └── models/          # TypeScript models
```

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 18+
- PostgreSQL 14+

### Quick Start

使用我们的一键启动脚本来同时启动前后端服务：

```bash
./start.sh
```

这个脚本会自动编译后端、启动前后端服务，并在浏览器中打开应用。

脚本支持多种选项，可以通过 `./start.sh --help` 查看帮助信息：

```
全栈应用启动脚本
用法: ./start.sh [选项]
选项:
  -h, --help         显示帮助信息
  -b, --backend-only 仅启动后端服务
  -f, --frontend-only 仅启动前端服务
  --backend-port PORT 设置后端端口 (默认: 8080)
  --frontend-port PORT 设置前端端口 (默认: 5555)
  --no-browser       不自动打开浏览器
```

### 手动设置

#### 后端设置

1. 进入后端目录：
   ```bash
   cd backend
   ```

2. 安装依赖：
   ```bash
   go mod tidy
   ```

   或者使用我们的脚本：
   ```bash
   ./update-go-deps.sh
   ```

3. 设置数据库：
   ```bash
   # 创建数据库
   createdb myapp

   # 运行迁移
   go run cmd/migrate/main.go
   ```

4. 启动服务器：
   ```bash
   go run cmd/api/main.go
   ```

API 服务器将在 http://localhost:8080 上可用。

#### 前端设置

1. 进入前端目录：
   ```bash
   cd frontend
   ```

2. 安装依赖：
   ```bash
   npm install
   ```

   或者使用我们的脚本：
   ```bash
   ./update-frontend-deps.sh
   ```

3. 启动开发服务器：
   ```bash
   npm run dev
   ```

前端将在 http://localhost:5555 上可用。

## Configuration

### Backend Configuration

The backend is configured using environment variables or a configuration file (`configs/app.yaml`). See `pkg/config/config.go` for available options.

### Frontend Configuration

The frontend is configured using environment variables in a `.env.local` file:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## 部署

### 一键部署

使用我们的一键部署脚本来构建和部署应用：

```bash
./deploy.sh
```

这个脚本会自动编译后端、构建前端，并创建 Docker 镜像（如果 Docker 已安装）。

脚本支持多种选项，可以通过 `./deploy.sh --help` 查看帮助信息：

```
全栈应用部署脚本
用法: ./deploy.sh [选项]
选项:
  -h, --help         显示帮助信息
  -b, --backend-only 仅部署后端
  -f, --frontend-only 仅部署前端
```

### Docker Compose 部署

我们提供了 Docker Compose 配置，可以一键启动整个应用栈，包括后端、前端、PostgreSQL 和 Redis：

```bash
docker-compose up -d
```

这将启动所有服务，并在后台运行。

要查看日志：

```bash
docker-compose logs -f
```

要停止所有服务：

```bash
docker-compose down
```

### 手动部署

#### 后端部署

后端可以作为 Docker 容器部署：

```bash
cd backend
docker build -t myapp-backend .
docker run -p 8080:8080 myapp-backend
```

或者直接运行编译后的二进制文件：

```bash
cd backend
go build -o api_binary ./cmd/api
./api_binary
```

#### 前端部署

前端可以部署到 Vercel 或作为 Docker 容器：

```bash
cd frontend
docker build -t myapp-frontend .
docker run -p 5555:5555 myapp-frontend
```

或者直接运行构建后的应用：

```bash
cd frontend
npm run build
npm start
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
