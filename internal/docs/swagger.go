package docs

import (
	"encoding/json"
)

// OpenAPISpec represents the OpenAPI specification
type OpenAPISpec struct {
	OpenAPI string   `json:"openapi"`
	Info    Info     `json:"info"`
	Servers []Server `json:"servers"`
	Paths   Paths    `json:"paths"`
}

// Info represents API info
type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// Server represents API server
type Server struct {
	URL string `json:"url"`
}

// Paths represents API paths
type Paths struct {
	Deployments  PathItem `json:"/api/deployments"`
	Compositions PathItem `json:"/api/compositions"`
	Cost         PathItem `json:"/api/cost/monthly"`
	Health       PathItem `json:"/api/health/status"`
	Stream       PathItem `json:"/api/stream"`
}

// PathItem represents a path item
type PathItem struct {
	Get Operation `json:"get"`
}

// Operation represents an operation
type Operation struct {
	Summary     string           `json:"summary"`
	Description string           `json:"description"`
	Responses   map[int]Response `json:"responses"`
	Tags        []string         `json:"tags"`
	Parameters  []Parameter      `json:"parameters"`
}

// Response represents a response
type Response struct {
	Description string             `json:"description"`
	Content     map[string]Content `json:"content"`
}

// Content represents content
type Content struct {
	Schema Schema `json:"schema"`
}

// Schema represents a schema
type Schema struct {
	Type  string  `json:"type"`
	Items *Schema `json:"items,omitempty"`
}

// Parameter represents a parameter
type Parameter struct {
	Name     string `json:"name"`
	In       string `json:"in"`
	Required bool   `json:"required"`
	Schema   Schema `json:"schema"`
}

// GetSpec returns the OpenAPI specification
func GetSpec() OpenAPISpec {
	return OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       "Sovereign Engine API",
			Version:     "1.0.0",
			Description: "API for Sovereign Engine cloud infrastructure management",
		},
		Servers: []Server{
			{URL: "http://localhost:8080"},
		},
		Paths: Paths{
			Deployments: PathItem{
				Get: Operation{
					Summary:     "List deployments",
					Description: "Returns a list of all deployments",
					Responses: map[int]Response{
						200: {
							Description: "Successful response",
							Content: map[string]Content{
								"application/json": {
									Schema: Schema{Type: "array", Items: &Schema{Type: "object"}},
								},
							},
						},
					},
					Tags: []string{"deployments"},
				},
			},
			Compositions: PathItem{
				Get: Operation{
					Summary:     "List compositions",
					Description: "Returns a list of available compositions",
					Responses: map[int]Response{
						200: {
							Description: "Successful response",
							Content: map[string]Content{
								"application/json": {
									Schema: Schema{Type: "array", Items: &Schema{Type: "object"}},
								},
							},
						},
					},
					Tags: []string{"compositions"},
				},
			},
			Cost: PathItem{
				Get: Operation{
					Summary:     "Get monthly cost",
					Description: "Returns monthly cost data for a team",
					Responses: map[int]Response{
						200: {
							Description: "Successful response",
							Content: map[string]Content{
								"application/json": {
									Schema: Schema{Type: "object"},
								},
							},
						},
					},
					Tags: []string{"cost"},
					Parameters: []Parameter{
						{Name: "team", In: "query", Required: false, Schema: Schema{Type: "string"}},
					},
				},
			},
			Health: PathItem{
				Get: Operation{
					Summary:     "Health check",
					Description: "Returns system health status",
					Responses: map[int]Response{
						200: {
							Description: "Successful response",
							Content: map[string]Content{
								"application/json": {
									Schema: Schema{Type: "object"},
								},
							},
						},
					},
					Tags: []string{"health"},
				},
			},
			Stream: PathItem{
				Get: Operation{
					Summary:     "SSE stream",
					Description: "Server-sent events for real-time updates",
					Responses: map[int]Response{
						200: {
							Description: "Successful response",
							Content: map[string]Content{
								"text/event-stream": {
									Schema: Schema{Type: "string"},
								},
							},
						},
					},
					Tags: []string{"stream"},
				},
			},
		},
	}
}

// ToJSON converts the spec to JSON
func (s OpenAPISpec) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}
