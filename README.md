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

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up the database:
   ```bash
   # Create the database
   createdb myapp
   
   # Run migrations
   go run cmd/migrate/main.go
   ```

4. Start the server:
   ```bash
   go run cmd/api/main.go
   ```

The API server will be available at http://localhost:8080.

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

The frontend will be available at http://localhost:3000.

## Configuration

### Backend Configuration

The backend is configured using environment variables or a configuration file (`configs/app.yaml`). See `pkg/config/config.go` for available options.

### Frontend Configuration

The frontend is configured using environment variables in a `.env.local` file:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Deployment

### Backend Deployment

The backend can be deployed as a Docker container:

```bash
cd backend
docker build -t myapp-backend .
docker run -p 8080:8080 myapp-backend
```

### Frontend Deployment

The frontend can be deployed to Vercel or as a Docker container:

```bash
cd frontend
npm run build
npm start
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
