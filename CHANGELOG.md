# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.21] - 2025-01-03

### Added

- **GraphQL Error Support**: Complete support for proper GraphQL error handling in mutations and queries
- **GraphQLError Types**: New error types specifically designed for GraphQL operations
  - `GraphQLError`: Main GraphQL error structure with message, extensions, path, and locations
  - `GraphQLErrorLocation`: Error location information for GraphQL errors
  - `GraphQLErrorExtensions`: Common error extension fields
- **GraphQL Error Constructors**: Comprehensive set of error creation functions
  - `GraphQLErrorWithMessage()`: Basic GraphQL error with message
  - `GraphQLErrorWithCode()`: GraphQL error with error code in extensions
  - `GraphQLErrorWithExtensions()`: GraphQL error with custom extensions
  - `GraphQLValidationError()`: Validation error for form/input validation
  - `GraphQLAuthenticationError()`: Authentication required error
  - `GraphQLAuthorizationError()`: Authorization/permission denied error
  - `GraphQLNotFoundError()`: Resource not found error
  - `GraphQLInternalError()`: Internal server error
  - `GraphQLBadUserInputError()`: Bad user input error
- **Resolver Helper Functions**: Convenient functions for returning GraphQL errors from resolvers
  - `ReturnGraphQLError()`: Return any GraphQL error from resolver
  - `ReturnValidationError()`: Return validation error from resolver
  - `ReturnAuthenticationError()`: Return authentication error from resolver
  - `ReturnAuthorizationError()`: Return authorization error from resolver
  - `ReturnNotFoundError()`: Return not found error from resolver
  - `ReturnInternalError()`: Return internal error from resolver
  - `ReturnBadUserInputError()`: Return bad input error from resolver
  - `ValidateAndReturn()`: Conditional validation with error or success result
  - `ValidateFieldAndReturn()`: Field validation with error or success result
  - `HandleErrorAndReturn()`: Convert any error to appropriate GraphQL error
- **Error Detection Functions**: Utility functions for error type checking
  - `IsGraphQLError()`: Check if error is a GraphQL error
  - `IsCodedError()`: Check if error is a coded HTTP error
  - `GetGraphQLError()`: Safely extract GraphQL error details

### Enhanced

- **Execute Method**: Modified to handle GraphQL errors differently from REST API errors
  - GraphQL operations now return properly formatted GraphQL error responses
  - REST API operations continue to use HTTP status code based error handling
  - Automatic conversion of regular errors to GraphQL format for GraphQL operations
- **Error Handling**: Backward compatible error handling that works with existing plugins
- **Developer Experience**: Simplified error handling in GraphQL resolvers with intuitive helper functions

### Fixed

- **GraphQL Mutation Errors**: Resolved issue where GraphQL mutation errors returned null response objects instead of proper GraphQL errors
- **Error Response Format**: Fixed error response format to comply with GraphQL specification

### Usage Example

```go
// In your GraphQL resolver function:
func (s *SuchokFunction) _changeEmployeePassword(ctx context.Context, userID, tenantID string, input map[string]interface{}) (interface{}, error) {
    // Validate input
    if employeeID := sdk.GetStringArg(input, "employee_id"); employeeID == "" {
        return sdk.ReturnValidationError("Employee ID is required", "employee_id")
    }
    
    // Check authentication
    if !isAuthenticated(userID) {
        return sdk.ReturnAuthenticationError("You must be logged in to change passwords")
    }
    
    // Check authorization
    if !hasPermission(userID, "change_password") {
        return sdk.ReturnAuthorizationError("You don't have permission to change employee passwords")
    }
    
    // Business logic...
    if err := changePassword(employeeID, newPassword); err != nil {
        return sdk.HandleErrorAndReturn(err, "Failed to change employee password")
    }
    
    return map[string]interface{}{
        "success": true,
        "message": "Password changed successfully",
    }, nil
}
```

### Breaking Changes

- None - All changes are backward compatible

## [0.1.17] - 2025-01-30

### Added

- **Array Object Type Helpers**: New convenience functions for creating GraphQL fields that return arrays of objects
- **NewArrayObjectType()**: Creates a GraphQL field that returns an array of the specified object type
- **NewArrayObjectTypeWithArgs()**: Creates a GraphQL field with arguments that returns an array of objects
- **Enhanced Documentation**: Added comprehensive examples for array object types in README.md

### Enhanced

- **Type System**: Improved support for complex object arrays in GraphQL schema building
- **Developer Experience**: Simplified creation of array object fields with convenience wrappers
- **Example Integration**: Added complete working example in hc-hello-world-plugin demonstrating array object usage

### Usage Example

```go
// Define an object type
taskType := sdk.NewObjectType("Task", "A task object").
    AddStringField("id", "Task ID", false).
    AddStringField("title", "Task title", false).
    AddBooleanField("completed", "Whether task is completed", false).
    Build()

// Method 1: Simple array field
plugin.RegisterQuery("getTasks",
    sdk.NewArrayObjectType(taskType),
    getTasksResolver)

// Method 2: Array field with arguments
plugin.RegisterQuery("getFilteredTasks",
    sdk.NewArrayObjectTypeWithArgs(taskType, map[string]interface{}{
        "status": sdk.StringArg("Filter by task status"),
        "limit":  sdk.IntArg("Maximum number of tasks to return"),
    }),
    getFilteredTasksResolver)
```

### Technical Features

- **Backward Compatibility**: New functions build on existing `ListOfObjectsFieldWithArgs` foundation
- **Zero Breaking Changes**: All existing functionality remains unchanged
- **Type Safety**: Maintains strong typing for object array responses
- **Documentation**: Complete examples and usage patterns included

## [0.1.16] - 2024-12-28

### Added

- **Enhanced File Upload Support**: Complete set of file upload helper functions for REST API endpoints
- **File Schema Functions**: `FileSchema()` and `MultipartFormSchema()` for defining file upload schemas
- **REST Endpoint File Upload Methods**: `WithFileUpload()` and `WithMultipartForm()` for configuring file uploads
- **File Content Extraction**: `GetFileUpload()`, `GetFileUploadBytes()`, `GetFileUploadInfo()` for accessing uploaded files
- **Base64 Content Handling**: Automatic detection and decoding of base64-encoded file content
- **Form Field Extraction**: `GetMultipartFormValue()` for extracting form data from multipart uploads

### Enhanced

- **LogRESTArgs()**: Upgraded to structured, categorized logging with emoji indicators and file upload detection
- **File Content Processing**: Handles both `[]byte` and base64 string content automatically
- **Parameter Parsing**: Enhanced `ParseRESTArgs()` for better organization of path, query, and body parameters
- **Error Handling**: Improved error messages and debugging information for file upload scenarios

### Technical Features

- **Backward Compatibility**: All existing functions continue to work unchanged - NO BREAKING CHANGES
- **Multiple Content Types**: Support for various file formats with proper MIME type detection
- **Size Handling**: Proper handling of file sizes as both `int64` and `float64` from JSON
- **Debug Logging**: Comprehensive logging with 📁, 📄, 🔍, 📦 emoji indicators for easy debugging

## [0.1.14] - 2025-01-26

### Enhanced

- **Comprehensive Health Check System**: Significantly upgraded the built-in health check functionality
- **Custom Health Checks**: Added `RegisterHealthCheck()` and `RegisterHealthChecks()` methods for plugin-specific health monitoring
- **Runtime Metrics**: Enhanced health check response includes runtime information (memory usage, goroutines, GC cycles, Go version)
- **Plugin Statistics**: Added detailed statistics about registered queries, mutations, REST APIs, functions, and object types
- **Environment Information**: Health check now reports PID, hostname, OS, and architecture information
- **Degraded Status Detection**: Smart status monitoring that detects degraded states from custom health checks
- **Flexible Health Check API**: New `HealthCheckFunc` type allows plugins to implement custom health monitoring logic

## [0.1.13] - 2024-01-11

### Added

- **Typed Array Helper Functions**: New helper functions for typed array argument extraction
- **GetStringArrayArg()**: Extracts string arrays with proper type conversion from `[]interface{}`
- **GetIntArrayArg()**: Extracts integer arrays with conversion from strings and floats
- **GetFloatArrayArg()**: Extracts float arrays with conversion from strings and integers
- **GetBoolArrayArg()**: Extracts boolean arrays with smart conversion from strings and numbers

### Enhanced

- **Type Safety**: All array helpers handle `[]interface{}` to typed array conversion automatically
- **Flexible Conversion**: Support for converting between compatible types (string to int, int to float, etc.)
- **Backward Compatibility**: Existing `GetArrayArg()` function remains unchanged

### Problem Solved

- Fixed issue where `GetArrayArg()` returns `[]interface{}` but actual values are typed arrays
- Eliminates empty array issues when dealing with typed arguments from GraphQL/REST APIs
- Provides safe type conversion with fallback mechanisms

### Usage Example

```go
// Before: orderIDs would be empty due to type mismatch
orderIDsRaw := sdk.GetArrayArg(args, "order_ids")  // Returns []interface{}

// After: proper string array extraction
orderIDs := sdk.GetStringArrayArg(args, "order_ids")  // Returns []string
```

## [0.1.11] - 2024-12-19

### Fixed

- **CRITICAL**: Fixed `GetBodyParam*` functions to properly handle `body_` prefix from engine
- **Code Simplification**: Refactored helper functions to reuse existing code and eliminate duplication
- **Parameter Extraction**: All body parameter helpers now correctly extract parameters sent with `body_` prefix

### Technical Details

- Simplified all `GetBodyParam*` functions to use existing `GetStringArg`, `GetIntArg`, etc. instead of duplicating logic
- Maintains backward compatibility while properly handling new engine parameter format
- Cleaner, more maintainable code without repetition

## [0.1.10] - 2024-12-19

### Added

- **REST API Helper Functions**: Comprehensive set of helper functions for parsing REST API parameters
- **Path Parameter Helpers**: `GetPathParam()` for extracting path parameters (`:id`, `:userId`, etc.)
- **Query Parameter Helpers**: `GetQueryParam()`, `GetQueryParamBool()`, `GetQueryParamInt()` for query string parameters
- **Body Parameter Helpers**: `GetBodyParam()`, `GetBodyParamInt()`, `GetBodyParamBool()`, `GetBodyParamObject()`, `GetBodyParamArray()` for request body data
- **Unified REST Parser**: `ParseRESTArgs()` categorizes all parameters into path, query, and body sections
- **Debug Logging**: `LogRESTArgs()` provides structured logging for REST API debugging
- **Endpoint Info**: `GetRESTEndpointInfo()` extracts HTTP method, path, and request metadata

### Enhanced

- **Type Safety**: All REST helpers include type conversion and validation
- **Flexible Input**: Handles multiple parameter naming conventions (with/without prefixes)
- **Default Values**: Support for default values in all parameter extraction functions
- **Structured Debugging**: Categorized parameter logging for easier troubleshooting

### Technical Features

- **Multiple Format Support**: Handles `:param`, `path_param`, `query_param` naming patterns
- **Boolean Parsing**: Smart boolean conversion from strings ("true", "1", "yes")
- **Integer Conversion**: Automatic conversion from strings and floats to integers
- **Context Integration**: Works with existing context helper functions

## [0.1.9] - 2024-12-19

### Fixed

- **REST API Function Routing**: Fixed compatibility with engine's new REST API function naming convention
- **Backward Compatibility**: Added support for both old (`METHOD_/path`) and new (`rest_method_path`) function naming formats
- **Automatic Conversion**: SDK now automatically converts new format function names to internal handler keys
- **Path Parameter Handling**: Properly handles path parameters like `:id` in function name conversion

### Enhanced

- **Debug Logging**: Added comprehensive logging for REST API function resolution
- **Error Messages**: Improved error messages to show available handlers when lookup fails
- **Smart Fallback**: Tries direct lookup first, then attempts format conversion
- **Documentation**: Updated README with detailed explanation of function name compatibility

### Technical Changes

- Updated `
