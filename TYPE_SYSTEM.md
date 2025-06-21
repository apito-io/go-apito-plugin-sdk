# GraphQL Type System

The Apito Plugin SDK now includes a comprehensive GraphQL type system that generates proper type structures compatible with the GraphQL engine, instead of simple string representations.

## Overview

Previously, the SDK used simple strings like `"String"` or `"[User]"` to represent types. This caused issues because the GraphQL engine expects structured type definitions that can be converted to actual GraphQL type objects like `graphql.NewObject()` and `graphql.NewList()`.

The new type system generates structured `GraphQLTypeDefinition` objects that the engine can properly interpret and convert to the appropriate GraphQL types.

## Type Structure

### GraphQLTypeDefinition

```go
type GraphQLTypeDefinition struct {
    Kind       string                 `json:"kind"`       // "scalar", "object", "list", "non_null"
    Name       string                 `json:"name"`       // For scalar and object types
    OfType     *GraphQLTypeDefinition `json:"ofType"`     // For list and non_null types
    Fields     map[string]interface{} `json:"fields"`     // For object types
    ScalarType string                 `json:"scalarType"` // For scalar types: "String", "Int", "Boolean", "Float"
}
```

### Type Kinds

1. **scalar** - Basic GraphQL scalar types (String, Int, Boolean, Float)
2. **object** - Complex object types with multiple fields
3. **list** - Array/list types containing other types
4. **non_null** - Non-nullable wrapper for other types

## Generated Type Structures

### Scalar Types

**Before:**

```json
{
  "type": "String"
}
```

**After:**

```json
{
  "type": {
    "kind": "scalar",
    "scalarType": "String",
    "name": "String"
  }
}
```

### Object Types

**Before:**

```json
{
  "type": "User"
}
```

**After:**

```json
{
  "type": {
    "kind": "object",
    "name": "User",
    "fields": {
      "id": {
        "type": {
          "kind": "non_null",
          "ofType": {
            "kind": "scalar",
            "scalarType": "String",
            "name": "String"
          }
        },
        "description": "User ID"
      },
      "name": {
        "type": {
          "kind": "scalar",
          "scalarType": "String",
          "name": "String"
        },
        "description": "User's full name"
      }
    }
  }
}
```

### List Types

**Before:**

```json
{
  "type": "[User]"
}
```

**After:**

```json
{
  "type": {
    "kind": "list",
    "ofType": {
      "kind": "object",
      "name": "User",
      "fields": {
        /* User fields */
      }
    }
  }
}
```

### Non-Null List Types

**Before:**

```json
{
  "type": "[User!]!"
}
```

**After:**

```json
{
  "type": {
    "kind": "non_null",
    "ofType": {
      "kind": "list",
      "ofType": {
        "kind": "non_null",
        "ofType": {
          "kind": "object",
          "name": "User",
          "fields": {
            /* User fields */
          }
        }
      }
    }
  }
}
```

## SDK Functions

All SDK field creation functions now generate proper type structures:

### Basic Fields

```go
// Generates scalar type structure
sdk.StringField("description")
sdk.IntField("description")
sdk.BooleanField("description")
sdk.FloatField("description")

// Generates list type structure
sdk.ListField("String", "description")

// Generates non-null type structure
sdk.NonNullField("String", "description")

// Generates non-null list type structure
sdk.NonNullListField("String", "description")
```

### Complex Object Fields

```go
// Define object type
userType := sdk.NewObjectType("User", "A user in the system").
    AddStringField("id", "User ID", false).  // Non-nullable
    AddStringField("name", "User's name", true). // Nullable
    Build()

// Generate object type structure
sdk.ComplexObjectField("description", userType)

// Generate list of objects type structure
sdk.ListOfObjectsField("description", userType)

// Generate non-null list of objects type structure
sdk.NonNullListOfObjectsField("description", userType)
```

## Engine Compatibility

The new type system generates structures that the GraphQL engine can convert to proper GraphQL types:

### Scalar Types → `graphql.String`, `graphql.Int`, etc.

```go
// SDK generates:
{
  "kind": "scalar",
  "scalarType": "String"
}

// Engine converts to:
graphql.String
```

### Object Types → `graphql.NewObject()`

```go
// SDK generates:
{
  "kind": "object",
  "name": "User",
  "fields": { /* field definitions */ }
}

// Engine converts to:
graphql.NewObject(graphql.ObjectConfig{
  Name: "User",
  Fields: /* converted fields */
})
```

### List Types → `graphql.NewList()`

```go
// SDK generates:
{
  "kind": "list",
  "ofType": { /* inner type */ }
}

// Engine converts to:
graphql.NewList(/* converted inner type */)
```

### Non-Null Types → `graphql.NewNonNull()`

```go
// SDK generates:
{
  "kind": "non_null",
  "ofType": { /* inner type */ }
}

// Engine converts to:
graphql.NewNonNull(/* converted inner type */)
```

## Backwards Compatibility

The SDK maintains backwards compatibility with string types:

```go
// Still works - converted to scalar type structure
plugin.RegisterQuery("test",
    sdk.FieldWithArgs("String", "description", args),
    resolver)

// Legacy string types are automatically converted:
// "String" → {"kind": "scalar", "scalarType": "String", "name": "String"}
```

## Migration Guide

### No Changes Required

Existing plugins using the SDK will continue to work without any code changes. The SDK automatically converts:

- Simple field types: `StringField()`, `IntField()`, etc.
- Legacy string types in `FieldWithArgs()`
- Complex object types from `ComplexObjectField()`

### New Features Available

With the new type system, you can now:

1. **Define Complex Objects**: Full object type definitions with nested fields
2. **Use Proper Lists**: Arrays of objects with correct type structures
3. **Control Nullability**: Precise control over nullable vs non-null fields
4. **Nested Objects**: Objects containing other objects
5. **Engine Compatibility**: Types that work correctly with the GraphQL engine

## Examples

### Simple Query

```go
plugin.RegisterQuery("hello",
    sdk.StringField("Returns a greeting"),
    helloResolver)
```

**Generated Type:**

```json
{
  "kind": "scalar",
  "scalarType": "String",
  "name": "String"
}
```

### Complex Object Query

```go
userType := sdk.NewObjectType("User", "A user").
    AddStringField("id", "User ID", false).
    AddStringField("email", "Email", true).
    Build()

plugin.RegisterQuery("getUser",
    sdk.ComplexObjectField("Get user", userType),
    getUserResolver)
```

**Generated Type:**

```json
{
  "kind": "object",
  "name": "User",
  "fields": {
    "id": {
      "type": {
        "kind": "non_null",
        "ofType": {
          "kind": "scalar",
          "scalarType": "String",
          "name": "String"
        }
      }
    },
    "email": {
      "type": {
        "kind": "scalar",
        "scalarType": "String",
        "name": "String"
      }
    }
  }
}
```

### List Query

```go
plugin.RegisterQuery("getUsers",
    sdk.ListOfObjectsField("Get users", userType),
    getUsersResolver)
```

**Generated Type:**

```json
{
  "kind": "list",
  "ofType": {
    "kind": "object",
    "name": "User",
    "fields": {
      /* User fields */
    }
  }
}
```

## Benefits

1. **Engine Compatibility**: Types work correctly with the GraphQL engine
2. **Type Safety**: Structured type definitions prevent errors
3. **Flexibility**: Support for complex nested types
4. **Maintainability**: Clear type structure for debugging
5. **Performance**: Efficient type conversion in the engine
6. **Standards Compliance**: Follows GraphQL type system standards
