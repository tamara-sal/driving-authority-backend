// Package docs exposes the OpenAPI (Swagger) specification for the Driving Authority API.
package docs

import (
	_ "embed"

	"github.com/swaggo/swag"
)

//go:embed swagger.json
var swaggerJSON string

// SwaggerInfo holds exported Swagger metadata.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "api-production-5e10.up.railway.app",
	BasePath:         "/api/v1",
	Schemes:          []string{"https", "http"},
	Title:            "Driving Authority API",
	Description:      "REST API for auth, JWT, RBAC, and identity verification.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  swaggerJSON,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
