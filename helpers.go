package sdk

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Field creates a basic GraphQL field
func Field(fieldType, description string) GraphQLField {
	return GraphQLField{
		Type:        createScalarType(fieldType),
		Description: description,
		Args:        make(map[string]interface{}),
	}
}

// FieldWithArgs creates a GraphQL field with arguments
func FieldWithArgs(fieldType, description string, args map[string]interface{}) GraphQLField {
	return GraphQLField{
		Type:        createScalarType(fieldType),
		Description: description,
		Args:        args,
	}
}

// StringField creates a String type GraphQL field
func StringField(description string) GraphQLField {
	return Field("String", description)
}

// IntField creates an Int type GraphQL field
func IntField(description string) GraphQLField {
	return Field("Int", description)
}

// BooleanField creates a Boolean type GraphQL field
func BooleanField(description string) GraphQLField {
	return Field("Boolean", description)
}

// FloatField creates a Float type GraphQL field
func FloatField(description string) GraphQLField {
	return Field("Float", description)
}

// ListField creates a list type GraphQL field
func ListField(itemType, description string) GraphQLField {
	return GraphQLField{
		Type:        createListType(createScalarType(itemType)),
		Description: description,
		Args:        make(map[string]interface{}),
	}
}

// NonNullField creates a non-null type GraphQL field
func NonNullField(fieldType, description string) GraphQLField {
	return GraphQLField{
		Type:        createNonNullType(createScalarType(fieldType)),
		Description: description,
		Args:        make(map[string]interface{}),
	}
}

// NonNullListField creates a non-null list type GraphQL field
func NonNullListField(itemType, description string) GraphQLField {
	return GraphQLField{
		Type:        createNonNullType(createListType(createNonNullType(createScalarType(itemType)))),
		Description: description,
		Args:        make(map[string]interface{}),
	}
}

// =====================================================
// ENHANCED COMPLEX TYPE SUPPORT
// =====================================================

// ObjectTypeDefinition represents a complex object type with multiple fields
type ObjectTypeDefinition struct {
	TypeName    string                    `json:"typeName"`
	Description string                    `json:"description"`
	Fields      map[string]ObjectFieldDef `json:"fields"`
}

// ObjectFieldDef represents a field within an object type
type ObjectFieldDef struct {
	Type          string `json:"type"`
	Description   string `json:"description"`
	Nullable      bool   `json:"nullable"`
	List          bool   `json:"list"`
	ListOfNonNull bool   `json:"listOfNonNull"`
}

// ComplexObjectField creates a GraphQL field that returns a complex object type
func ComplexObjectField(description string, objectDef ObjectTypeDefinition) GraphQLField {
	// Convert ObjectTypeDefinition to GraphQLTypeDefinition
	objectFields := convertObjectFieldsToGraphQLFields(objectDef.Fields)

	return GraphQLField{
		Type:        createObjectType(objectDef.TypeName, objectFields),
		Description: description,
		Args: map[string]interface{}{
			"objectType": map[string]interface{}{
				"typeName":    objectDef.TypeName,
				"description": objectDef.Description,
				"fields":      serializeObjectFields(objectDef.Fields),
			},
		},
	}
}

// ComplexObjectFieldWithArgs creates a GraphQL field with args that returns a complex object type
func ComplexObjectFieldWithArgs(description string, objectDef ObjectTypeDefinition, args map[string]interface{}) GraphQLField {
	field := ComplexObjectField(description, objectDef)

	// Merge the arguments with the object type definition
	for key, value := range args {
		field.Args[key] = value
	}

	return field
}

// ListOfObjectsField creates a GraphQL field that returns a list of complex objects
func ListOfObjectsField(description string, objectDef ObjectTypeDefinition) GraphQLField {
	// Convert ObjectTypeDefinition to GraphQLTypeDefinition
	objectFields := convertObjectFieldsToGraphQLFields(objectDef.Fields)
	objectType := createObjectType(objectDef.TypeName, objectFields)

	return GraphQLField{
		Type:        createListType(objectType),
		Description: description,
		Args: map[string]interface{}{
			"objectType": map[string]interface{}{
				"typeName":    objectDef.TypeName,
				"description": objectDef.Description,
				"fields":      serializeObjectFields(objectDef.Fields),
			},
		},
	}
}

// ListOfObjectsFieldWithArgs creates a GraphQL field with args that returns a list of complex objects
func ListOfObjectsFieldWithArgs(description string, objectDef ObjectTypeDefinition, args map[string]interface{}) GraphQLField {
	field := ListOfObjectsField(description, objectDef)

	// Merge the arguments with the object type definition
	for key, value := range args {
		field.Args[key] = value
	}

	return field
}

// NonNullListOfObjectsField creates a non-null list of objects field
func NonNullListOfObjectsField(description string, objectDef ObjectTypeDefinition) GraphQLField {
	// Convert ObjectTypeDefinition to GraphQLTypeDefinition
	objectFields := convertObjectFieldsToGraphQLFields(objectDef.Fields)
	objectType := createObjectType(objectDef.TypeName, objectFields)

	return GraphQLField{
		Type:        createNonNullType(createListType(createNonNullType(objectType))),
		Description: description,
		Args: map[string]interface{}{
			"objectType": map[string]interface{}{
				"typeName":    objectDef.TypeName,
				"description": objectDef.Description,
				"fields":      serializeObjectFields(objectDef.Fields),
			},
		},
	}
}

// serializeObjectFields converts ObjectFieldDef map to serializable format
func serializeObjectFields(fields map[string]ObjectFieldDef) map[string]interface{} {
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

// =====================================================
// OBJECT TYPE DEFINITION BUILDERS
// =====================================================

// NewObjectType creates a new object type definition
func NewObjectType(typeName, description string) *ObjectTypeBuilder {
	return &ObjectTypeBuilder{
		def: ObjectTypeDefinition{
			TypeName:    typeName,
			Description: description,
			Fields:      make(map[string]ObjectFieldDef),
		},
	}
}

// ObjectTypeBuilder helps build complex object type definitions
type ObjectTypeBuilder struct {
	def ObjectTypeDefinition
}

// AddStringField adds a string field to the object type
func (b *ObjectTypeBuilder) AddStringField(name, description string, nullable bool) *ObjectTypeBuilder {
	b.def.Fields[name] = ObjectFieldDef{
		Type:          "String",
		Description:   description,
		Nullable:      nullable,
		List:          false,
		ListOfNonNull: false,
	}
	return b
}

// AddIntField adds an integer field to the object type
func (b *ObjectTypeBuilder) AddIntField(name, description string, nullable bool) *ObjectTypeBuilder {
	b.def.Fields[name] = ObjectFieldDef{
		Type:          "Int",
		Description:   description,
		Nullable:      nullable,
		List:          false,
		ListOfNonNull: false,
	}
	return b
}

// AddBooleanField adds a boolean field to the object type
func (b *ObjectTypeBuilder) AddBooleanField(name, description string, nullable bool) *ObjectTypeBuilder {
	b.def.Fields[name] = ObjectFieldDef{
		Type:          "Boolean",
		Description:   description,
		Nullable:      nullable,
		List:          false,
		ListOfNonNull: false,
	}
	return b
}

// AddFloatField adds a float field to the object type
func (b *ObjectTypeBuilder) AddFloatField(name, description string, nullable bool) *ObjectTypeBuilder {
	b.def.Fields[name] = ObjectFieldDef{
		Type:          "Float",
		Description:   description,
		Nullable:      nullable,
		List:          false,
		ListOfNonNull: false,
	}
	return b
}

// AddObjectField adds a nested object field to the object type
func (b *ObjectTypeBuilder) AddObjectField(name, description, objectType string, nullable bool) *ObjectTypeBuilder {
	b.def.Fields[name] = ObjectFieldDef{
		Type:          objectType,
		Description:   description,
		Nullable:      nullable,
		List:          false,
		ListOfNonNull: false,
	}
	return b
}

// AddListField adds an array/list field to the object type
func (b *ObjectTypeBuilder) AddListField(name, description, itemType string, nullable, listOfNonNull bool) *ObjectTypeBuilder {
	b.def.Fields[name] = ObjectFieldDef{
		Type:          itemType,
		Description:   description,
		Nullable:      nullable,
		List:          true,
		ListOfNonNull: listOfNonNull,
	}
	return b
}

// AddStringListField adds a list of strings field
func (b *ObjectTypeBuilder) AddStringListField(name, description string, nullable, listOfNonNull bool) *ObjectTypeBuilder {
	return b.AddListField(name, description, "String", nullable, listOfNonNull)
}

// AddIntListField adds a list of integers field
func (b *ObjectTypeBuilder) AddIntListField(name, description string, nullable, listOfNonNull bool) *ObjectTypeBuilder {
	return b.AddListField(name, description, "Int", nullable, listOfNonNull)
}

// AddObjectListField adds a list of objects field
func (b *ObjectTypeBuilder) AddObjectListField(name, description, objectType string, nullable, listOfNonNull bool) *ObjectTypeBuilder {
	return b.AddListField(name, description, objectType, nullable, listOfNonNull)
}

// Build returns the completed object type definition
func (b *ObjectTypeBuilder) Build() ObjectTypeDefinition {
	return b.def
}

// =====================================================
// COMMON COMPLEX TYPE DEFINITIONS
// =====================================================

// UserObjectType creates a standard User object type
func UserObjectType() ObjectTypeDefinition {
	return NewObjectType("User", "A user in the system").
		AddStringField("id", "User ID", false).
		AddStringField("name", "User's full name", true).
		AddStringField("email", "User's email address", true).
		AddStringField("username", "User's username", true).
		AddBooleanField("active", "Whether the user is active", false).
		AddStringField("createdAt", "When the user was created", true).
		AddStringField("updatedAt", "When the user was last updated", true).
		Build()
}

// PaginationInfoType creates a standard pagination info object type
func PaginationInfoType() ObjectTypeDefinition {
	return NewObjectType("PaginationInfo", "Information about pagination").
		AddIntField("total", "Total number of items", false).
		AddIntField("limit", "Number of items per page", false).
		AddIntField("offset", "Current offset", false).
		AddIntField("page", "Current page number", false).
		AddIntField("totalPages", "Total number of pages", false).
		AddBooleanField("hasNext", "Whether there is a next page", false).
		AddBooleanField("hasPrevious", "Whether there is a previous page", false).
		Build()
}

// ErrorObjectType creates a standard error object type
func ErrorObjectType() ObjectTypeDefinition {
	return NewObjectType("Error", "An error object").
		AddStringField("code", "Error code", false).
		AddStringField("message", "Error message", false).
		AddStringField("field", "Field that caused the error", true).
		AddStringListField("details", "Additional error details", true, false).
		Build()
}

// ResponseWrapperType creates a generic response wrapper type
func ResponseWrapperType(dataType string) ObjectTypeDefinition {
	return NewObjectType("Response", "A generic response wrapper").
		AddBooleanField("success", "Whether the operation was successful", false).
		AddStringField("message", "Response message", true).
		AddObjectField("data", "The response data", dataType, true).
		AddObjectListField("errors", "List of errors if any", "Error", true, false).
		Build()
}

// PaginatedResponseType creates a paginated response type
func PaginatedResponseType(itemType string) ObjectTypeDefinition {
	return NewObjectType("PaginatedResponse", "A paginated response").
		AddObjectListField("items", "List of items", itemType, false, false).
		AddObjectField("pagination", "Pagination information", "PaginationInfo", false).
		AddBooleanField("success", "Whether the operation was successful", false).
		AddStringField("message", "Response message", true).
		Build()
}

// =====================================================
// BACKWARD COMPATIBILITY - OLD OBJECTFIELD
// =====================================================

// ObjectField creates an object type GraphQL field with properties (legacy)
// Deprecated: Use ComplexObjectField instead for better type safety
func ObjectField(description string, properties map[string]interface{}) GraphQLField {
	field := Field("Object", description)
	field.Args = map[string]interface{}{
		"properties": properties,
	}
	return field
}

// Arg creates a GraphQL argument definition
func Arg(argType, description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        argType,
		"description": description,
	}
}

// StringArg creates a String type argument
func StringArg(description string) map[string]interface{} {
	return Arg("String", description)
}

// IntArg creates an Int type argument
func IntArg(description string) map[string]interface{} {
	return Arg("Int", description)
}

// BooleanArg creates a Boolean type argument
func BooleanArg(description string) map[string]interface{} {
	return Arg("Boolean", description)
}

// FloatArg creates a Float type argument
func FloatArg(description string) map[string]interface{} {
	return Arg("Float", description)
}

// NonNullArg creates a non-null type argument
func NonNullArg(argType, description string) map[string]interface{} {
	return Arg(argType+"!", description)
}

// ObjectArg creates an Object type argument with properties
func ObjectArg(description string, properties map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":        "Object",
		"description": description,
		"properties":  properties,
	}
}

// ListArg creates a list type argument
func ListArg(itemType, description string) map[string]interface{} {
	return Arg("["+itemType+"]", description)
}

// Property creates a property definition for object types
func Property(propType, description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        propType,
		"description": description,
	}
}

// StringProperty creates a String type property
func StringProperty(description string) map[string]interface{} {
	return Property("String", description)
}

// IntProperty creates an Int type property
func IntProperty(description string) map[string]interface{} {
	return Property("Int", description)
}

// BooleanProperty creates a Boolean type property
func BooleanProperty(description string) map[string]interface{} {
	return Property("Boolean", description)
}

// FloatProperty creates a Float type property
func FloatProperty(description string) map[string]interface{} {
	return Property("Float", description)
}

// RESTEndpointBuilder helps build REST endpoint definitions
type RESTEndpointBuilder struct {
	endpoint RESTEndpoint
}

// NewRESTEndpoint creates a new REST endpoint builder
func NewRESTEndpoint(method, path, description string) *RESTEndpointBuilder {
	return &RESTEndpointBuilder{
		endpoint: RESTEndpoint{
			Method:      method,
			Path:        path,
			Description: description,
			Schema:      make(map[string]interface{}),
		},
	}
}

// WithRequestSchema adds request schema to the REST endpoint
func (b *RESTEndpointBuilder) WithRequestSchema(schema map[string]interface{}) *RESTEndpointBuilder {
	b.endpoint.Schema["request"] = schema
	return b
}

// WithResponseSchema adds response schema to the REST endpoint
func (b *RESTEndpointBuilder) WithResponseSchema(schema map[string]interface{}) *RESTEndpointBuilder {
	b.endpoint.Schema["response"] = schema
	return b
}

// Build returns the constructed REST endpoint
func (b *RESTEndpointBuilder) Build() RESTEndpoint {
	return b.endpoint
}

// Common REST endpoint helpers

// GETEndpoint creates a GET REST endpoint
func GETEndpoint(path, description string) *RESTEndpointBuilder {
	return NewRESTEndpoint("GET", path, description)
}

// POSTEndpoint creates a POST REST endpoint
func POSTEndpoint(path, description string) *RESTEndpointBuilder {
	return NewRESTEndpoint("POST", path, description)
}

// PUTEndpoint creates a PUT REST endpoint
func PUTEndpoint(path, description string) *RESTEndpointBuilder {
	return NewRESTEndpoint("PUT", path, description)
}

// DELETEEndpoint creates a DELETE REST endpoint
func DELETEEndpoint(path, description string) *RESTEndpointBuilder {
	return NewRESTEndpoint("DELETE", path, description)
}

// PATCHEndpoint creates a PATCH REST endpoint
func PATCHEndpoint(path, description string) *RESTEndpointBuilder {
	return NewRESTEndpoint("PATCH", path, description)
}

// Schema helpers for REST APIs

// ObjectSchema creates an object schema for REST API
func ObjectSchema(properties map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}
}

// ArraySchema creates an array schema for REST API
func ArraySchema(itemSchema map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":  "array",
		"items": itemSchema,
	}
}

// StringSchema creates a string schema for REST API
func StringSchema(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": description,
	}
}

// IntegerSchema creates an integer schema for REST API
func IntegerSchema(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "integer",
		"description": description,
	}
}

// BooleanSchema creates a boolean schema for REST API
func BooleanSchema(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "boolean",
		"description": description,
	}
}

// NumberSchema creates a number schema for REST API
func NumberSchema(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "number",
		"description": description,
	}
}

// ArgParser helps parse and convert GraphQL arguments based on field definitions
type ArgParser struct {
	fieldDef GraphQLField
}

// NewArgParser creates a new argument parser for a GraphQL field
func NewArgParser(field GraphQLField) *ArgParser {
	return &ArgParser{fieldDef: field}
}

// ParseArgs converts raw GraphQL arguments to properly typed Go values based on field definition
func (p *ArgParser) ParseArgs(rawArgs map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for argName, argDef := range p.fieldDef.Args {
		if rawValue, exists := rawArgs[argName]; exists && rawValue != nil {
			result[argName] = p.parseValue(rawValue, argDef)
		}
	}

	return result
}

// parseValue converts a raw value based on argument definition
func (p *ArgParser) parseValue(rawValue interface{}, argDef interface{}) interface{} {
	// Handle argument definition as map
	if argDefMap, ok := argDef.(map[string]interface{}); ok {
		argType, _ := argDefMap["type"].(string)

		switch {
		case argType == "Object":
			return p.parseObject(rawValue, argDefMap)
		case argType == "[Object]" || argType == "[Object!]":
			return p.parseObjectArray(rawValue, argDefMap)
		case argType == "[String]" || argType == "[String!]":
			return p.parseStringArray(rawValue)
		case argType == "[Int]" || argType == "[Int!]":
			return p.parseIntArray(rawValue)
		case argType == "[Boolean]" || argType == "[Boolean!]":
			return p.parseBooleanArray(rawValue)
		case argType == "String" || argType == "String!":
			return p.parseString(rawValue)
		case argType == "Int" || argType == "Int!":
			return p.parseInt(rawValue)
		case argType == "Boolean" || argType == "Boolean!":
			return p.parseBoolean(rawValue)
		case argType == "Float" || argType == "Float!":
			return p.parseFloat(rawValue)
		default:
			return rawValue
		}
	}

	return rawValue
}

// parseObject converts raw object data to structured map
func (p *ArgParser) parseObject(rawValue interface{}, argDef map[string]interface{}) map[string]interface{} {
	if objMap, ok := rawValue.(map[string]interface{}); ok {
		result := make(map[string]interface{})

		// If properties are defined, validate and convert them
		if properties, exists := argDef["properties"]; exists {
			if propMap, ok := properties.(map[string]interface{}); ok {
				for propName, propValue := range objMap {
					if propDef, exists := propMap[propName]; exists {
						result[propName] = p.parseValue(propValue, propDef)
					} else {
						result[propName] = propValue // Keep unknown properties as-is
					}
				}
				return result
			}
		}

		return objMap
	}

	if objMap, ok := rawValue.(map[string]interface{}); ok {
		return objMap
	}
	return make(map[string]interface{})
}

// parseObjectArray converts raw array of objects
func (p *ArgParser) parseObjectArray(rawValue interface{}, argDef map[string]interface{}) []interface{} {
	if arr, ok := rawValue.([]interface{}); ok {
		result := make([]interface{}, len(arr))
		for i, item := range arr {
			result[i] = p.parseObject(item, argDef)
		}
		return result
	}

	return []interface{}{rawValue}
}

// parseStringArray converts raw array to string array
func (p *ArgParser) parseStringArray(rawValue interface{}) []string {
	if arr, ok := rawValue.([]interface{}); ok {
		result := make([]string, len(arr))
		for i, item := range arr {
			result[i] = p.parseString(item)
		}
		return result
	}

	return []string{p.parseString(rawValue)}
}

// parseIntArray converts raw array to int array
func (p *ArgParser) parseIntArray(rawValue interface{}) []int {
	if arr, ok := rawValue.([]interface{}); ok {
		result := make([]int, len(arr))
		for i, item := range arr {
			result[i] = p.parseInt(item)
		}
		return result
	}

	return []int{p.parseInt(rawValue)}
}

// parseBooleanArray converts raw array to boolean array
func (p *ArgParser) parseBooleanArray(rawValue interface{}) []bool {
	if arr, ok := rawValue.([]interface{}); ok {
		result := make([]bool, len(arr))
		for i, item := range arr {
			result[i] = p.parseBoolean(item)
		}
		return result
	}

	return []bool{p.parseBoolean(rawValue)}
}

// parseString safely converts value to string
func (p *ArgParser) parseString(rawValue interface{}) string {
	if str, ok := rawValue.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", rawValue)
}

// parseInt safely converts value to int
func (p *ArgParser) parseInt(rawValue interface{}) int {
	switch v := rawValue.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// parseBoolean safely converts value to boolean
func (p *ArgParser) parseBoolean(rawValue interface{}) bool {
	if b, ok := rawValue.(bool); ok {
		return b
	}
	return false
}

// parseFloat safely converts value to float64
func (p *ArgParser) parseFloat(rawValue interface{}) float64 {
	switch v := rawValue.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.0
}

// Convenience functions for common parsing scenarios

// ParseGraphQLArgs is a convenience function to parse arguments for a specific field
func ParseGraphQLArgs(field GraphQLField, rawArgs map[string]interface{}) map[string]interface{} {
	parser := NewArgParser(field)
	return parser.ParseArgs(rawArgs)
}

// GetStringArg safely extracts a string argument
func GetStringArg(args map[string]interface{}, name string, defaultValue ...string) string {
	if val, exists := args[name]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// GetIntArg safely extracts an int argument
func GetIntArg(args map[string]interface{}, name string, defaultValue ...int) int {
	if val, exists := args[name]; exists {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetBoolArg safely extracts a boolean argument
func GetBoolArg(args map[string]interface{}, name string, defaultValue ...bool) bool {
	if val, exists := args[name]; exists {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// GetObjectArg safely extracts an object argument
func GetObjectArg(args map[string]interface{}, name string) map[string]interface{} {
	if val, exists := args[name]; exists {
		if obj, ok := val.(map[string]interface{}); ok {
			return obj
		}
	}
	return make(map[string]interface{})
}

// GetArrayArg safely extracts an array argument
func GetArrayArg(args map[string]interface{}, name string) []interface{} {
	if val, exists := args[name]; exists {
		if arr, ok := val.([]interface{}); ok {
			return arr
		}
	}
	return []interface{}{}
}

// ParseArgsForResolver automatically parses arguments for the current resolver based on field definition
// This is the main function that plugins should use in their resolvers
func ParseArgsForResolver(resolverName string, rawArgs map[string]interface{}) map[string]interface{} {
	if currentPlugin == nil {
		log.Printf("SDK Warning: No current plugin instance available for argument parsing")
		return rawArgs
	}

	// Try to find the field definition in queries first
	if field, exists := currentPlugin.GetQueryField(resolverName); exists {
		return ParseGraphQLArgs(field, rawArgs)
	}

	// Then try mutations
	if field, exists := currentPlugin.GetMutationField(resolverName); exists {
		return ParseGraphQLArgs(field, rawArgs)
	}

	log.Printf("SDK Warning: No field definition found for resolver '%s', returning raw args", resolverName)
	return rawArgs
}

// Context data access helpers - these help plugins access sensitive data passed from the host

// GetContextString safely extracts a string value from context data in args
func GetContextString(args map[string]interface{}, key string, defaultValue ...string) string {
	contextKey := fmt.Sprintf("context_%s", key)
	if val, exists := args[contextKey]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// GetContextFromContext safely extracts a string value directly from context
func GetContextFromContext(ctx context.Context, key string, defaultValue ...string) string {
	if val := ctx.Value(key); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// GetPluginID extracts the plugin ID from context data
func GetPluginID(args map[string]interface{}) string {
	return GetContextString(args, "plugin_id")
}

// GetPluginIDFromContext extracts the plugin ID directly from context
func GetPluginIDFromContext(ctx context.Context) string {
	return GetContextFromContext(ctx, "plugin_id")
}

// GetProjectID extracts the project ID from context data
func GetProjectID(args map[string]interface{}) string {
	return GetContextString(args, "project_id")
}

// GetProjectIDFromContext extracts the project ID directly from context
func GetProjectIDFromContext(ctx context.Context) string {
	return GetContextFromContext(ctx, "project_id")
}

// GetUserID extracts the user ID from context data
func GetUserID(args map[string]interface{}) string {
	return GetContextString(args, "user_id")
}

// GetUserIDFromContext extracts the user ID directly from context
func GetUserIDFromContext(ctx context.Context) string {
	return GetContextFromContext(ctx, "user_id")
}

// GetTenantID extracts the tenant ID from context data
func GetTenantID(args map[string]interface{}) string {
	return GetContextString(args, "tenant_id")
}

// GetTenantIDFromContext extracts the tenant ID directly from context
func GetTenantIDFromContext(ctx context.Context) string {
	return GetContextFromContext(ctx, "tenant_id")
}

// GetAllContextData extracts all context data from args
func GetAllContextData(args map[string]interface{}) map[string]interface{} {
	contextData := make(map[string]interface{})
	for key, value := range args {
		if strings.HasPrefix(key, "context_") {
			actualKey := strings.TrimPrefix(key, "context_")
			contextData[actualKey] = value
		}
	}
	return contextData
}

// =====================================================
// TYPE CREATION HELPERS
// =====================================================

// createScalarType creates a scalar type definition
func createScalarType(scalarType string) GraphQLTypeDefinition {
	return GraphQLTypeDefinition{
		Kind:       "scalar",
		ScalarType: scalarType,
		Name:       scalarType,
	}
}

// createObjectType creates an object type definition
func createObjectType(name string, fields map[string]interface{}) GraphQLTypeDefinition {
	return GraphQLTypeDefinition{
		Kind:   "object",
		Name:   name,
		Fields: fields,
	}
}

// createListType creates a list type definition
func createListType(ofType GraphQLTypeDefinition) GraphQLTypeDefinition {
	return GraphQLTypeDefinition{
		Kind:   "list",
		OfType: &ofType,
	}
}

// createNonNullType creates a non-null type definition
func createNonNullType(ofType GraphQLTypeDefinition) GraphQLTypeDefinition {
	return GraphQLTypeDefinition{
		Kind:   "non_null",
		OfType: &ofType,
	}
}

// convertObjectFieldsToGraphQLFields converts ObjectFieldDef map to GraphQL field definitions
func convertObjectFieldsToGraphQLFields(fields map[string]ObjectFieldDef) map[string]interface{} {
	result := make(map[string]interface{})

	for fieldName, fieldDef := range fields {
		var fieldType GraphQLTypeDefinition

		// Start with the base type
		if isScalarType(fieldDef.Type) {
			fieldType = createScalarType(fieldDef.Type)
		} else {
			// For object types, create a reference
			fieldType = GraphQLTypeDefinition{
				Kind: "object",
				Name: fieldDef.Type,
			}
		}

		// Apply list wrapper if needed
		if fieldDef.List {
			if fieldDef.ListOfNonNull {
				fieldType = createListType(createNonNullType(fieldType))
			} else {
				fieldType = createListType(fieldType)
			}
		}

		// Apply non-null wrapper if needed
		if !fieldDef.Nullable {
			fieldType = createNonNullType(fieldType)
		}

		result[fieldName] = map[string]interface{}{
			"type":        fieldType,
			"description": fieldDef.Description,
		}
	}

	return result
}

// isScalarType checks if a type is a GraphQL scalar type
func isScalarType(typeName string) bool {
	switch typeName {
	case "String", "Int", "Boolean", "Float", "ID":
		return true
	default:
		return false
	}
}
