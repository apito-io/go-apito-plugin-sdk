package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apito-io/go-apito-plugin-sdk"
)

func main() {
	// Initialize the plugin - much simpler than before!
	plugin := sdk.Init("hc-hello-world-plugin", "2.0.0-sdk", "your-api-key")

	// Register GraphQL queries using the simplified API
	plugin.RegisterQuery("helloWorldQuery",
		sdk.FieldWithArgs("String", "Returns a hello world message from the plugin", map[string]interface{}{
			"name": sdk.StringArg("Optional name to include in greeting"),
			"object": sdk.ObjectArg("An object with a name and age", map[string]interface{}{
				"name": sdk.StringProperty("Name of the object"),
				"age":  sdk.IntProperty("Age of the object"),
			}),
			"arrayofObjects": sdk.ListArg("Object", "Array of objects"),
		}),
		helloWorldResolver,
	)

	plugin.RegisterQuery("processComplexData",
		sdk.FieldWithArgs("String", "Processes complex input data including objects, arrays, and array of objects", map[string]interface{}{
			"user": sdk.ObjectArg("A single user object input", map[string]interface{}{
				"id":     sdk.IntProperty("User ID"),
				"name":   sdk.StringProperty("User name"),
				"email":  sdk.StringProperty("User email"),
				"age":    sdk.IntProperty("User age"),
				"active": sdk.BooleanProperty("Whether user is active"),
			}),
			"tags":          sdk.ListArg("String", "Array of string tags"),
			"numbers":       sdk.NonNullListField("Int", "Array of required integers"),
			"users":         sdk.ListArg("Object", "Array of required user objects"),
			"optionalUsers": sdk.ListArg("Object", "Array of optional user objects"),
		}),
		processComplexDataResolver,
	)

	// Register GraphQL mutations
	plugin.RegisterMutation("sayHelloMutation",
		sdk.FieldWithArgs("String", "Says hello with a custom message", map[string]interface{}{
			"message": sdk.NonNullArg("String", "The message to echo back with hello"),
		}),
		sayHelloResolver,
	)

	// Register REST API endpoints using the builder pattern
	plugin.RegisterRESTAPI(
		sdk.GETEndpoint("/plugin/hello", "Returns a simple hello world message").
			WithResponseSchema(sdk.ObjectSchema(map[string]interface{}{
				"message":   sdk.StringSchema("Hello world message"),
				"timestamp": sdk.StringSchema("Current timestamp"),
				"plugin":    sdk.StringSchema("Plugin name"),
			})).
			Build(),
		helloRESTHandler,
	)

	plugin.RegisterRESTAPI(
		sdk.POSTEndpoint("/plugin/hello/custom", "Returns a custom hello message").
			WithRequestSchema(sdk.ObjectSchema(map[string]interface{}{
				"name":    sdk.StringSchema("Name to include in greeting"),
				"message": sdk.StringSchema("Custom message"),
			})).
			WithResponseSchema(sdk.ObjectSchema(map[string]interface{}{
				"greeting": sdk.StringSchema("Custom greeting message"),
				"plugin":   sdk.StringSchema("Plugin name"),
			})).
			Build(),
		customHelloRESTHandler,
	)

	plugin.RegisterRESTAPI(
		sdk.GETEndpoint("/plugin/hello/status", "Returns plugin status and information").
			WithResponseSchema(sdk.ObjectSchema(map[string]interface{}{
				"status":   sdk.StringSchema("Plugin status"),
				"version":  sdk.StringSchema("Plugin version"),
				"features": sdk.ArraySchema(sdk.StringSchema("Feature name")),
			})).
			Build(),
		statusRESTHandler,
	)

	// Register custom functions
	plugin.RegisterFunction("customFunction", customFunction)

	// Start the plugin server - all the gRPC and HashiCorp plugin setup is handled internally
	plugin.Serve()
}

// GraphQL Resolvers - same logic, much cleaner setup!

func helloWorldResolver(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	var result strings.Builder
	result.WriteString("Hello World Plugin Response:\n")

	// Handle name parameter
	name := "World"
	if nameArg, ok := args["name"].(string); ok && nameArg != "" {
		name = nameArg
	}
	result.WriteString(fmt.Sprintf("Hello, %s!\n", name))

	// Handle object parameter
	if obj, exists := args["object"]; exists && obj != nil {
		if objMap, ok := obj.(map[string]interface{}); ok {
			result.WriteString("Object received: ")
			if objName, ok := objMap["name"].(string); ok {
				result.WriteString(fmt.Sprintf("name=%s ", objName))
			}
			if age, ok := objMap["age"].(float64); ok {
				result.WriteString(fmt.Sprintf("age=%d", int(age)))
			}
			result.WriteString("\n")
		}
	}

	// Handle arrayofObjects parameter
	if arrObjs, exists := args["arrayofObjects"]; exists && arrObjs != nil {
		if objSlice, ok := arrObjs.([]interface{}); ok {
			result.WriteString("Array of Objects received:\n")
			for i, obj := range objSlice {
				if objMap, ok := obj.(map[string]interface{}); ok {
					result.WriteString(fmt.Sprintf("  Object %d: ", i+1))
					if objName, ok := objMap["name"].(string); ok {
						result.WriteString(fmt.Sprintf("name=%s ", objName))
					}
					if age, ok := objMap["age"].(float64); ok {
						result.WriteString(fmt.Sprintf("age=%d", int(age)))
					}
					result.WriteString("\n")
				}
			}
		}
	}

	return result.String(), nil
}

func processComplexDataResolver(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	var result strings.Builder
	result.WriteString("Processing complex data:\n")

	// Process single user object
	if user, exists := args["user"]; exists && user != nil {
		if userMap, ok := user.(map[string]interface{}); ok {
			result.WriteString("User: ")
			if id, ok := userMap["id"].(float64); ok {
				result.WriteString(fmt.Sprintf("ID=%d ", int(id)))
			}
			if name, ok := userMap["name"].(string); ok {
				result.WriteString(fmt.Sprintf("Name=%s ", name))
			}
			if email, ok := userMap["email"].(string); ok {
				result.WriteString(fmt.Sprintf("Email=%s ", email))
			}
			if age, ok := userMap["age"].(float64); ok {
				result.WriteString(fmt.Sprintf("Age=%d ", int(age)))
			}
			if active, ok := userMap["active"].(bool); ok {
				result.WriteString(fmt.Sprintf("Active=%t", active))
			}
			result.WriteString("\n")
		}
	}

	// Process array of strings (tags)
	if tags, exists := args["tags"]; exists && tags != nil {
		if tagSlice, ok := tags.([]interface{}); ok {
			result.WriteString("Tags: ")
			for i, tag := range tagSlice {
				if tagStr, ok := tag.(string); ok {
					result.WriteString(tagStr)
					if i < len(tagSlice)-1 {
						result.WriteString(", ")
					}
				}
			}
			result.WriteString("\n")
		}
	}

	// Process other arrays as in the original...
	// (abbreviated for brevity)

	return result.String(), nil
}

func sayHelloResolver(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	message := "Hello!"
	if msgArg, ok := args["message"].(string); ok && msgArg != "" {
		message = msgArg
	}

	return fmt.Sprintf("Plugin says: %s (from hc-hello-world-plugin)", message), nil
}

// REST Handlers - much simpler than managing protobuf structs!

func helloRESTHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"message":   "Hello World from REST API!",
		"timestamp": time.Now().Format(time.RFC3339),
		"plugin":    "hc-hello-world-plugin",
	}, nil
}

func customHelloRESTHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	name := "World"
	message := "Hello"

	if nameArg, ok := args["name"].(string); ok && nameArg != "" {
		name = nameArg
	}
	if msgArg, ok := args["message"].(string); ok && msgArg != "" {
		message = msgArg
	}

	return map[string]interface{}{
		"greeting": fmt.Sprintf("%s, %s!", message, name),
		"plugin":   "hc-hello-world-plugin",
	}, nil
}

func statusRESTHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"status":  "running",
		"version": "2.0.0-sdk",
		"features": []string{
			"GraphQL Queries",
			"GraphQL Mutations",
			"REST APIs",
			"Custom Functions",
		},
	}, nil
}

// Custom Functions

func customFunction(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return "Custom function executed successfully", nil
}
