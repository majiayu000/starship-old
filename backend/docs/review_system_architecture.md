# Review System Architecture

This document describes the architecture of the review system designed for handling different data sources like SAT Oneprep and IXL.

## Overview

The review system is designed to be extensible, allowing for easy addition of new data sources while maintaining a consistent API for frontend interactions. It follows a clean architecture pattern with clear separation of concerns.

## Architecture Layers

### 1. Domain Layer

The domain layer contains the core business entities and business rules.

- **Domain Models**:
  - `ReviewableItem`: Base struct with common fields for all reviewable items
  - `SATOneprepItem`: Specialization for Oneprep SAT data
  - `SATIXLItem`: Specialization for IXL SAT data
  - `ReviewStatusUpdate`: DTO for updating review status
  - `QueryParams`: DTO for filtering and pagination
  - `PaginatedResult`: DTO for paginated results

### 2. Repository Layer

The repository layer handles data access and persistence.

- **Interfaces**:
  - `ReviewRepository`: Interface for data access operations
  
- **Implementations**:
  - `PostgresReviewRepository`: Implementation using PostgreSQL

The repository layer abstracts away the data source details, allowing the service layer to work with domain objects without knowing where the data comes from.

### 3. Service Layer

The service layer contains business logic and orchestration.

- **Interfaces**:
  - `ReviewService`: Interface for business operations
  
- **Implementations**:
  - `ReviewServiceImpl`: Implementation of business logic

The service layer enforces business rules, validates input, and coordinates operations between the repository and the API.

### 4. API Layer

The API layer handles HTTP requests and responses.

- **Handlers**:
  - `ReviewHandler`: HTTP handler for review operations
  
- **Routes**:
  - `/api/v1/review/sources`: Get available data sources
  - `/api/v1/review/:dataSource/filters`: Get filter options for a data source
  - `/api/v1/review/:dataSource/items`: Get all items with pagination and filtering
  - `/api/v1/review/:dataSource/items/:id`: Get a specific item
  - `/api/v1/review/:dataSource/items/:id/review`: Review an item (update status)

## Data Flow

1. The client sends a request to one of the review endpoints
2. The handler parses the request and calls the appropriate service method
3. The service validates the input and calls the repository
4. The repository executes the database query and returns the result
5. The service processes the result and returns it to the handler
6. The handler formats the response and sends it to the client

## Extending for New Data Sources

To add a new data source:

1. Add a new constant in the `domain` package for the data source
2. Create a new domain model by embedding `ReviewableItem`
3. Add a new entry to the `sourceToTable` map in the repository
4. Implement the scanning methods for the new data source
5. Update the `IsValidDataSource` function

## Error Handling

The system uses a consistent error handling approach:

1. Each layer catches errors from the layer below and adds context
2. The handler layer translates errors into appropriate HTTP status codes
3. Errors are logged for debugging and monitoring

## Configuration

The system is configurable through the application's configuration system:

1. Database connection details
2. Pagination defaults
3. Logging level
4. CORS settings
5. Authentication requirements

## Authentication and Authorization

The system integrates with the existing authentication system:

1. Review endpoints can be protected by the authentication middleware
2. Authorization can be enforced based on user roles
3. The reviewer's identity is captured when a review is submitted

## Monitoring and Logging

The system includes comprehensive logging:

1. All operations are logged with appropriate context
2. Errors are logged with full details
3. Performance metrics can be collected for monitoring

## Conclusion

This architecture provides a flexible and extensible foundation for the review system, allowing for easy addition of new data sources while maintaining a consistent API for frontend interactions. 