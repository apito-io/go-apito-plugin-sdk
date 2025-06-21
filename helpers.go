package sdk

// Field creates a basic GraphQL field
func Field(fieldType, description string) GraphQLField {
	return GraphQLField{
		Type:        fieldType,
		Description: description,
		Args:        make(map[string]interface{}),
	}
}

// FieldWithArgs creates a GraphQL field with arguments
func FieldWithArgs(fieldType, description string, args map[string]interface{}) GraphQLField {
	return GraphQLField{
		Type:        fieldType,
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
	return Field("["+itemType+"]", description)
}

// NonNullField creates a non-null type GraphQL field
func NonNullField(fieldType, description string) GraphQLField {
	return Field(fieldType+"!", description)
}

// NonNullListField creates a non-null list type GraphQL field
func NonNullListField(itemType, description string) GraphQLField {
	return Field("["+itemType+"!]!", description)
}

// ObjectField creates an object type GraphQL field with properties
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
