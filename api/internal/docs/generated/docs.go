// Package generated Code generated by swaggo/swag. DO NOT EDIT
package generated

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
        "/client": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "summary": "Update a client.",
                "operationId": "jwt.Auth =\u003e user.UpdateClient",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Client ID",
                        "name": "id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "default": false,
                        "description": "Newsletter",
                        "name": "newsletter",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Password updated"
                    },
                    "400": {
                        "description": "Invalid email, password or token"
                    },
                    "404": {
                        "description": "Client not found"
                    },
                    "409": {
                        "description": "Client already validated"
                    },
                    "410": {
                        "description": "Token expired"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/client/export": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "summary": "Export all data of the connected client.",
                "operationId": "jwt.Auth =\u003e user.ExportClients",
                "responses": {
                    "200": {
                        "description": "Clients exported"
                    },
                    "500": {
                        "description": "Internal"
                    }
                }
            }
        },
        "/client/register": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "summary": "Register a client.",
                "operationId": "user.RegisterClient",
                "parameters": [
                    {
                        "type": "string",
                        "format": "email",
                        "default": "user-thetiptop@yopmail.com",
                        "description": "Email address",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Aa1@azetyuiop",
                        "description": "Password",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "default": true,
                        "description": "CGU",
                        "name": "cgu",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "default": false,
                        "description": "Newsletter",
                        "name": "newsletter",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Client created"
                    },
                    "400": {
                        "description": "Invalid email or password"
                    },
                    "409": {
                        "description": "Client already exists"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/client/{id}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "summary": "Get a client by ID.",
                "operationId": "jwt.Auth =\u003e user.GetClient",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Client ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Client details"
                    },
                    "400": {
                        "description": "Invalid client ID"
                    },
                    "404": {
                        "description": "Client not found"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "summary": "Delete a client by ID.",
                "operationId": "jwt.Auth =\u003e user.DeleteClient",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Client ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Client deleted"
                    },
                    "400": {
                        "description": "Invalid client ID"
                    },
                    "404": {
                        "description": "Client not found"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/code/error": {
            "get": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Error"
                ],
                "summary": "List all code errors.",
                "operationId": "code.ListErrors",
                "responses": {}
            }
        },
        "/employee": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Employee"
                ],
                "summary": "Update a employee.",
                "operationId": "jwt.Auth =\u003e user.UpdateEmployee",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Employee ID",
                        "name": "id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "default": false,
                        "description": "Newsletter",
                        "name": "newsletter",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Password updated"
                    },
                    "400": {
                        "description": "Invalid email, password or token"
                    },
                    "404": {
                        "description": "Employee not found"
                    },
                    "409": {
                        "description": "Employee already validated"
                    },
                    "410": {
                        "description": "Token expired"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/employee/register": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Employee"
                ],
                "summary": "Register a employee.",
                "operationId": "user.RegisterEmployee",
                "parameters": [
                    {
                        "type": "string",
                        "format": "email",
                        "default": "user-thetiptop@yopmail.com",
                        "description": "Email address",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Aa1@azetyuiop",
                        "description": "Password",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Employee created"
                    },
                    "400": {
                        "description": "Invalid email or password"
                    },
                    "409": {
                        "description": "Employee already exists"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/employee/{id}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Employee"
                ],
                "summary": "Get a employee by ID.",
                "operationId": "jwt.Auth =\u003e user.GetEmployee",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Employee ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Employee details"
                    },
                    "400": {
                        "description": "Invalid employee ID"
                    },
                    "404": {
                        "description": "Employee not found"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Employee"
                ],
                "summary": "Delete a client by ID.",
                "operationId": "jwt.Auth =\u003e user.DeleteEmployee",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Employee ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Employee deleted"
                    },
                    "400": {
                        "description": "Invalid employee ID"
                    },
                    "404": {
                        "description": "Employee not found"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/game/random": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Game"
                ],
                "summary": "Get a random ticket.",
                "operationId": "jwt.Auth =\u003e game.GetTicket",
                "responses": {
                    "200": {
                        "description": "Ticket details"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/game/ticket": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Game"
                ],
                "summary": "Update a ticket.",
                "operationId": "jwt.Auth =\u003e game.UpdateTicket",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Ticket ID",
                        "name": "id",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ticket details"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Not found"
                    }
                }
            }
        },
        "/game/ticket/{id}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Game"
                ],
                "summary": "Get ticket by id.",
                "operationId": "jwt.Auth =\u003e game.GetTicketById",
                "responses": {
                    "200": {
                        "description": "Tickets details"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "404": {
                        "description": "Not found"
                    }
                }
            }
        },
        "/game/tickets": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Game"
                ],
                "summary": "List all tickets likend to the authenticated user.",
                "operationId": "jwt.Auth =\u003e game.GetTickets",
                "responses": {
                    "200": {
                        "description": "Tickets details"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "404": {
                        "description": "Not found"
                    }
                }
            }
        },
        "/status/healthcheck": {
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
                "operationId": "status.HealthCheck",
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/status/ip": {
            "get": {
                "description": "get the ip of user.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Status"
                ],
                "summary": "Show the ip of user.",
                "operationId": "status.IP",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/user/auth": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Authenticate a client/employees.",
                "operationId": "user.UserAuth",
                "parameters": [
                    {
                        "type": "string",
                        "format": "email",
                        "default": "user-thetiptop@yopmail.com",
                        "description": "Email address",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Aa1@azetyuiop",
                        "description": "Password",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Client signed in"
                    },
                    "400": {
                        "description": "Invalid email or password"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/user/auth/renew": {
            "get": {
                "consumes": [
                    "*/*",
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Renew JWT for a client/employees.",
                "operationId": "user.UserAuthRenew",
                "parameters": [
                    {
                        "type": "string",
                        "description": "With the bearer started",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT token renewed"
                    },
                    "400": {
                        "description": "Invalid token"
                    },
                    "401": {
                        "description": "Token expired"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/user/password": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Update a client/employees password.",
                "operationId": "jwt.Auth =\u003e user.CredentialUpdate",
                "parameters": [
                    {
                        "type": "string",
                        "format": "email",
                        "default": "user-thetiptop@yopmail.com",
                        "description": "Email address",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Aa1@azetyuiop",
                        "description": "Password",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Token",
                        "name": "token",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Password updated"
                    },
                    "400": {
                        "description": "Invalid email, password or token"
                    },
                    "404": {
                        "description": "Client not found"
                    },
                    "409": {
                        "description": "Client already validated"
                    },
                    "410": {
                        "description": "Token expired"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/user/register/validation": {
            "put": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Validate a client/employees email.",
                "operationId": "user.MailValidation",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token",
                        "name": "token",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "email",
                        "default": "user-thetiptop@yopmail.com",
                        "description": "Email address",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Client email validate"
                    },
                    "400": {
                        "description": "Invalid email or token"
                    },
                    "404": {
                        "description": "Client not found"
                    },
                    "409": {
                        "description": "Client already validated"
                    },
                    "410": {
                        "description": "Token expired"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/user/validation/renew": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Recover a client/employees validation type.",
                "operationId": "user.ValidationRecover",
                "parameters": [
                    {
                        "type": "string",
                        "format": "email",
                        "default": "user-thetiptop@yopmail.com",
                        "description": "Email address",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "enum": [
                            "mail",
                            "password",
                            "phone"
                        ],
                        "type": "string",
                        "description": "Type of validation",
                        "name": "type",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "dev",
	Host:             "localhost",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "TheTipTop",
	Description:      "TheTipTop API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
