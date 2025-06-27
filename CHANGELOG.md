# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

### File Upload Example

```go
// Register endpoint with file upload support
plugin.RegisterRESTAPI(sdk.POSTEndpoint("/share/invoice", "Upload invoice PDF").
    WithFileUpload("pdf", "Invoice PDF file", map[string]interface{}{
        "customer_phone": sdk.StringSchema("Customer phone number"),
        "invoice_no":     sdk.StringSchema("Invoice number"),
        "order_id":       sdk.StringSchema("Order ID"),
    }).Build(), shareInvoiceHandler)

// In handler function
func shareInvoiceHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    // Extract file content (handles both []byte and base64 automatically)
    fileContent := sdk.GetFileUploadBytes(args, "pdf")
    filename, contentType, size := sdk.GetFileUploadInfo(args, "pdf")

    // Extract form fields
    customerPhone := sdk.GetMultipartFormValue(args, "customer_phone")

    // Process file upload...
    return uploadResponse, nil
}
```

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
