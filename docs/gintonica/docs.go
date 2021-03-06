// Package gintonica GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package gintonica

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "description": "Login, got it?.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Do user login and return JWT token.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The email of the citizen",
                        "name": "email",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Citizen's password",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/products": {
            "get": {
                "description": "Returns a list of products",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "All products.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer ",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/localdb.Product"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a product",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Create Product document.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer ",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "The data",
                        "name": "localdb.Product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "$ref": "#/definitions/localdb.Product"
                            }
                        }
                    }
                }
            }
        },
        "/products/category/{category}": {
            "get": {
                "description": "Returns a list of products from this category",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "All products from this category.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer ",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "The category you want",
                        "name": "category",
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
                                "$ref": "#/definitions/localdb.Product"
                            }
                        }
                    }
                }
            }
        },
        "/products/update/{id}": {
            "patch": {
                "description": "For real dude, it catchs the document that represents the Product, and update it.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Update Product document.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer ",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "The id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "The data",
                        "name": "localdb.Product{}",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "$ref": "#/definitions/localdb.Product"
                            }
                        }
                    }
                }
            }
        },
        "/products/view/{id}": {
            "get": {
                "description": "For real dude, it catchs the document that represents the Product.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Retrieve Product document.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer ",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "The id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "$ref": "#/definitions/localdb.Product"
                            }
                        }
                    }
                }
            }
        },
        "/products/{id}": {
            "delete": {
                "description": "For real dude, it catchs the document that represents the Product, and update it.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Delete Product document.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer ",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "The id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "localdb.Product": {
            "type": "object",
            "properties": {
                "brand": {
                    "type": "string"
                },
                "category": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "timestamp": {
                    "$ref": "#/definitions/primitive.Timestamp"
                }
            }
        },
        "primitive.Timestamp": {
            "type": "object",
            "properties": {
                "i": {
                    "type": "integer"
                },
                "t": {
                    "type": "integer"
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
	Schemes:          []string{"http"},
	Title:            "Gin Swagger Example API",
	Description:      "This is a sample server server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
