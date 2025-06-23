# Apito Plugin SDK

A simplified SDK for building HashiCorp plugins for the Apito Engine. This SDK abstracts away all the boilerplate code and provides a clean, easy-to-use interface for plugin developers.

## Installation

```bash
go mod init your-plugin-name
go get github.com/apito-io/go-apito-plugin-sdk
```

## Quick Start

### Basic Plugin Structure

```go
package main

import (
    "context"
    "fmt"

    "github.com/apito-io/go-apito-plugin-sdk"
)

func main() {
    // Initialize the plugin
    plugin := sdk.Init("my-awesome-plugin", "1.0.0", "your-api-key")

    // Register GraphQL queries
    plugin.RegisterQuery("hello",
        sdk.FieldWithArgs("String", "Returns a greeting", map[string]interface{}{
            "name": sdk.StringArg("Name to greet"),
        }),
        helloResolver,
    )

    // Register GraphQL mutations
    plugin.RegisterMutation("createUser",
        sdk.FieldWithArgs("String", "Creates a new user", map[string]interface{}{
            "user": sdk.ObjectArg("User data", map[string]interface{}{
                "name":  sdk.StringProperty("User name"),
                "email": sdk.StringProperty("User email"),
                "age":   sdk.IntProperty("User age"),
            }),
        }),
        createUserResolver,
    )

    // Register REST API endpoints
    plugin.RegisterRESTAPI(
        sdk.GETEndpoint("/hello", "Simple hello endpoint").
            WithResponseSchema(sdk.ObjectSchema(map[string]interface{}{
                "message": sdk.StringSchema("Hello message"),
                "timestamp": sdk.StringSchema("Current timestamp"),
            })).
            Build(),
        helloRESTHandler,
    )

    // Register custom functions
    plugin.RegisterFunction("processData", processDataFunction)

    // Start the plugin server
    plugin.Serve()
}

// GraphQL Resolvers
func helloResolver(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    name := "World"
    if nameArg, ok := args["name"].(string); ok && nameArg != "" {
        name = nameArg
    }
    return fmt.Sprintf("Hello, %s!", name), nil
}

func createUserResolver(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    if user, ok := args["user"].(map[string]interface{}); ok {
        name := user["name"].(string)
        email := user["email"].(string)
        age := int(user["age"].(float64))

        return fmt.Sprintf("Created user: %s <%s> (age: %d)", name, email, age), nil
    }
    return nil, fmt.Errorf("invalid user data")
}

// REST Handlers
func helloRESTHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    return map[string]interface{}{
        "message":   "Hello from REST API!",
        "timestamp": "2024-01-01T00:00:00Z",
    }, nil
}

// Custom Functions
func processDataFunction(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    return "Data processed successfully", nil
}
```

## API Reference

### Plugin Initialization

#### `sdk.Init(name, version, apiKey string) *Plugin`

Initializes a new plugin instance.

- `name`: Plugin name
- `version`: Plugin version
- `apiKey`: API key for authentication

### GraphQL Schema Registration

#### Individual Registration

```go
// Register a single query
plugin.RegisterQuery(name string, field GraphQLField, resolver ResolverFunc)

// Register a single mutation
plugin.RegisterMutation(name string, field GraphQLField, resolver ResolverFunc)
```

#### Batch Registration

```go
// Register multiple queries at once
queries := map[string]sdk.GraphQLField{
    "getUser": sdk.FieldWithArgs("User", "Get user by ID", map[string]interface{}{
        "id": sdk.NonNullArg("ID", "User ID"),
    }),
    "getUsers": sdk.ListField("User", "Get all users"),
}

resolvers := map[string]sdk.ResolverFunc{
    "getUser":  getUserResolver,
    "getUsers": getUsersResolver,
}

plugin.RegisterQueries(queries, resolvers)
```

### GraphQL Field Helpers

#### Basic Fields

```go
sdk.StringField("description")          // String
sdk.IntField("description")             // Int
sdk.BooleanField("description")         // Boolean
sdk.FloatField("description")           // Float
sdk.ListField("String", "description")  // [String]
sdk.NonNullField("String", "description") // String!
sdk.NonNullListField("String", "description") // [String!]!
```

#### Object Fields

```go
sdk.ObjectField("User object", map[string]interface{}{
    "id":    sdk.IntProperty("User ID"),
    "name":  sdk.StringProperty("User name"),
    "email": sdk.StringProperty("User email"),
})
```

#### Fields with Arguments

```go
sdk.FieldWithArgs("String", "Get user greeting", map[string]interface{}{
    "name": sdk.StringArg("User name"),
    "age":  sdk.IntArg("User age"),
    "user": sdk.ObjectArg("User data", map[string]interface{}{
        "id":   sdk.IntProperty("User ID"),
        "name": sdk.StringProperty("User name"),
    }),
})
```

### REST API Registration

#### Individual Registration

```go
endpoint := sdk.GETEndpoint("/users", "Get all users").
    WithResponseSchema(sdk.ObjectSchema(map[string]interface{}{
        "users": sdk.ArraySchema(sdk.ObjectSchema(map[string]interface{}{
            "id":   sdk.IntegerSchema("User ID"),
            "name": sdk.StringSchema("User name"),
        })),
    })).
    Build()

plugin.RegisterRESTAPI(endpoint, getUsersHandler)
```

#### Batch Registration

```go
endpoints := []sdk.RESTEndpoint{
    sdk.GETEndpoint("/health", "Health check").Build(),
    sdk.POSTEndpoint("/users", "Create user").Build(),
}

handlers := map[string]sdk.RESTHandlerFunc{
    "GET_/health": healthHandler,
    "POST_/users": createUserHandler,
}

plugin.RegisterRESTAPIs(endpoints, handlers)
```

#### Function Name Compatibility (v0.1.9+)

The SDK automatically handles both old and new REST API function naming conventions for compatibility with different engine versions:

**Engine Naming Conventions:**

- **Old Format**: `METHOD_/path` (e.g., `GET_/hello`, `POST_/users/:id`)
- **New Format**: `rest_method_path` (e.g., `rest_get_hello`, `rest_post_users_:id`)

**How It Works:**

```go
// You register REST APIs the same way regardless of engine version
plugin.RegisterRESTAPI(sdk.RESTEndpoint{
    Method: "GET",
    Path:   "/hello",
    Description: "Simple hello endpoint",
}, helloHandler)

// The SDK internally stores the handler with key: "GET_/hello"
// But can handle calls from engines using either:
// - Old engine: calls with "GET_/hello" ✅
// - New engine: calls with "rest_get_hello" ✅ (automatically converted)
```

**Automatic Conversion Examples:**

- `rest_get_hello` → `GET_/hello`
- `rest_post_users` → `POST_/users`
- `rest_get_users_:id` → `GET_/users/:id`
- `rest_put_users_:id_profile` → `PUT_/users/:id/profile`

**Debug Information:**
The SDK provides detailed logging when function name conversion occurs:

```
Plugin SDK: Trying to convert REST function name 'rest_get_hello' to old format 'GET_/hello'
Plugin SDK: Found REST handler using old format conversion
```

This ensures your plugins work with both older and newer versions of the Apito Engine without any code changes.

#### REST API Helper Functions (v0.1.10+)

The SDK provides specialized helper functions for parsing different types of REST API parameters:

**Path Parameters:**

```go
// Extract path parameters like /users/:id
userID := sdk.GetPathParam(args, "id", "default-id")

// Handles multiple formats: ":id", "id", "path_id"
func getUserHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    userID := sdk.GetPathParam(args, "id")
    if userID == "" {
        return map[string]interface{}{"error": "User ID required"}, nil
    }
    return map[string]interface{}{"user_id": userID}, nil
}
```

**Query Parameters:**

```go
// Extract query parameters like ?limit=10&active=true
limit := sdk.GetQueryParamInt(args, "limit", 20)
active := sdk.GetQueryParamBool(args, "active", true)
search := sdk.GetQueryParam(args, "search", "")

func listUsersHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    // Parse query parameters with defaults
    limit := sdk.GetQueryParamInt(args, "limit", 20)
    offset := sdk.GetQueryParamInt(args, "offset", 0)
    active := sdk.GetQueryParamBool(args, "active", true)

    return map[string]interface{}{
        "filters": map[string]interface{}{
            "limit": limit,
            "offset": offset,
            "active": active,
        },
    }, nil
}
```

**Body Parameters:**

```go
// Extract POST/PUT body parameters
name := sdk.GetBodyParam(args, "name")
age := sdk.GetBodyParamInt(args, "age", 0)
active := sdk.GetBodyParamBool(args, "active", true)
metadata := sdk.GetBodyParamObject(args, "metadata")

func createUserHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    // Extract and validate body parameters
    name := sdk.GetBodyParam(args, "name")
    email := sdk.GetBodyParam(args, "email")
    age := sdk.GetBodyParamInt(args, "age", 18)

    if name == "" || email == "" {
        return map[string]interface{}{
            "success": false,
            "error": "Name and email are required",
        }, nil
    }

    return map[string]interface{}{
        "success": true,
        "user": map[string]interface{}{
            "name": name,
            "email": email,
            "age": age,
        },
    }, nil
}
```

**Unified Parameter Parsing:**

```go
// Parse all parameters into categorized groups
func complexRESTHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    // Parse and categorize all parameters
    parsed := sdk.ParseRESTArgs(args)

    pathParams := parsed["path"].(map[string]interface{})
    queryParams := parsed["query"].(map[string]interface{})
    bodyParams := parsed["body"].(map[string]interface{})

    return map[string]interface{}{
        "path_params": pathParams,
        "query_params": queryParams,
        "body_params": bodyParams,
    }, nil
}
```

**Debug Logging:**

```go
func debugRESTHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    // Log all parameters in a structured way
    sdk.LogRESTArgs("debugRESTHandler", args)

    // Extract endpoint information
    endpointInfo := sdk.GetRESTEndpointInfo(args)

    return map[string]interface{}{
        "endpoint_info": endpointInfo,
        "message": "Check logs for detailed parameter breakdown",
    }, nil
}
```

**Available REST API Helpers:**

| Function                              | Purpose                     | Example                                 |
| ------------------------------------- | --------------------------- | --------------------------------------- |
| `GetPathParam(args, "id")`            | Extract path parameter      | `/users/:id` → `args[":id"]`            |
| `GetQueryParam(args, "search")`       | Extract query parameter     | `?search=john` → `args["query_search"]` |
| `GetQueryParamInt(args, "limit", 20)` | Extract integer query param | `?limit=10` with default 20             |
| `GetQueryParamBool(args, "active")`   | Extract boolean query param | `?active=true`                          |
| `GetBodyParam(args, "name")`          | Extract body parameter      | POST `{"name": "John"}`                 |
| `GetBodyParamObject(args, "user")`    | Extract object from body    | POST `{"user": {...}}`                  |
| `ParseRESTArgs(args)`                 | Categorize all parameters   | Returns `{path, query, body, raw}`      |
| `LogRESTArgs("handler", args)`        | Debug log all parameters    | Structured console output               |

### REST Endpoint Builders

```go
sdk.GETEndpoint(path, description)
sdk.POSTEndpoint(path, description)
sdk.PUTEndpoint(path, description)
sdk.DELETEEndpoint(path, description)
sdk.PATCHEndpoint(path, description)
```

### REST Schema Helpers

```go
sdk.ObjectSchema(properties)            // Object schema
sdk.ArraySchema(itemSchema)             // Array schema
sdk.StringSchema(description)           // String schema
sdk.IntegerSchema(description)          // Integer schema
sdk.BooleanSchema(description)          // Boolean schema
sdk.NumberSchema(description)           // Number schema
```

### Function Registration

#### Individual Registration

```go
plugin.RegisterFunction("processData", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    // Function logic here
    return "result", nil
})
```

#### Batch Registration

```go
functions := map[string]sdk.FunctionHandlerFunc{
    "processData":   processDataFunction,
    "validateData":  validateDataFunction,
    "transformData": transformDataFunction,
}

plugin.RegisterFunctions(functions)
```

### Function Signatures

```go
type ResolverFunc func(ctx context.Context, args map[string]interface{}) (interface{}, error)
type RESTHandlerFunc func(ctx context.Context, args map[string]interface{}) (interface{}, error)
type FunctionHandlerFunc func(ctx context.Context, args map[string]interface{}) (interface{}, error)
```

## Advanced Examples

### Complex GraphQL Query with Nested Objects

```go
plugin.RegisterQuery("processComplexData",
    sdk.FieldWithArgs("String", "Process complex input data", map[string]interface{}{
        "user": sdk.ObjectArg("Single user", map[string]interface{}{
            "id":     sdk.IntProperty("User ID"),
            "name":   sdk.StringProperty("User name"),
            "email":  sdk.StringProperty("User email"),
            "active": sdk.BooleanProperty("Is user active"),
        }),
        "tags": sdk.ListArg("String", "Array of tags"),
        "users": sdk.ListArg("Object", "Array of user objects"),
    }),
    processComplexDataResolver,
)
```

### REST API with Complex Schema

```go
endpoint := sdk.POSTEndpoint("/api/users", "Create new user").
    WithRequestSchema(sdk.ObjectSchema(map[string]interface{}{
        "user": sdk.ObjectSchema(map[string]interface{}{
            "name":     sdk.StringSchema("User name"),
            "email":    sdk.StringSchema("User email"),
            "age":      sdk.IntegerSchema("User age"),
            "metadata": sdk.ObjectSchema(map[string]interface{}{
                "department": sdk.StringSchema("User department"),
                "role":       sdk.StringSchema("User role"),
            }),
        }),
        "tags": sdk.ArraySchema(sdk.StringSchema("Tag name")),
    })).
    WithResponseSchema(sdk.ObjectSchema(map[string]interface{}{
        "success": sdk.BooleanSchema("Operation success"),
        "user_id": sdk.IntegerSchema("Created user ID"),
        "message": sdk.StringSchema("Response message"),
    })).
    Build()

plugin.RegisterRESTAPI(endpoint, createUserWithMetadataHandler)
```

## Error Handling

All resolver functions, REST handlers, and custom functions should return an error as the second return value:

```go
func myResolver(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    if someCondition {
        return nil, fmt.Errorf("validation failed: %s", reason)
    }

    result := processData(args)
    return result, nil
}
```

## Context Usage

The context parameter provides access to the request context and can be used for:

- Request timeouts and cancellation
- Passing request-scoped data
- Logging and tracing

```go
func myResolver(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    // Check for cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    // Extract request ID if available
    if requestID := ctx.Value("request_id"); requestID != nil {
        log.Printf("Processing request: %s", requestID)
    }

    return processWithContext(ctx, args), nil
}
```

## Building and Running

1. Create your plugin using the SDK
2. Build as a Go binary:
   ```bash
   go build -o my-plugin main.go
   ```
3. The Apito Engine will execute your plugin binary as a HashiCorp plugin

## Best Practices

1. **Use descriptive names** for GraphQL fields and REST endpoints
2. **Validate input data** in your resolvers and handlers
3. **Handle errors gracefully** and return meaningful error messages
4. **Use context** for request-scoped operations and cancellation
5. **Keep resolvers simple** and delegate complex logic to separate functions
6. **Test your plugins** thoroughly before deployment

## License

This SDK is part of the Apito Engine project.
