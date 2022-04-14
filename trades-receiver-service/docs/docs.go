// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/trades/api/v1": {
            "get": {
                "description": "Get all transactions objects with cursor pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "trades"
                ],
                "summary": "Get all transactions",
                "operationId": "getAllTransactions",
                "parameters": [
                    {
                        "type": "string",
                        "description": "pagination",
                        "name": "cursor",
                        "in": "path"
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.Transaction"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "integer"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.Transaction": {
            "type": "object",
            "properties": {
                "averagePrice": {
                    "type": "number"
                },
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "secFilingsID": {
                    "type": "string"
                },
                "totalShares": {
                    "type": "integer"
                },
                "totalValue": {
                    "type": "number"
                },
                "transactionTypeName": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Insider-trades trades-receiver API",
	Description:      "Receives insider trades sec forms from external api and serves out structured trades information",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
