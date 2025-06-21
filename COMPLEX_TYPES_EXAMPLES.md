# Complex Types Examples

This guide shows how to use the enhanced complex type support in the Apito Plugin SDK. Instead of just returning simple strings, your GraphQL resolvers can now return complex objects, arrays, nested objects, and more.

## Basic Usage

### Simple Object Type

```go
// Define a User object type
userType := sdk.NewObjectType("User", "A user in the system").
    AddStringField("id", "User ID", false).
    AddStringField("name", "User's full name", true).
    AddStringField("email", "User's email address", true).
    AddBooleanField("active", "Whether the user is active", false).
    Build()

// Register a query that returns a User object
plugin.RegisterQuery("getUser",
    sdk.ComplexObjectFieldWithArgs("Get a user by ID", userType, map[string]interface{}{
        "id": sdk.StringArg("User ID to fetch"),
    }),
    getUserResolver)

// Resolver function
func getUserResolver(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    userID := sdk.GetStringArg(args, "id", "")

    // Return a map that matches the User object type structure
    return map[string]interface{}{
        "id":     userID,
        "name":   "John Doe",
        "email":  "john@example.com",
        "active": true,
    }, nil
}
```

### Array of Objects

```go
// Register a query that returns a list of users
plugin.RegisterQuery("getUsers",
    sdk.ListOfObjectsFieldWithArgs("Get a list of users", userType, map[string]interface{}{
        "limit":  sdk.IntArg("Maximum number of users to return"),
        "offset": sdk.IntArg("Number of users to skip"),
    }),
    getUsersResolver)

// Resolver function
func getUsersResolver(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    limit := sdk.GetIntArg(args, "limit", 10)
    offset := sdk.GetIntArg(args, "offset", 0)

    // Return an array of maps that match the User object type
    users := []interface{}{
        map[string]interface{}{
            "id":     "1",
            "name":   "John Doe",
            "email":  "john@example.com",
            "active": true,
        },
        map[string]interface{}{
            "id":     "2",
            "name":   "Jane Smith",
            "email":  "jane@example.com",
            "active": false,
        },
    }

    return users, nil
}
```

## Built-in Common Types

The SDK provides several built-in common object types:

### UserObjectType()

```go
userType := sdk.UserObjectType()
// Contains: id, name, email, username, active, createdAt, updatedAt
```

### PaginationInfoType()

```go
paginationType := sdk.PaginationInfoType()
// Contains: total, limit, offset, page, totalPages, hasNext, hasPrevious
```

## Field Types Reference

When building object types, you can use these field types:

- `AddStringField(name, description, nullable)` - String field
- `AddIntField(name, description, nullable)` - Integer field
- `AddBooleanField(name, description, nullable)` - Boolean field
- `AddFloatField(name, description, nullable)` - Float field
- `AddObjectField(name, description, objectType, nullable)` - Nested object field
- `AddListField(name, description, itemType, nullable, listOfNonNull)` - Generic list field
- `AddStringListField(name, description, nullable, listOfNonNull)` - List of strings
- `AddIntListField(name, description, nullable, listOfNonNull)` - List of integers
- `AddObjectListField(name, description, objectType, nullable, listOfNonNull)` - List of objects

## Tips

1. **Return Structure Must Match**: Your resolver function's return value must match the object type structure you defined.
2. **Use Nullable Fields**: Set `nullable: true` for optional fields that might not always have values.
3. **Reuse Object Types**: Define object types once and reuse them across multiple queries/mutations.
