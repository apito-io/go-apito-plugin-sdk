# Apito Plugin Development and Extension Guide

This comprehensive guide covers everything you need to know about developing and extending Apito plugins using the Go Plugin SDK. This guide is designed to be a complete reference for developers and AI systems.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Plugin Architecture](#plugin-architecture)
3. [GraphQL Development](#graphql-development)
4. [Error Handling](#error-handling)
5. [REST API Development](#rest-api-development)
6. [Complex Data Types](#complex-data-types)
7. [Helper Functions](#helper-functions)
8. [Best Practices](#best-practices)
9. [Testing](#testing)
10. [Deployment](#deployment)
11. [Extension Guidelines](#extension-guidelines)

## Getting Started

### Basic Plugin Structure

Every Apito plugin follows this basic structure:

```go
package main

import (
    "context"
    "log"
    sdk "github.com/apito-io/go-apito-plugin-sdk"
)

func main() {
    // Initialize plugin
    plugin := sdk.Init("plugin-name", "1.0.0", "your-api-key")
    
    // Register your GraphQL operations
    plugin.RegisterMutation("exampleMutation", 
        sdk.ComplexObjectFieldWithArgs("Example mutation", responseType, inputArgs),
        exampleMutationResolver)
    
    // Register REST endpoints (optional)
    plugin.RegisterRESTAPI(
        sdk.POSTEndpoint("/api/example", "Example endpoint").Build(),
        exampleRESTHandler)
    
    // Start the plugin server
    plugin.Serve()
}

func exampleMutationResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
    // Your resolver logic here
    userID := sdk.GetUserID(rawArgs)
    tenantID := sdk.GetTenantID(rawArgs)
    args := sdk.ParseArgsForResolver("exampleMutation", rawArgs)
    
    // Business logic...
    
    return map[string]interface{}{
        "success": true,
        "message": "Operation completed successfully",
    }, nil
}
```

### Dependencies

Add to your `go.mod`:

```go
module your-plugin-name

go 1.21

require (
    github.com/apito-io/go-apito-plugin-sdk v0.1.21
    github.com/apito-io/types latest
)
```

## Plugin Architecture

### HashiCorp Plugin System

Apito uses HashiCorp's go-plugin system for process isolation and fault tolerance:

- **Process Isolation**: Each plugin runs in a separate process
- **RPC Communication**: Communication via gRPC
- **Fault Tolerance**: Plugin crashes don't affect the main engine
- **Hot Reloading**: Plugins can be restarted without engine downtime

### Plugin Lifecycle

1. **Initialization**: Plugin registers its schema and endpoints
2. **Schema Registration**: GraphQL schema is merged with engine schema
3. **Execution**: Plugin handles incoming requests
4. **Health Checks**: Continuous monitoring of plugin health

## GraphQL Development

### Creating Mutations

```go
// Define response type
responseType := sdk.NewObjectType("ExampleResponse", "Response for example operation").
    AddBooleanField("success", "Operation success status", false).
    AddStringField("message", "Response message", true).
    AddStringField("id", "Created resource ID", true).
    Build()

// Define input arguments
inputArgs := map[string]interface{}{
    "input": sdk.ObjectArg("Example input data", map[string]interface{}{
        "name":        sdk.StringProperty("Resource name"),
        "description": sdk.StringProperty("Resource description"),
        "active":      sdk.BooleanProperty("Active status"),
    }),
}

// Register mutation
plugin.RegisterMutation("createExample",
    sdk.ComplexObjectFieldWithArgs("Create example resource", responseType, inputArgs),
    createExampleResolver)
```

### Creating Queries

```go
// Register query
plugin.RegisterQuery("getExample",
    sdk.ComplexObjectFieldWithArgs("Get example resource", responseType, map[string]interface{}{
        "id": sdk.StringArg("Resource ID"),
    }),
    getExampleResolver)

func getExampleResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
    args := sdk.ParseArgsForResolver("getExample", rawArgs)
    id := sdk.GetStringArg(args, "id")
    
    if id == "" {
        return sdk.ReturnValidationError("ID is required", "id")
    }
    
    // Fetch data...
    
    return result, nil
}
```

### List/Array Queries

```go
// For returning arrays of objects
listResponseType := sdk.NewArrayObjectTypeWithArgs(itemType, map[string]interface{}{
    "filter": sdk.ObjectArg("Filter criteria", map[string]interface{}{
        "status": sdk.StringProperty("Status filter"),
        "limit":  sdk.IntProperty("Limit results"),
        "offset": sdk.IntProperty("Offset for pagination"),
    }),
})

plugin.RegisterQuery("listExamples", listResponseType, listExamplesResolver)
```

## Error Handling

### GraphQL Error Types

The SDK provides comprehensive GraphQL error handling:

#### Basic Error Types

```go
// Validation errors
return sdk.ReturnValidationError("Name is required", "name")

// Authentication errors
return sdk.ReturnAuthenticationError("You must be logged in")

// Authorization errors
return sdk.ReturnAuthorizationError("You don't have permission to perform this action")

// Not found errors
return sdk.ReturnNotFoundError("Resource not found")

// Internal errors
return sdk.ReturnInternalError("Something went wrong")

// Bad user input
return sdk.ReturnBadUserInputError("Invalid email format", "email")
```

#### Custom Error Extensions

```go
// Custom error with extensions
extensions := map[string]interface{}{
    "code":        "BUSINESS_RULE_VIOLATION",
    "businessRule": "INSUFFICIENT_BALANCE",
    "currentBalance": 100.00,
    "requiredAmount": 150.00,
}

return sdk.ReturnGraphQLErrorWithExtensions("Insufficient balance for this operation", extensions)
```

#### Error Conversion

```go
// Convert any error to GraphQL error
if err := someOperation(); err != nil {
    return sdk.HandleErrorAndReturn(err, "Failed to complete operation")
}
```

#### Conditional Validation

```go
// Validate conditions
return sdk.ValidateAndReturn(
    user.HasPermission("admin"), 
    "Admin permission required", 
    successResult,
    "PERMISSION_DENIED")

// Validate fields
return sdk.ValidateFieldAndReturn(email, "email", true, successResult)
```

### Error Response Format

GraphQL errors are automatically formatted according to the GraphQL specification:

```json
{
  "data": null,
  "errors": [
    {
      "message": "Validation failed",
      "extensions": {
        "code": "VALIDATION_ERROR",
        "field": "email"
      },
      "path": ["createUser"],
      "locations": [{"line": 2, "column": 3}]
    }
  ]
}
```

## REST API Development

### Creating REST Endpoints

```go
// Simple GET endpoint
endpoint := sdk.GETEndpoint("/api/examples", "List examples").
    WithResponseSchema(sdk.ObjectSchema(map[string]interface{}{
        "data": sdk.ArraySchema(sdk.ObjectSchema(map[string]interface{}{
            "id":   sdk.StringSchema("Example ID"),
            "name": sdk.StringSchema("Example name"),
        })),
    })).
    Build()

plugin.RegisterRESTAPI(endpoint, listExamplesRESTHandler)

// POST endpoint with file upload
uploadEndpoint := sdk.POSTEndpoint("/api/upload", "Upload file").
    WithFileUpload("file", "File to upload", map[string]interface{}{
        "description": sdk.StringSchema("File description"),
        "category":    sdk.StringSchema("File category"),
    }).
    Build()

plugin.RegisterRESTAPI(uploadEndpoint, uploadFileHandler)
```

### REST Handler Implementation

```go
func listExamplesRESTHandler(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
    // Parse REST arguments
    parsed := sdk.ParseRESTArgs(rawArgs)
    
    // Extract query parameters
    limit := sdk.GetQueryParamInt(rawArgs, "limit", 10)
    offset := sdk.GetQueryParamInt(rawArgs, "offset", 0)
    filter := sdk.GetQueryParam(rawArgs, "filter", "")
    
    // Extract path parameters
    id := sdk.GetPathParam(rawArgs, "id")
    
    // Business logic...
    
    return map[string]interface{}{
        "data": results,
        "pagination": map[string]interface{}{
            "limit":  limit,
            "offset": offset,
            "total":  total,
        },
    }, nil
}
```

### File Upload Handling

```go
func uploadFileHandler(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
    // Get file upload
    fileBytes := sdk.GetFileUploadBytes(rawArgs, "file")
    filename, contentType, size := sdk.GetFileUploadInfo(rawArgs, "file")
    
    // Get form data
    description := sdk.GetMultipartFormValue(rawArgs, "description")
    category := sdk.GetMultipartFormValue(rawArgs, "category")
    
    // Process file...
    
    return map[string]interface{}{
        "fileId":   fileId,
        "filename": filename,
        "size":     size,
        "url":      fileUrl,
    }, nil
}
```

## Complex Data Types

### Object Type Definition

```go
// Define complex object
userType := sdk.NewObjectType("User", "User object").
    AddStringField("id", "User ID", false).
    AddStringField("email", "User email", false).
    AddStringField("firstName", "First name", true).
    AddStringField("lastName", "Last name", true).
    AddBooleanField("active", "User active status", false).
    AddStringListField("roles", "User roles", false, true).
    Build()

// Register the type
plugin.RegisterObjectType(userType)
```

### Nested Objects

```go
// Define address type
addressType := sdk.NewObjectType("Address", "Address information").
    AddStringField("street", "Street address", false).
    AddStringField("city", "City", false).
    AddStringField("zipCode", "ZIP code", false).
    Build()

// Use in user type
userType := sdk.NewObjectType("User", "User object").
    AddStringField("id", "User ID", false).
    AddObjectField("address", "User address", addressType, true).
    Build()
```

### JSON Fields

```go
// For storing complex JSON data
userType := sdk.NewObjectType("User", "User object").
    AddStringField("id", "User ID", false).
    AddJSONField("metadata", "User metadata", userMetadataType, true).
    AddJSONArrayField("preferences", "User preferences", preferenceType, true).
    Build()
```

## Helper Functions

### Argument Extraction

```go
// String arguments
name := sdk.GetStringArg(args, "name", "default")

// Integer arguments
limit := sdk.GetIntArg(args, "limit", 10)

// Boolean arguments
active := sdk.GetBoolArg(args, "active", true)

// Object arguments
input := sdk.GetObjectArg(args, "input")

// Array arguments
tags := sdk.GetStringArrayArg(args, "tags")
ids := sdk.GetIntArrayArg(args, "ids")
```

### Context Data Access

```go
// Get context information
userID := sdk.GetUserID(rawArgs)
tenantID := sdk.GetTenantID(rawArgs)
projectID := sdk.GetProjectID(rawArgs)
pluginID := sdk.GetPluginID(rawArgs)

// Get all context data
contextData := sdk.GetAllContextData(rawArgs)
```

### Logging and Debugging

```go
// Log REST API calls
sdk.LogRESTArgs("uploadFile", rawArgs)

// Get endpoint information
endpointInfo := sdk.GetRESTEndpointInfo(rawArgs)
log.Printf("HTTP Method: %s, Path: %s", endpointInfo["method"], endpointInfo["path"])
```

## Best Practices

### 1. Error Handling

```go
// Always use appropriate error types
func createUserResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
    args := sdk.ParseArgsForResolver("createUser", rawArgs)
    input := sdk.GetObjectArg(args, "input")
    
    // Validate input
    email := sdk.GetStringArg(input, "email")
    if email == "" {
        return sdk.ReturnValidationError("Email is required", "email")
    }
    
    // Check authentication
    userID := sdk.GetUserID(rawArgs)
    if userID == "" {
        return sdk.ReturnAuthenticationError("Authentication required")
    }
    
    // Business logic with error conversion
    if err := createUser(email); err != nil {
        return sdk.HandleErrorAndReturn(err, "Failed to create user")
    }
    
    return result, nil
}
```

### 2. Input Validation

```go
// Comprehensive input validation
func validateUserInput(input map[string]interface{}) error {
    email := sdk.GetStringArg(input, "email")
    password := sdk.GetStringArg(input, "password")
    
    if email == "" {
        return sdk.GraphQLValidationError("Email is required", "email")
    }
    
    if !isValidEmail(email) {
        return sdk.GraphQLValidationError("Invalid email format", "email")
    }
    
    if len(password) < 8 {
        return sdk.GraphQLValidationError("Password must be at least 8 characters", "password")
    }
    
    return nil
}
```

### 3. Resource Management

```go
// Proper resource cleanup
func processFileResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
    fileBytes := sdk.GetFileUploadBytes(rawArgs, "file")
    
    // Create temporary file
    tempFile, err := createTempFile(fileBytes)
    if err != nil {
        return sdk.ReturnInternalError("Failed to create temporary file")
    }
    defer os.Remove(tempFile.Name()) // Cleanup
    
    // Process file...
    
    return result, nil
}
```

### 4. Pagination

```go
// Standard pagination pattern
func listItemsResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
    args := sdk.ParseArgsForResolver("listItems", rawArgs)
    
    limit := sdk.GetIntArg(args, "limit", 20)
    offset := sdk.GetIntArg(args, "offset", 0)
    
    // Validate limits
    if limit > 100 {
        return sdk.ReturnValidationError("Limit cannot exceed 100", "limit")
    }
    
    // Fetch data with pagination
    items, total, err := fetchItems(limit, offset)
    if err != nil {
        return sdk.HandleErrorAndReturn(err, "Failed to fetch items")
    }
    
    return map[string]interface{}{
        "items": items,
        "pagination": map[string]interface{}{
            "total":       total,
            "limit":       limit,
            "offset":      offset,
            "hasNext":     offset+limit < total,
            "hasPrevious": offset > 0,
        },
    }, nil
}
```

## Testing

### Unit Testing

```go
package main

import (
    "context"
    "testing"
    sdk "github.com/apito-io/go-apito-plugin-sdk"
)

func TestCreateUserResolver(t *testing.T) {
    // Setup test args
    rawArgs := map[string]interface{}{
        "input": map[string]interface{}{
            "email":     "test@example.com",
            "firstName": "John",
            "lastName":  "Doe",
        },
        "context_user_id":   "user123",
        "context_tenant_id": "tenant123",
    }
    
    // Call resolver
    result, err := createUserResolver(context.Background(), rawArgs)
    
    // Assertions
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    response := result.(map[string]interface{})
    if response["success"] != true {
        t.Errorf("Expected success to be true")
    }
}

func TestCreateUserValidation(t *testing.T) {
    // Test missing email
    rawArgs := map[string]interface{}{
        "input": map[string]interface{}{},
        "context_user_id": "user123",
    }
    
    _, err := createUserResolver(context.Background(), rawArgs)
    
    // Should return validation error
    if !sdk.IsGraphQLError(err) {
        t.Errorf("Expected GraphQL error, got %v", err)
    }
    
    gqlErr := sdk.GetGraphQLError(err)
    if gqlErr.Extensions["code"] != "VALIDATION_ERROR" {
        t.Errorf("Expected validation error code")
    }
}
```

### Integration Testing

```go
func TestPluginIntegration(t *testing.T) {
    // Initialize plugin
    plugin := sdk.Init("test-plugin", "1.0.0", "test-key")
    
    // Register test operations
    plugin.RegisterMutation("testMutation", testField, testResolver)
    
    // Test schema registration
    if _, exists := plugin.GetMutationField("testMutation"); !exists {
        t.Error("Mutation not registered")
    }
    
    // Test health check
    healthResult, err := plugin.performHealthCheck(context.Background())
    if err != nil {
        t.Fatalf("Health check failed: %v", err)
    }
    
    if healthResult["status"] != "healthy" {
        t.Error("Plugin not healthy")
    }
}
```

## Deployment

### Build Configuration

```makefile
# Makefile example
PLUGIN_NAME=your-plugin
VERSION=1.0.0

build:
	go mod tidy
	go build -o $(PLUGIN_NAME) .

test:
	go test -v ./...

clean:
	rm -f $(PLUGIN_NAME)

.PHONY: build test clean
```

### Docker Support

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o plugin .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/plugin .
CMD ["./plugin"]
```

### Environment Configuration

```go
// Environment-based configuration
func getConfig() Config {
    return Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        APIKey:      os.Getenv("API_KEY"),
        Debug:       os.Getenv("DEBUG") == "true",
    }
}
```

## Extension Guidelines

### Adding New Functionality

When extending the plugin:

1. **Maintain Backward Compatibility**: Never break existing APIs
2. **Follow Naming Conventions**: Use consistent naming patterns
3. **Add Comprehensive Tests**: Test all new functionality
4. **Update Documentation**: Keep docs current
5. **Version Appropriately**: Follow semantic versioning

### SDK Extension Example

```go
// If extending the SDK itself, follow this pattern:

// 1. Add new error types (if needed)
type CustomError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// 2. Add helper functions
func NewCustomError(code, message string, data interface{}) error {
    return &CustomError{Code: code, Message: message, Data: data}
}

// 3. Add to Execute method handling
if customErr, ok := err.(*CustomError); ok {
    // Handle custom error type
}

// 4. Add convenience functions
func ReturnCustomError(code, message string, data interface{}) (interface{}, error) {
    return nil, NewCustomError(code, message, data)
}
```

### Plugin Extension Checklist

- [ ] Functionality implemented and tested
- [ ] Error handling follows SDK patterns
- [ ] Input validation comprehensive
- [ ] Documentation updated
- [ ] Backward compatibility maintained
- [ ] Performance implications considered
- [ ] Security implications reviewed
- [ ] Integration tests pass

## Common Patterns

### Authentication Middleware

```go
func requireAuth(next sdk.ResolverFunc) sdk.ResolverFunc {
    return func(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
        userID := sdk.GetUserID(rawArgs)
        if userID == "" {
            return sdk.ReturnAuthenticationError("Authentication required")
        }
        return next(ctx, rawArgs)
    }
}

// Usage
plugin.RegisterMutation("protectedOperation",
    field,
    requireAuth(protectedOperationResolver))
```

### Caching Pattern

```go
var cache = make(map[string]interface{})

func cachedResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
    args := sdk.ParseArgsForResolver("cachedOperation", rawArgs)
    key := generateCacheKey(args)
    
    // Check cache
    if cached, exists := cache[key]; exists {
        return cached, nil
    }
    
    // Fetch data
    result, err := fetchExpensiveData(args)
    if err != nil {
        return sdk.HandleErrorAndReturn(err, "Failed to fetch data")
    }
    
    // Cache result
    cache[key] = result
    return result, nil
}
```

### Batch Operations

```go
func batchCreateResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
    args := sdk.ParseArgsForResolver("batchCreate", rawArgs)
    items := sdk.GetArrayObjectArg(args, "items")
    
    if len(items) > 100 {
        return sdk.ReturnValidationError("Cannot process more than 100 items at once", "items")
    }
    
    results := make([]interface{}, 0, len(items))
    errors := make([]interface{}, 0)
    
    for i, item := range items {
        result, err := createSingleItem(item)
        if err != nil {
            errors = append(errors, map[string]interface{}{
                "index": i,
                "error": err.Error(),
            })
            continue
        }
        results = append(results, result)
    }
    
    return map[string]interface{}{
        "results":    results,
        "errors":     errors,
        "successful": len(results),
        "failed":     len(errors),
    }, nil
}
```

This guide provides comprehensive coverage of plugin development for the Apito platform. Use it as a reference for creating new plugins or extending existing ones. Always refer to the latest SDK documentation for any updates or new features.
