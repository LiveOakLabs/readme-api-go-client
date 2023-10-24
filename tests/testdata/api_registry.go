package testdata

import "github.com/liveoaklabs/readme-api-go-client/readme"

var APIDefinition = map[string]interface{}{
	"openapi": "3.0.2",
	"info": map[string]interface{}{
		"description": "OpenAPI Specification for Testing.",
		"version":     "2.0.0",
		"title":       "API Endpoints",
		"contact": map[string]interface{}{
			"name":  "API Support",
			"url":   "https://docs.example.com/docs/contact-support",
			"email": "support@example.com",
		},
	},
}

func APIRegistrySaved(uuid string) readme.APIRegistrySaved {
	return readme.APIRegistrySaved{
		Definition:   APIDefinition,
		RegistryUUID: uuid,
	}
}
