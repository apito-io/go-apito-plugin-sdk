# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
