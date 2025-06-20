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
        "/api/v1/files/upload": {
            "post": {
                "description": "Upload one or more files",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Files"
                ],
                "summary": "Upload files",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "file"
                        },
                        "collectionFormat": "multi",
                        "description": "Files to upload",
                        "name": "files",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully uploaded files",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/services.UploadResult"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request or file upload failed",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/messages/{messageId}": {
            "delete": {
                "description": "Delete a message by message ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Delete a message by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Message ID",
                        "name": "messageId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Message successfully deleted"
                    },
                    "400": {
                        "description": "Invalid message ID format",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Message not found",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/threads/{threadId}": {
            "get": {
                "description": "Get a thread and its messages by thread ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Threads"
                ],
                "summary": "Get a thread by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Thread ID",
                        "name": "threadId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Thread details with messages",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "created_at": {
                                    "type": "string"
                                },
                                "id": {
                                    "type": "string"
                                },
                                "messages": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "properties": {
                                            "created_at": {
                                                "type": "string"
                                            },
                                            "event_type": {
                                                "type": "string"
                                            },
                                            "id": {
                                                "type": "string"
                                            },
                                            "message_index": {
                                                "type": "integer"
                                            },
                                            "meta_data": {
                                                "type": "string"
                                            },
                                            "query_text": {
                                                "type": "string"
                                            },
                                            "response_text": {
                                                "type": "string"
                                            },
                                            "stream_status": {
                                                "type": "string"
                                            },
                                            "version": {
                                                "type": "integer"
                                            }
                                        }
                                    }
                                },
                                "slug": {
                                    "type": "string"
                                },
                                "updated_at": {
                                    "type": "string"
                                },
                                "version": {
                                    "type": "integer"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid thread ID format",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Thread not found",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a thread by thread ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Threads"
                ],
                "summary": "Delete a thread by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Thread ID",
                        "name": "threadId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Thread successfully deleted"
                    },
                    "400": {
                        "description": "Invalid thread ID format",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Thread not found",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "helpers.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "object",
                    "properties": {
                        "code": {
                            "type": "string"
                        },
                        "message": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "services.UploadResult": {
            "type": "object",
            "properties": {
                "file_name": {
                    "type": "string"
                },
                "file_size_bytes": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "mime_type": {
                    "type": "string"
                },
                "original_file_name": {
                    "type": "string"
                },
                "uploaded_at": {
                    "type": "string"
                },
                "version": {
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
	Schemes:          []string{},
	Title:            "Agios API Documentation",
	Description:      "This is the API documentation for the Agios application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
