// Package docs is the package that contains the swagger docs
// this file is the placeholder for the generated docs
// when you run the swagger documentation generator,
// it will be replaced with the generated docs
// DO NOT EDIT THIS FILE
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "PlanVX",
            "url": "https://github.com/PlanVX"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "https://github.com/PlanVX/aweme/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {},
    "definitions": {}
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "aweme",
	Description:      "aweme api",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
