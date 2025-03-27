# Implementation Steps for Review System

This document outlines the necessary steps to fully implement the review system in the application.

## 1. Update Database Tables

Ensure the database tables have the necessary columns for review functionality:

```sql
-- For t_math_sat_temp_classify_oneprep table, add or ensure these columns exist:
ALTER TABLE t_math_sat_temp_classify_oneprep
ADD COLUMN IF NOT EXISTS review_status TEXT NOT NULL DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS review_comment TEXT,
ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS reviewed_by TEXT;

-- Create similar columns for other tables that need review functionality
```

## 2. Update Application Configuration

Add configuration for the review system in `config.yaml` or environment variables:

```yaml
features:
  enableReview: true

review:
  dataSources:
    - name: sat_oneprep
      table: t_math_sat_temp_classify_oneprep
      enabled: true
    - name: sat_ixl
      table: t_math_sat_temp_classify_ixl
      enabled: true
  defaultPageSize: 10
```

## 3. Wire Up Dependencies in Main Application

Update the application's dependency injection in `cmd/api/main.go`:

```go
// Create review repository
reviewRepo := postgres.NewReviewRepository(db, logger)

// Create review service
reviewService := services.NewReviewService(reviewRepo, logger)

// Create HTTP server with review service
server := api.NewServer(
    cfg,
    logger,
    userService,
    authService,
    logManager,
    cacheService,
    reviewService,
)
```

## 4. Frontend Implementation

### 4.1. API Client

Create API client functions in the frontend:

```typescript
// Get all items with pagination and filtering
async function getReviewItems(dataSource, params) {
  const queryString = new URLSearchParams(params).toString();
  const response = await fetch(`/api/v1/review/${dataSource}/items?${queryString}`);
  return response.json();
}

// Get a specific item
async function getReviewItem(dataSource, id) {
  const response = await fetch(`/api/v1/review/${dataSource}/items/${id}`);
  return response.json();
}

// Review an item
async function reviewItem(dataSource, id, reviewData) {
  const response = await fetch(`/api/v1/review/${dataSource}/items/${id}/review`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(reviewData)
  });
  return response.json();
}

// Get filter options
async function getFilterOptions(dataSource) {
  const response = await fetch(`/api/v1/review/${dataSource}/filters`);
  return response.json();
}

// Get available data sources
async function getDataSources() {
  const response = await fetch('/api/v1/review/sources');
  return response.json();
}
```

### 4.2. UI Components

Create the following UI components:

1. **Data Source Selector**: Dropdown to select the data source
2. **Review List**: Table or card list showing items with pagination
3. **Filter Panel**: UI for applying filters (by status, domain, skill, etc.)
4. **Review Detail**: Modal or page for reviewing a specific item
5. **Review Form**: Form for submitting reviews with status and comments

### 4.3. State Management

Set up state management (Redux, Context API, etc.) to handle:

1. Current data source
2. Filter criteria
3. Pagination state
4. Current items
5. Loading states
6. Error states

## 5. Testing

### 5.1. Unit Tests

Create unit tests for:

1. Domain models
2. Service layer
3. Repository layer
4. API handlers

### 5.2. Integration Tests

Create integration tests for:

1. Database operations
2. API endpoints

### 5.3. End-to-End Tests

Create end-to-end tests for:

1. Complete review workflow
2. Error handling
3. Authentication

## 6. Deployment

### 6.1. Database Migrations

Create database migrations for adding review columns to existing tables.

### 6.2. Environment Configuration

Set up environment-specific configuration for the review system.

### 6.3. Monitoring

Set up monitoring for:

1. API endpoint usage
2. Error rates
3. Database performance

## 7. Documentation

### 7.1. API Documentation

Document the API endpoints using Swagger or similar tool.

### 7.2. User Guide

Create a user guide for the review system.

### 7.3. Developer Guide

Create a developer guide for extending the review system for new data sources. 