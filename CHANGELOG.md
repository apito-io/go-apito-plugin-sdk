# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

- Updated `Execute()` method in `sdk.go` to handle both REST API naming conventions
- Added path reconstruction logic for complex REST endpoints
- Enhanced error handling with detailed debugging information

## [0.1.6] - 2024-12-19

### Fixed

- **Pagination Type Issue**: Fixed `PaginatedResponseType` to avoid undefined type references that caused GraphQL schema errors
- **Type Reference Problem**: Simplified paginated response to use scalar fields instead of complex object references
- **Schema Validation**: Resolved "fields must be an object" errors in complex nested types

### Changed

- `PaginatedResponseType()` now creates a self-contained type with scalar fields instead of referencing other object types
- Simplified pagination structure to avoid circular type dependencies

### Technical Notes

- This is a temporary fix to resolve immediate schema validation issues
- Future versions will implement proper type resolution for complex nested object references

## [0.1.5] - 2024-12-19

### Fixed

- **Engine Compatibility**: Fixed schema registration issues where complex GraphQL types weren't being recognized by the engine
- **Type System**: Improved compatibility between SDK's structured GraphQL types and engine's type conversion system
- **Backwards Compatibility**: Maintained full backwards compatibility with existing string-based type definitions

### Technical Changes

- Enhanced engine's `convertGraphQLType` function to handle both string and structured type definitions
- Added new `convertGraphQLTypeFromData`, `convertStructuredGraphQLType`, and related helper functions
- Updated schema registration to properly convert SDK's `GraphQLTypeDefinition` objects to native GraphQL types
- Improved argument type handling in `extractGraphQLArgs` function

## [0.1.4] - 2024-12-19

### Added

- **Proper GraphQL Type System**: Complete rewrite of type handling to generate proper GraphQL type structures
- **Engine Compatibility**: Types now generate proper `GraphQLTypeDefinition` objects instead of simple strings
- **Complex Type Support**: Full support for nested objects, arrays, and non-null types
- **Type Documentation**: Added `TYPE_SYSTEM.md` with comprehensive type structure documentation

### Technical Changes

- Replaced string-based types with structured `GraphQLTypeDefinition` objects
- Added `createScalarType()`, `createObjectType()`, `createListType()`, `createNonNullType()` helpers
- Enhanced serialization to handle complex type structures
- Maintained backwards compatibility with string types

## [0.1.3] - 2024-12-19

### Added

- **Complex Type Support**: Added comprehensive support for complex return types beyond simple strings
- **Object Type System**: New `ObjectTypeDefinition` and `ObjectFieldDef` for defining complex structures
- **Type Builders**: Fluent API with `ObjectTypeBuilder` for building complex types
- **Built-in Types**: Common types like `UserObjectType()`, `PaginationInfoType()`, `ErrorObjectType()`
- **Advanced Field Functions**: `ComplexObjectField()`, `ListOfObjectsField()`, `PaginatedResponseType()`
- **Documentation**: Added `COMPLEX_TYPES_EXAMPLES.md` with comprehensive examples

### Enhanced

- Updated hello-world plugin to demonstrate all complex type features
- Added examples for objects, arrays, pagination, and error handling

## [0.1.1] - 2024-12-19

### Fixed

- **Critical Bug**: Fixed protobuf serialization error when registering GraphQL schemas
- **Type Serialization**: Added proper serialization for `GraphQLField` and nested argument structures
- **Engine Compatibility**: Resolved "proto: invalid type: sdk.GraphQLField" error

### Added

- `serializeGraphQLField()` method for converting GraphQL fields to protobuf-compatible maps
- `serializeArgs()` method for recursive argument serialization
- `serializeValue()` method for handling nested maps, arrays, and custom types

### Technical Changes

- Updated `SchemaRegister()` to use proper serialization instead of direct field assignment
- Enhanced error handling and debugging output
- Improved protobuf compatibility

## [0.1.0] - 2024-01-20

### Added

- Initial release of the Apito Plugin SDK
- Simplified API for creating HashiCorp plugins for Apito Engine
- Support for GraphQL queries and mutations registration
- Support for REST API endpoints registration
- Support for custom function registration
- Helper functions for common GraphQL field types
- Builder pattern for REST endpoint definitions
- Comprehensive documentation and examples
- Example plugin demonstrating SDK usage
- Automatic handling of gRPC and HashiCorp plugin boilerplate

### Features

- **Simple Initialization**: `sdk.Init(name, version, apiKey)`
- **GraphQL Registration**: Easy registration of queries and mutations
- **REST API Registration**: Fluent builder pattern for endpoints
- **Function Registration**: Custom function support
- **Type Helpers**: Built-in helpers for common GraphQL types
- **Batch Registration**: Register multiple items at once
- **Error Handling**: Proper error handling and context support

### Developer Experience

- Reduces plugin boilerplate from ~600 lines to ~100 lines
- Declarative API design
- Type-safe function signatures
- Comprehensive helper functions
- Clear documentation with examples
