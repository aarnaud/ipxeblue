// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/Bootentries/{username}": {
            "delete": {
                "description": "Delete Bootentry",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Delete Bootentry",
                "parameters": [
                    {
                        "maxLength": 36,
                        "minLength": 36,
                        "type": "string",
                        "description": "Bootentry UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "Failed to parse UUID",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Bootentry UUID not found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/bootentries": {
            "get": {
                "description": "List of Bootentry filtered or not",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "List Bootentries",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "_start",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Bootentry"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a Bootentry",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create Bootentry",
                "parameters": [
                    {
                        "description": "json format Bootentry",
                        "name": "bootentry",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Bootentry"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Bootentry"
                        }
                    },
                    "400": {
                        "description": "Failed to create account in DB",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Unmarshall error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/bootentries/{id}": {
            "get": {
                "description": "Get a Bootentry by Id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get Bootentry",
                "parameters": [
                    {
                        "maxLength": 36,
                        "minLength": 36,
                        "type": "string",
                        "description": "Bootentry UUID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Bootentry"
                        }
                    },
                    "404": {
                        "description": "Computer with uuid %s not found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/bootentries/{username}": {
            "put": {
                "description": "Update a Bootentry",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update Bootentry",
                "parameters": [
                    {
                        "maxLength": 36,
                        "minLength": 36,
                        "type": "string",
                        "description": "Bootentry UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "json format of Bootentry",
                        "name": "bootentry",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Bootentry"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Bootentry"
                        }
                    },
                    "400": {
                        "description": "Query uuid and uuid miss match",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Bootentry UUID not found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Unmarshall error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/computers": {
            "get": {
                "description": "List of computers filtered or not",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "List computers",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "_start",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Computer"
                            }
                        }
                    }
                }
            }
        },
        "/computers/{id}": {
            "get": {
                "description": "Get a computer by Id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get computer",
                "parameters": [
                    {
                        "maxLength": 36,
                        "minLength": 36,
                        "type": "string",
                        "description": "Computer UUID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Computer"
                        }
                    },
                    "404": {
                        "description": "Computer with uuid %s not found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a computer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update computer",
                "parameters": [
                    {
                        "maxLength": 36,
                        "minLength": 36,
                        "type": "string",
                        "description": "Computer UUID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "json format computer",
                        "name": "computer",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Computer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Computer"
                        }
                    },
                    "400": {
                        "description": "Query ID and UUID miss match",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Unmarshall error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a computer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Delete computer",
                "parameters": [
                    {
                        "maxLength": 36,
                        "minLength": 36,
                        "type": "string",
                        "description": "Computer UUID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "Failed to parse UUID",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/ipxeaccount": {
            "get": {
                "description": "List of accounts for ipxe",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "List iPXE account",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "_start",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Ipxeaccount"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a iPXE account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create iPXE account",
                "parameters": [
                    {
                        "description": "json format iPXE account",
                        "name": "ipxeaccount",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Ipxeaccount"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Ipxeaccount"
                        }
                    },
                    "400": {
                        "description": "Failed to create account in DB",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Unmarshall error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/ipxeaccount/{username}": {
            "get": {
                "description": "Get iPXE account by username",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get iPXE account",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Ipxeaccount"
                        }
                    },
                    "404": {
                        "description": "iPXE account not found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a iPXE account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update iPXE account",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "json format iPXE account",
                        "name": "ipxeaccount",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Ipxeaccount"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Ipxeaccount"
                        }
                    },
                    "400": {
                        "description": "Query username and username miss match",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "iPXE account not found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Unmarshall error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a iPXE account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Delete iPXE account",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Ipxeaccount"
                        }
                    },
                    "404": {
                        "description": "iPXE account not found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Bootentry": {
            "type": "object",
            "properties": {
                "computers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Computer"
                    }
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "files": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.BootentryFile"
                    }
                },
                "id": {
                    "type": "string"
                },
                "ipxe_script": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.BootentryFile": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "protected": {
                    "type": "boolean"
                }
            }
        },
        "models.Computer": {
            "type": "object",
            "properties": {
                "asset": {
                    "type": "string"
                },
                "bootentry_uuid": {
                    "type": "string"
                },
                "build_arch": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "hostname": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_seen": {
                    "type": "string"
                },
                "manufacturer": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "platform": {
                    "type": "string"
                },
                "product": {
                    "type": "string"
                },
                "serial": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Tag"
                    }
                },
                "updated_at": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "models.Error": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "models.Ipxeaccount": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "is_admin": {
                    "type": "boolean"
                },
                "last_login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.Tag": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1",
	Host:        "localhost:8080",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "ipxeblue API",
	Description: "Manage PXE boot",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
