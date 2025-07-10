package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	hcplugin "github.com/hashicorp/go-plugin"
	"gitlab.com/apito.io/buffers/protobuff"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// ========================================
// ERROR HANDLING WITH HTTP STATUS CODES
// ========================================

// CodedError represents an error with an HTTP status code
type CodedError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *CodedError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("HTTP %d: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

// ErrorWithCode creates a new error with HTTP status code
func ErrorWithCode(code int, message string, details ...string) error {
	err := &CodedError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// Common HTTP error constructors
func BadRequestError(message string, details ...string) error {
	return ErrorWithCode(400, message, details...)
}

func UnauthorizedError(message string, details ...string) error {
	return ErrorWithCode(401, message, details...)
}

func ForbiddenError(message string, details ...string) error {
	return ErrorWithCode(403, message, details...)
}

func NotFoundError(message string, details ...string) error {
	return ErrorWithCode(404, message, details...)
}

func InternalServerError(message string, details ...string) error {
	return ErrorWithCode(500, message, details...)
}

// GetErrorCode extracts HTTP status code from error, returns 500 for unknown errors
func GetErrorCode(err error) int {
	if codedErr, ok := err.(*CodedError); ok {
		return codedErr.Code
	}
	return 500 // Default to internal server error
}

// GetErrorMessage extracts error message, handling both coded and regular errors
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// Global plugin instance for resolver access
var currentPlugin *Plugin

// ResolverFunc is the function signature for GraphQL resolvers
type ResolverFunc func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// RESTHandlerFunc is the function signature for REST API handlers
type RESTHandlerFunc func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// FunctionHandlerFunc is the function signature for custom functions
type FunctionHandlerFunc func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// HealthCheckFunc is the function signature for custom health checks
type HealthCheckFunc func(ctx context.Context) (map[string]interface{}, error)

// GraphQLField represents a GraphQL field definition
type GraphQLField struct {
	Type        interface{}            `json:"type"` // Can be string or GraphQLTypeDefinition
	Description string                 `json:"description"`
	Args        map[string]interface{} `json:"args,omitempty"`
	Resolve     string                 `json:"resolve"`
}

// GraphQLTypeDefinition represents a complex GraphQL type
type GraphQLTypeDefinition struct {
	Kind       string                 `json:"kind"`       // "scalar", "object", "list", "non_null"
	Name       string                 `json:"name"`       // For scalar and object types
	OfType     *GraphQLTypeDefinition `json:"ofType"`     // For list and non_null types
	Fields     map[string]interface{} `json:"fields"`     // For object types
	ScalarType string                 `json:"scalarType"` // For scalar types: "String", "Int", "Boolean", "Float"
}

// RESTEndpoint represents a REST API endpoint definition
type RESTEndpoint struct {
	Method      string
	Path        string
	Description string
	Schema      map[string]interface{}
	Handler     string
}

// Plugin represents the SDK plugin instance
type Plugin struct {
	name         string
	version      string
	apiKey       string
	queries      map[string]GraphQLField
	mutations    map[string]GraphQLField
	restAPIs     []RESTEndpoint
	resolvers    map[string]ResolverFunc
	restHandlers map[string]RESTHandlerFunc
	functions    map[string]FunctionHandlerFunc
	healthChecks []HealthCheckFunc

	// Type registry for nested objects
	objectTypes map[string]ObjectTypeDefinition

	// Internal implementation
	impl *pluginImpl
}

// pluginImpl implements the protobuff.PluginServiceServer
type pluginImpl struct {
	protobuff.UnimplementedPluginServiceServer
	plugin *Plugin
}

// Init initializes a new plugin instance
func Init(name, version, apiKey string) *Plugin {
	p := &Plugin{
		name:         name,
		version:      version,
		apiKey:       apiKey,
		queries:      make(map[string]GraphQLField),
		mutations:    make(map[string]GraphQLField),
		restAPIs:     make([]RESTEndpoint, 0),
		resolvers:    make(map[string]ResolverFunc),
		restHandlers: make(map[string]RESTHandlerFunc),
		functions:    make(map[string]FunctionHandlerFunc),
		healthChecks: make([]HealthCheckFunc, 0),
		objectTypes:  make(map[string]ObjectTypeDefinition),
	}

	p.impl = &pluginImpl{plugin: p}

	// Register built-in health check function
	p.functions["health_check"] = func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return p.performHealthCheck(ctx)
	}

	// Set the global plugin instance for resolver access
	currentPlugin = p

	return p
}

// RegisterQuery registers a GraphQL query
func (p *Plugin) RegisterQuery(name string, field GraphQLField, resolver ResolverFunc) {
	field.Resolve = name + "Resolver"
	p.queries[name] = field
	p.resolvers[name] = resolver

}

// RegisterMutation registers a GraphQL mutation
func (p *Plugin) RegisterMutation(name string, field GraphQLField, resolver ResolverFunc) {
	field.Resolve = name + "Resolver"
	p.mutations[name] = field
	p.resolvers[name] = resolver

}

// RegisterQueries registers multiple GraphQL queries at once
func (p *Plugin) RegisterQueries(queries map[string]GraphQLField, resolvers map[string]ResolverFunc) {
	for name, field := range queries {
		if resolver, exists := resolvers[name]; exists {
			p.RegisterQuery(name, field, resolver)
		}
	}
}

// RegisterMutations registers multiple GraphQL mutations at once
func (p *Plugin) RegisterMutations(mutations map[string]GraphQLField, resolvers map[string]ResolverFunc) {
	for name, field := range mutations {
		if resolver, exists := resolvers[name]; exists {
			p.RegisterMutation(name, field, resolver)
		}
	}
}

// RegisterRESTAPI registers a REST API endpoint
func (p *Plugin) RegisterRESTAPI(endpoint RESTEndpoint, handler RESTHandlerFunc) {
	endpoint.Handler = endpoint.Method + "_" + endpoint.Path
	p.restAPIs = append(p.restAPIs, endpoint)
	p.restHandlers[endpoint.Handler] = handler
	log.Printf("Plugin SDK: Registered REST API %s %s", endpoint.Method, endpoint.Path)
}

// RegisterRESTAPIs registers multiple REST API endpoints at once
func (p *Plugin) RegisterRESTAPIs(endpoints []RESTEndpoint, handlers map[string]RESTHandlerFunc) {
	for _, endpoint := range endpoints {
		handlerKey := endpoint.Method + "_" + endpoint.Path
		if handler, exists := handlers[handlerKey]; exists {
			p.RegisterRESTAPI(endpoint, handler)
		}
	}
}

// RegisterFunction registers a custom function
func (p *Plugin) RegisterFunction(name string, function FunctionHandlerFunc) {
	p.functions[name] = function

}

// RegisterFunctions registers multiple custom functions at once
func (p *Plugin) RegisterFunctions(functions map[string]FunctionHandlerFunc) {
	for name, function := range functions {
		p.RegisterFunction(name, function)
	}
}

// GetQueryField returns the field definition for a query
func (p *Plugin) GetQueryField(name string) (GraphQLField, bool) {
	field, exists := p.queries[name]
	return field, exists
}

// GetMutationField returns the field definition for a mutation
func (p *Plugin) GetMutationField(name string) (GraphQLField, bool) {
	field, exists := p.mutations[name]
	return field, exists
}

// RegisterObjectType registers an object type definition for nested object support
func (p *Plugin) RegisterObjectType(objectType ObjectTypeDefinition) {
	p.objectTypes[objectType.TypeName] = objectType

}

// GetObjectType returns the object type definition for a given name
func (p *Plugin) GetObjectType(name string) (ObjectTypeDefinition, bool) {
	objectType, exists := p.objectTypes[name]
	return objectType, exists
}

// GetAllObjectTypes returns all registered object types
func (p *Plugin) GetAllObjectTypes() map[string]ObjectTypeDefinition {
	return p.objectTypes
}

// Serve starts the plugin server
func (p *Plugin) Serve() {

	handshakeConfig := hcplugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "APITO_PLUGIN",
		MagicCookieValue: "apito_plugin_magic_cookie_v1",
	}

	pluginMap := map[string]hcplugin.Plugin{
		"Plugin": &grpcPlugin{Impl: p.impl},
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   p.name,
		Output: os.Stderr,
		Level:  hclog.Error, // Only show errors
	})

	hcplugin.Serve(&hcplugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      hcplugin.DefaultGRPCServer,
		Logger:          logger,
	})
}

// grpcPlugin implements the hcplugin.GRPCPlugin interface
type grpcPlugin struct {
	hcplugin.Plugin
	Impl *pluginImpl
}

func (p *grpcPlugin) GRPCServer(broker *hcplugin.GRPCBroker, s *grpc.Server) error {
	protobuff.RegisterPluginServiceServer(s, p.Impl)
	return nil
}

func (p *grpcPlugin) GRPCClient(ctx context.Context, broker *hcplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return protobuff.NewPluginServiceClient(c), nil
}

// Implementation of protobuff.PluginServiceServer methods

func (impl *pluginImpl) Init(ctx context.Context, req *protobuff.InitRequest) (*protobuff.InitResponse, error) {
	// Set environment variables
	for _, env := range req.EnvVars {
		os.Setenv(env.Key, env.Value)
	}

	return &protobuff.InitResponse{
		Success: true,
		Message: fmt.Sprintf("Plugin '%s' initialized successfully", impl.plugin.name),
	}, nil
}

func (impl *pluginImpl) Migration(ctx context.Context, req *protobuff.MigrationRequest) (*protobuff.MigrationResponse, error) {
	return &protobuff.MigrationResponse{
		Success: true,
		Message: fmt.Sprintf("No migration needed for plugin '%s'", impl.plugin.name),
	}, nil
}

func (impl *pluginImpl) SchemaRegister(ctx context.Context, req *protobuff.SchemaRegisterRequest) (*protobuff.SchemaRegisterResponse, error) {

	// Convert queries to protobuf struct
	queriesMap := make(map[string]interface{})
	for name, field := range impl.plugin.queries {
		queriesMap[name] = impl.serializeGraphQLField(field)
	}

	// Convert mutations to protobuf struct
	mutationsMap := make(map[string]interface{})
	for name, field := range impl.plugin.mutations {
		mutationsMap[name] = impl.serializeGraphQLField(field)
	}

	// Convert object types to protobuf struct
	objectTypesMap := make(map[string]interface{})
	for name, objectType := range impl.plugin.objectTypes {
		serialized := impl.serializeObjectTypeDefinition(objectType)
		objectTypesMap[name] = serialized
		//log.Printf("[NESTED-OBJECT-DEBUG] [SDK] Serializing object type %s: %+v", name, serialized)
	}

	queriesStruct, err := structpb.NewStruct(queriesMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create queries struct: %v", err)
	}

	mutationsStruct, err := structpb.NewStruct(mutationsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create mutations struct: %v", err)
	}

	// For now, include object types in a custom field or extend the existing schema
	// We'll add object types as a special query field that the engine can recognize
	if len(objectTypesMap) > 0 {
		objectTypesField := map[string]interface{}{
			"type":        "String",
			"description": "Object type definitions for nested objects",
			"objectTypes": objectTypesMap,
		}
		queriesMap["__objectTypes"] = objectTypesField
		//log.Printf("[NESTED-OBJECT-DEBUG] [SDK] Adding __objectTypes field with %d types: %+v", len(objectTypesMap), objectTypesField)

		// Recreate the queries struct with object types included
		queriesStruct, err = structpb.NewStruct(queriesMap)
		if err != nil {
			return nil, fmt.Errorf("failed to recreate queries struct with object types: %v", err)
		}
	}

	schema := &protobuff.ThirdPartyGraphQLSchemas{
		Queries:   queriesStruct,
		Mutations: mutationsStruct,
	}

	log.Printf("Plugin SDK: GraphQL schema registered successfully for plugin '%s'", impl.plugin.name)
	return &protobuff.SchemaRegisterResponse{
		Schema: schema,
	}, nil
}

// serializeGraphQLField converts a GraphQLField to a protobuf-compatible map
func (impl *pluginImpl) serializeGraphQLField(field GraphQLField) map[string]interface{} {
	result := map[string]interface{}{
		"type":        impl.serializeType(field.Type),
		"description": field.Description,
		"resolve":     field.Resolve,
	}

	// Serialize arguments if they exist
	if len(field.Args) > 0 {
		result["args"] = impl.serializeArgs(field.Args)
	}

	return result
}

// serializeType converts a type (string or GraphQLTypeDefinition) to protobuf-compatible format
func (impl *pluginImpl) serializeType(fieldType interface{}) interface{} {
	switch t := fieldType.(type) {
	case string:
		// Legacy string type - convert to simple scalar type
		return map[string]interface{}{
			"kind":       "scalar",
			"scalarType": t,
			"name":       t,
		}
	case GraphQLTypeDefinition:
		// New structured type definition
		return impl.serializeTypeDefinition(t)
	default:
		// Fallback to string representation
		return fmt.Sprintf("%v", fieldType)
	}
}

// serializeTypeDefinition converts a GraphQLTypeDefinition to protobuf-compatible format
func (impl *pluginImpl) serializeTypeDefinition(typeDef GraphQLTypeDefinition) map[string]interface{} {
	result := map[string]interface{}{
		"kind": typeDef.Kind,
	}

	if typeDef.Name != "" {
		result["name"] = typeDef.Name
	}

	if typeDef.ScalarType != "" {
		result["scalarType"] = typeDef.ScalarType
	}

	if typeDef.OfType != nil {
		result["ofType"] = impl.serializeTypeDefinition(*typeDef.OfType)
	}

	if len(typeDef.Fields) > 0 {
		result["fields"] = impl.serializeArgs(typeDef.Fields)
	}

	return result
}

// serializeObjectTypeDefinition converts an ObjectTypeDefinition to protobuf-compatible format
func (impl *pluginImpl) serializeObjectTypeDefinition(objectType ObjectTypeDefinition) map[string]interface{} {
	// Convert ObjectFieldDef to the engine's expected format
	engineFields := make(map[string]interface{})
	for fieldName, fieldDef := range objectType.Fields {
		// Convert ObjectFieldDef to engine's GraphQL field format
		var fieldType map[string]interface{}

		// Start with the base type
		if impl.isScalarType(fieldDef.Type) {
			fieldType = map[string]interface{}{
				"kind":       "scalar",
				"name":       fieldDef.Type,
				"scalarType": fieldDef.Type,
			}
		} else {
			// For object types, create a reference
			fieldType = map[string]interface{}{
				"kind": "object",
				"name": fieldDef.Type,
			}
		}

		// Apply list wrapper if needed
		if fieldDef.List {
			if fieldDef.ListOfNonNull {
				fieldType = map[string]interface{}{
					"kind": "list",
					"ofType": map[string]interface{}{
						"kind":   "non_null",
						"ofType": fieldType,
					},
				}
			} else {
				fieldType = map[string]interface{}{
					"kind":   "list",
					"ofType": fieldType,
				}
			}
		}

		// Apply non-null wrapper if needed
		if !fieldDef.Nullable {
			fieldType = map[string]interface{}{
				"kind":   "non_null",
				"ofType": fieldType,
			}
		}

		engineFields[fieldName] = map[string]interface{}{
			"type":        fieldType,
			"description": fieldDef.Description,
		}
	}

	return map[string]interface{}{
		"kind":        "object",
		"name":        objectType.TypeName,
		"description": objectType.Description,
		"fields":      engineFields,
	}
}

// isScalarType checks if a type is a GraphQL scalar type
func (impl *pluginImpl) isScalarType(typeName string) bool {
	switch typeName {
	case "String", "Int", "Boolean", "Float", "ID":
		return true
	default:
		return false
	}
}

// serializeObjectFields converts ObjectFieldDef map to protobuf-compatible format
func (impl *pluginImpl) serializeObjectFields(fields map[string]ObjectFieldDef) map[string]interface{} {
	result := make(map[string]interface{})
	for fieldName, fieldDef := range fields {
		result[fieldName] = map[string]interface{}{
			"type":          fieldDef.Type,
			"description":   fieldDef.Description,
			"nullable":      fieldDef.Nullable,
			"list":          fieldDef.List,
			"listOfNonNull": fieldDef.ListOfNonNull,
		}
	}
	return result
}

// serializeArgs recursively serializes argument structures for protobuf compatibility
func (impl *pluginImpl) serializeArgs(args map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range args {
		result[key] = impl.serializeValue(value)
	}

	return result
}

// serializeValue recursively serializes any value to be protobuf-compatible
func (impl *pluginImpl) serializeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		// Recursively serialize nested maps
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = impl.serializeValue(val)
		}
		return result

	case []interface{}:
		// Recursively serialize arrays
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = impl.serializeValue(val)
		}
		return result

	case GraphQLField:
		// Convert GraphQLField to map
		return impl.serializeGraphQLField(v)

	case GraphQLTypeDefinition:
		// Convert GraphQLTypeDefinition to map
		return impl.serializeTypeDefinition(v)

	default:
		// Return primitive types as-is (string, int, bool, float, etc.)
		return value
	}
}

func (impl *pluginImpl) RESTApiRegister(ctx context.Context, req *protobuff.RESTApiRegisterRequest) (*protobuff.RESTApiRegisterResponse, error) {
	log.Printf("Plugin SDK: Registering REST APIs for plugin '%s'...", impl.plugin.name)

	apis := make([]*protobuff.ThirdPartyRESTApi, len(impl.plugin.restAPIs))
	for i, endpoint := range impl.plugin.restAPIs {
		schema, err := structpb.NewStruct(endpoint.Schema)
		if err != nil {
			return nil, fmt.Errorf("failed to create schema struct for %s %s: %v", endpoint.Method, endpoint.Path, err)
		}

		apis[i] = &protobuff.ThirdPartyRESTApi{
			Method:      endpoint.Method,
			Path:        endpoint.Path,
			Description: endpoint.Description,
			Schema:      schema,
		}
	}

	log.Printf("Plugin SDK: Registered %d REST API endpoints for plugin '%s'", len(apis), impl.plugin.name)
	return &protobuff.RESTApiRegisterResponse{
		Apis: apis,
	}, nil
}

func (impl *pluginImpl) GetVersion(ctx context.Context, req *protobuff.GetVersionRequest) (*protobuff.GetVersionResponse, error) {
	return &protobuff.GetVersionResponse{
		Version: impl.plugin.version,
	}, nil
}

// isComplexArrayData checks if the result contains complex data that structpb.NewStruct can't handle
func isComplexArrayData(data interface{}) bool {
	val := reflect.ValueOf(data)

	// Check if it's a slice/array
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		if val.Len() == 0 {
			return false // Empty arrays are fine
		}

		// Check the first element to see if it's complex
		firstElem := val.Index(0).Interface()
		switch firstElem.(type) {
		case map[string]interface{}:
			return true // Array of maps - complex
		case []interface{}:
			return true // Array of arrays - complex
		default:
			// Check if it's a struct type
			elemVal := reflect.ValueOf(firstElem)
			if elemVal.Kind() == reflect.Struct {
				return true // Array of structs - complex
			}
		}
	}

	// Check if it's a map containing complex arrays
	if val.Kind() == reflect.Map {
		for _, key := range val.MapKeys() {
			mapVal := val.MapIndex(key).Interface()
			if isComplexArrayData(mapVal) {
				return true
			}
		}
	}

	return false
}

// serializeComplexData serializes complex data as JSON bytes wrapped in anypb.Any
func serializeComplexData(data interface{}, functionName, functionType string) (*anypb.Any, error) {
	// Create the result map
	resultMap := map[string]interface{}{
		"data":          data,
		"function_name": functionName,
		"function_type": functionType,
		"serialization": "json_bytes", // Flag to indicate this is JSON serialized
	}

	// JSON serialize the entire result
	jsonBytes, err := json.Marshal(resultMap)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal complex data: %v", err)
	}

	// Pack JSON bytes as anypb.Any with type indication
	anyResult, err := anypb.New(&structpb.Value{
		Kind: &structpb.Value_StringValue{
			StringValue: string(jsonBytes),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create any result from JSON: %v", err)
	}

	return anyResult, nil
}

func (impl *pluginImpl) Execute(ctx context.Context, req *protobuff.ExecuteRequest) (*protobuff.ExecuteResponse, error) {

	// Extract arguments from the request
	args := make(map[string]interface{})
	if req.Args != nil {
		args = req.Args.AsMap()
	}

	// Extract context data and merge it with arguments
	// This allows plugins to access sensitive data passed from the host
	if req.Context != nil {
		contextData := req.Context.AsMap()

		// Add context data to args with a "context_" prefix to avoid conflicts
		for key, value := range contextData {
			contextKey := fmt.Sprintf("context_%s", key)
			args[contextKey] = value
		}

		// Also create a new context with the values for proper context propagation
		for key, value := range contextData {
			ctx = context.WithValue(ctx, key, value)
		}
	}

	var result interface{}
	var err error

	// Handle different function types
	switch req.FunctionType {
	case "graphql_query", "graphql_mutation":
		if resolver, exists := impl.plugin.resolvers[req.FunctionName]; exists {
			result, err = resolver(ctx, args)
		} else {
			return &protobuff.ExecuteResponse{
				Success: false,
				Message: fmt.Sprintf("Unknown GraphQL resolver: %s", req.FunctionName),
			}, nil
		}

	case "rest_api":
		// Try to find the handler using the function name directly first
		handler, exists := impl.plugin.restHandlers[req.FunctionName]

		// If not found, try to convert from new format (rest_method_path) to old format (METHOD_path)
		if !exists && strings.HasPrefix(req.FunctionName, "rest_") {
			// Convert from "rest_get_hello" to "GET_/hello"
			// Or from "rest_post_users_:id" to "POST_/users/:id"
			parts := strings.SplitN(req.FunctionName, "_", 3) // Split into ["rest", "method", "path"]
			if len(parts) >= 3 {
				method := strings.ToUpper(parts[1])
				pathParts := strings.Split(parts[2], "_")

				// Reconstruct the path with slashes
				var path strings.Builder
				path.WriteString("/")
				for i, part := range pathParts {
					if i > 0 {
						path.WriteString("/")
					}
					path.WriteString(part)
				}

				oldFormatKey := method + "_" + path.String()

				if h, found := impl.plugin.restHandlers[oldFormatKey]; found {
					handler = h
					exists = true
				}
			}
		}

		if exists {
			result, err = handler(ctx, args)
		} else {
			return &protobuff.ExecuteResponse{
				Success: false,
				Message: fmt.Sprintf("Unknown REST handler: %s", req.FunctionName),
			}, nil
		}

	case "function", "system":
		if function, exists := impl.plugin.functions[req.FunctionName]; exists {
			result, err = function(ctx, args)
		} else {
			return &protobuff.ExecuteResponse{
				Success: false,
				Message: fmt.Sprintf("Unknown function: %s", req.FunctionName),
			}, nil
		}

	default:
		return &protobuff.ExecuteResponse{
			Success: false,
			Message: fmt.Sprintf("Unsupported function type: %s", req.FunctionType),
		}, nil
	}

	if err != nil {
		return &protobuff.ExecuteResponse{
			Success: false,
			Message: fmt.Sprintf("Execution failed: %v", err),
		}, nil
	}

	// Convert result to protobuf Any
	if isComplexArrayData(result) {
		log.Printf("🎯 [SDK] Detected complex array data, using JSON bytes serialization")
		anyResult, err := serializeComplexData(result, req.FunctionName, req.FunctionType)
		if err != nil {
			return &protobuff.ExecuteResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to serialize complex data: %v", err),
			}, nil
		}

		return &protobuff.ExecuteResponse{
			Success: true,
			Message: "Execution completed successfully (complex data)",
			Result:  anyResult,
		}, nil
	}

	// Handle simple data with existing structpb approach
	resultMap := map[string]interface{}{
		"data":          result,
		"function_name": req.FunctionName,
		"function_type": req.FunctionType,
	}

	resultStruct, err := structpb.NewStruct(resultMap)
	if err != nil {
		return &protobuff.ExecuteResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create result struct: %v", err),
		}, nil
	}

	anyResult, err := anypb.New(resultStruct)
	if err != nil {
		return &protobuff.ExecuteResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create any result: %v", err),
		}, nil
	}

	return &protobuff.ExecuteResponse{
		Success: true,
		Message: "Execution completed successfully",
		Result:  anyResult,
	}, nil
}

func (impl *pluginImpl) Debug(ctx context.Context, req *protobuff.DebugRequest) (*protobuff.DebugResponse, error) {

	result := map[string]interface{}{
		"plugin":  impl.plugin.name,
		"version": impl.plugin.version,
		"stage":   req.Stage,
		"message": "Debug method called successfully",
	}

	resultStruct, err := structpb.NewStruct(result)
	if err != nil {
		return nil, fmt.Errorf("failed to create result struct: %v", err)
	}

	anyResult, err := anypb.New(resultStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to create any result: %v", err)
	}

	return &protobuff.DebugResponse{
		Result: anyResult,
	}, nil
}

// performHealthCheck performs a comprehensive health check of the plugin
func (p *Plugin) performHealthCheck(ctx context.Context) (interface{}, error) {
	startTime := time.Now()

	// Basic plugin health information
	healthInfo := map[string]interface{}{
		"status":    "healthy",
		"plugin":    p.name,
		"version":   p.version,
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(startTime).Milliseconds(),
	}

	// Runtime information
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	healthInfo["runtime"] = map[string]interface{}{
		"goroutines":       runtime.NumGoroutine(),
		"memory_allocated": memStats.Alloc,
		"memory_total":     memStats.TotalAlloc,
		"memory_sys":       memStats.Sys,
		"gc_cycles":        memStats.NumGC,
		"go_version":       runtime.Version(),
	}

	// Plugin registration statistics
	healthInfo["statistics"] = map[string]interface{}{
		"queries_registered":       len(p.queries),
		"mutations_registered":     len(p.mutations),
		"rest_apis_registered":     len(p.restAPIs),
		"functions_registered":     len(p.functions),
		"object_types_defined":     len(p.objectTypes),
		"health_checks_registered": len(p.healthChecks),
	}

	// Check if plugin can respond to basic operations
	healthInfo["capabilities"] = map[string]interface{}{
		"graphql_queries":   len(p.queries) > 0,
		"graphql_mutations": len(p.mutations) > 0,
		"rest_endpoints":    len(p.restAPIs) > 0,
		"custom_functions":  len(p.functions) > 0,
		"health_checks":     len(p.healthChecks) > 0,
	}

	// Environment information
	healthInfo["environment"] = map[string]interface{}{
		"pid":      os.Getpid(),
		"hostname": getHostname(),
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
	}

	// Run custom health checks
	customHealthResults := make(map[string]interface{})
	overallStatus := "healthy"

	for i, healthCheck := range p.healthChecks {
		checkName := fmt.Sprintf("custom_check_%d", i)
		checkResult, err := healthCheck(ctx)
		if err != nil {
			customHealthResults[checkName] = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
			overallStatus = "degraded"
		} else {
			customHealthResults[checkName] = checkResult
			// Check if the custom health check indicates an issue
			if status, ok := checkResult["status"].(string); ok && status != "healthy" {
				overallStatus = "degraded"
			}
		}
	}

	if len(p.healthChecks) > 0 {
		healthInfo["custom_health_checks"] = customHealthResults
	}

	// Update overall status based on custom checks
	healthInfo["status"] = overallStatus

	return healthInfo, nil
}

// getHostname safely gets the hostname
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// RegisterHealthCheck registers a custom health check function
func (p *Plugin) RegisterHealthCheck(healthCheck HealthCheckFunc) {
	p.healthChecks = append(p.healthChecks, healthCheck)
}

// RegisterHealthChecks registers multiple custom health check functions at once
func (p *Plugin) RegisterHealthChecks(healthChecks []HealthCheckFunc) {
	for _, healthCheck := range healthChecks {
		p.RegisterHealthCheck(healthCheck)
	}
}
