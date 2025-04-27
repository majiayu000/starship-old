# Backend Architecture

This document outlines the architecture of the backend application.

## Overview

The backend follows a clean architecture approach, separating concerns into distinct layers:

1. **API Layer**: Handles HTTP requests and responses
2. **Core Layer**: Contains business logic and domain models
3. **Infrastructure Layer**: Provides implementations for external services
4. **Repository Layer**: Handles data access and persistence

## Layers

### API Layer

The API layer is responsible for handling HTTP requests and responses. It includes:

- **Handlers**: Process HTTP requests and return responses
- **Middleware**: Provide cross-cutting concerns like authentication, logging, etc.
- **Routes**: Define the API endpoints

### Core Layer

The core layer contains the business logic and domain models. It includes:

- **Domain**: Define the domain models and business rules
- **Ports**: Define interfaces for external dependencies
- **Services**: Implement business logic

### Infrastructure Layer

The infrastructure layer provides implementations for external services. It includes:

- **Auth**: Authentication and authorization services
- **Cache**: Caching services
- **Database**: Database connections and management

### Repository Layer

The repository layer handles data access and persistence. It includes:

- **Postgres**: PostgreSQL repository implementations

## Flow of Control

1. HTTP request comes in through the API layer
2. API handler calls the appropriate service in the core layer
3. Service uses repositories through ports to access data
4. Service applies business logic and returns result
5. API handler formats the response and returns it

## Dependency Injection

The application uses manual dependency injection to wire up the components:

1. Infrastructure components are created first
2. Repositories are created with infrastructure dependencies
3. Services are created with repository dependencies
4. API handlers are created with service dependencies

This approach ensures that dependencies flow inward, with the core layer having no dependencies on outer layers.
