package sdk

import (
	"context"
	"fmt"
	"log"
	"os"

	hcplugin "github.com/hashicorp/go-plugin"
	"gitlab.com/apito.io/buffers/protobuff"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// Global plugin instance for resolver access
var currentPlugin *Plugin

// ResolverFunc is the function signature for GraphQL resolvers
type ResolverFunc func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// RESTHandlerFunc is the function signature for REST API handlers
type RESTHandlerFunc func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// FunctionHandlerFunc is the function signature for custom functions
type FunctionHandlerFunc func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// GraphQLField represents a GraphQL field definition
type GraphQLField struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Args        map[string]interface{} `json:"args,omitempty"`
	Resolve     string                 `json:"resolve"`
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
	}

	p.impl = &pluginImpl{plugin: p}

	// Set the global plugin instance for resolver access
	currentPlugin = p

	log.Printf("Plugin SDK: Initialized plugin '%s' version %s", name, version)
	return p
}

// RegisterQuery registers a GraphQL query
func (p *Plugin) RegisterQuery(name string, field GraphQLField, resolver ResolverFunc) {
	field.Resolve = name + "Resolver"
	p.queries[name] = field
	p.resolvers[name] = resolver
	log.Printf("Plugin SDK: Registered GraphQL query '%s'", name)
}

// RegisterMutation registers a GraphQL mutation
func (p *Plugin) RegisterMutation(name string, field GraphQLField, resolver ResolverFunc) {
	field.Resolve = name + "Resolver"
	p.mutations[name] = field
	p.resolvers[name] = resolver
	log.Printf("Plugin SDK: Registered GraphQL mutation '%s'", name)
}

// RegisterQueries registers multiple GraphQL queries at once
func (p *Plugin) RegisterQueries(queries map[string]GraphQLField, resolvers map[string]ResolverFunc) {
	for name, field := range queries {
		if resolver, exists := resolvers[name]; exists {
			p.RegisterQuery(name, field, resolver)
		} else {
			log.Printf("Plugin SDK: Warning - No resolver found for query '%s'", name)
		}
	}
}

// RegisterMutations registers multiple GraphQL mutations at once
func (p *Plugin) RegisterMutations(mutations map[string]GraphQLField, resolvers map[string]ResolverFunc) {
	for name, field := range mutations {
		if resolver, exists := resolvers[name]; exists {
			p.RegisterMutation(name, field, resolver)
		} else {
			log.Printf("Plugin SDK: Warning - No resolver found for mutation '%s'", name)
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
		} else {
			log.Printf("Plugin SDK: Warning - No handler found for REST API %s %s", endpoint.Method, endpoint.Path)
		}
	}
}

// RegisterFunction registers a custom function
func (p *Plugin) RegisterFunction(name string, function FunctionHandlerFunc) {
	p.functions[name] = function
	log.Printf("Plugin SDK: Registered function '%s'", name)
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

// Serve starts the plugin server
func (p *Plugin) Serve() {
	log.Printf("Plugin SDK: Starting plugin '%s' server...", p.name)

	handshakeConfig := hcplugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "APITO_PLUGIN",
		MagicCookieValue: "apito_plugin_magic_cookie_v1",
	}

	pluginMap := map[string]hcplugin.Plugin{
		"Plugin": &grpcPlugin{Impl: p.impl},
	}

	hcplugin.Serve(&hcplugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      hcplugin.DefaultGRPCServer,
	})

	log.Printf("Plugin SDK: Plugin '%s' shutting down...", p.name)
	os.Exit(0)
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
	log.Printf("Plugin SDK: Initializing plugin '%s'...", impl.plugin.name)

	// Log environment variables if needed
	for _, env := range req.EnvVars {
		log.Printf("Plugin SDK: Env %s=%s", env.Key, env.Value)
	}

	return &protobuff.InitResponse{
		Success: true,
		Message: fmt.Sprintf("Plugin '%s' initialized successfully", impl.plugin.name),
	}, nil
}

func (impl *pluginImpl) Migration(ctx context.Context, req *protobuff.MigrationRequest) (*protobuff.MigrationResponse, error) {
	log.Printf("Plugin SDK: Running migration for plugin '%s'...", impl.plugin.name)
	return &protobuff.MigrationResponse{
		Success: true,
		Message: fmt.Sprintf("No migration needed for plugin '%s'", impl.plugin.name),
	}, nil
}

func (impl *pluginImpl) SchemaRegister(ctx context.Context, req *protobuff.SchemaRegisterRequest) (*protobuff.SchemaRegisterResponse, error) {
	log.Printf("Plugin SDK: Registering GraphQL schema for plugin '%s'...", impl.plugin.name)

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

	queriesStruct, err := structpb.NewStruct(queriesMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create queries struct: %v", err)
	}

	mutationsStruct, err := structpb.NewStruct(mutationsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create mutations struct: %v", err)
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
		"type":        field.Type,
		"description": field.Description,
		"resolve":     field.Resolve,
	}

	// Serialize arguments if they exist
	if len(field.Args) > 0 {
		result["args"] = impl.serializeArgs(field.Args)
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

func (impl *pluginImpl) Execute(ctx context.Context, req *protobuff.ExecuteRequest) (*protobuff.ExecuteResponse, error) {
	log.Printf("Plugin SDK: Execute called - function: %s, type: %s", req.FunctionName, req.FunctionType)

	// Extract arguments from the request
	args := make(map[string]interface{})
	if req.Args != nil {
		args = req.Args.AsMap()
	}

	// Extract context data and merge it with arguments
	// This allows plugins to access sensitive data passed from the host
	if req.Context != nil {
		contextData := req.Context.AsMap()
		log.Printf("Plugin SDK: Context data received: %+v", contextData)

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
		if handler, exists := impl.plugin.restHandlers[req.FunctionName]; exists {
			result, err = handler(ctx, args)
		} else {
			return &protobuff.ExecuteResponse{
				Success: false,
				Message: fmt.Sprintf("Unknown REST handler: %s", req.FunctionName),
			}, nil
		}

	case "function":
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
	log.Printf("Plugin SDK: Debug called with stage: %s", req.Stage)

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
