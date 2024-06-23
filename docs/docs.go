// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Vũ Kỳ",
            "url": "github.com/vukyn",
            "email": "vukynpro@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/storage/download/{file_name}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Download media file",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Storage"
                ],
                "summary": "Download media",
                "parameters": [
                    {
                        "type": "string",
                        "description": "file id",
                        "name": "file_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/medioa_internal_storage_models.DownloadResponse"
                        }
                    }
                }
            }
        },
        "/storage/secret": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create new secret for upload media",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Storage"
                ],
                "summary": "Create new secret",
                "parameters": [
                    {
                        "description": "create secret request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/medioa_internal_storage_models.CreateSecretRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/medioa_internal_storage_models.CreateSecretResponse"
                        }
                    }
                }
            }
        },
        "/storage/secret/download/{file_name}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Download media file",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Storage"
                ],
                "summary": "Download media with secret",
                "parameters": [
                    {
                        "type": "string",
                        "description": "file id",
                        "name": "file_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "secret",
                        "name": "secret",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/medioa_internal_storage_models.DownloadResponse"
                        }
                    }
                }
            }
        },
        "/storage/secret/pin": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Reset pin code for secret",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Storage"
                ],
                "summary": "Reset pin code",
                "parameters": [
                    {
                        "description": "reset pin request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/medioa_internal_storage_models.ResetPinCodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/storage/secret/retrieve": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retrieve secret with new access token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Storage"
                ],
                "summary": "Retrieve secret",
                "parameters": [
                    {
                        "description": "retrieve secrect request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/medioa_internal_storage_models.RetrieveSecretRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/medioa_internal_storage_models.RetrieveSecretResponse"
                        }
                    }
                }
            }
        },
        "/storage/secret/upload": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Upload media file (images, videos, etc.)",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Storage"
                ],
                "summary": "Upload media with secret",
                "parameters": [
                    {
                        "type": "string",
                        "description": "session id",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "secret",
                        "name": "secret",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "binary file",
                        "name": "chunk",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "file name",
                        "name": "filename",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/medioa_internal_storage_models.UploadResponse"
                        }
                    }
                }
            }
        },
        "/storage/upload": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Upload media file (images, videos, etc.)",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Storage"
                ],
                "summary": "Upload media",
                "parameters": [
                    {
                        "type": "string",
                        "description": "session id",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "file",
                        "description": "binary file",
                        "name": "chunk",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "file name",
                        "name": "file_name",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/medioa_internal_storage_models.UploadResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "medioa_internal_storage_models.CreateSecretRequest": {
            "type": "object",
            "properties": {
                "master_key": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "pin_code": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "medioa_internal_storage_models.CreateSecretResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "medioa_internal_storage_models.DownloadResponse": {
            "type": "object",
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        },
        "medioa_internal_storage_models.ResetPinCodeRequest": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "new_pin_code": {
                    "type": "string"
                }
            }
        },
        "medioa_internal_storage_models.RetrieveSecretRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "medioa_internal_storage_models.RetrieveSecretResponse": {
            "type": "object",
            "properties": {
                "access_key": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "medioa_internal_storage_models.UploadResponse": {
            "type": "object",
            "properties": {
                "ext": {
                    "type": "string"
                },
                "file_id": {
                    "type": "string"
                },
                "file_name": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Medioa API",
	Description:      "Medioa REST API (with gin-gonic).",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
