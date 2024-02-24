// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
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
        "/api/v1/fizzbuzz/:int1/:int2/:limit/:str1/:str2": {
            "get": {
                "description": "Returns a list of strings with numbers from 1 to limit, where: \\n all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "FizzBuzz"
                ],
                "summary": "Return FizzBuzz result.",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 3,
                        "description": "Give the first number",
                        "name": "int1",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 5,
                        "description": "Give the second number",
                        "name": "int2",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Give limit of fizzbuzz",
                        "name": "limit",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "fizz",
                        "description": "Give the first word",
                        "name": "str1",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "buzz",
                        "description": "Give the second word",
                        "name": "str2",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/fizzbuzz/stats": {
            "get": {
                "description": "Return the parameters corresponding to the most used request, as well as the number of hits for this request.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "FizzBuzz"
                ],
                "summary": "Return FizzBuzz statistics.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/fizzbuzz.Stats"
                        }
                    }
                }
            }
        },
        "/api/v1/status/healthcheck": {
            "get": {
                "description": "get the status of server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Status"
                ],
                "summary": "Show the status of server.",
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        }
    },
    "definitions": {
        "fizzbuzz.Request": {
            "type": "object",
            "properties": {
                "int1": {
                    "description": "Int1 is the first integer to be used as a divisor.",
                    "type": "integer"
                },
                "int2": {
                    "description": "Int2 is the second integer to be used as a divisor.",
                    "type": "integer"
                },
                "limit": {
                    "description": "Limit is the maximum number of iterations to perform.",
                    "type": "integer"
                },
                "str1": {
                    "description": "Str1 is the string to print when the current iteration is divisible by Int1.",
                    "type": "string"
                },
                "str2": {
                    "description": "Str2 is the string to print when the current iteration is divisible by Int2.",
                    "type": "string"
                }
            }
        },
        "fizzbuzz.Stats": {
            "type": "object",
            "properties": {
                "hits": {
                    "description": "Hits is the number of times the Request has been made.",
                    "type": "integer"
                },
                "request": {
                    "description": "Request is the Request object associated with these statistics.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/fizzbuzz.Request"
                        }
                    ]
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
