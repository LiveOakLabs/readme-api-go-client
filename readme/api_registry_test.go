package readme_test

import (
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/internal/testutil"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

const registryEndpoint = "http://readme-test.local/api/v1/api-registry"

func Test_APIRegistry_Get(t *testing.T) {
	t.Run("when called with an existing uuid", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    registryEndpoint + "/3bbeunznlboryu0o",
			Status: 200,
			Body: `
			{
				"openapi": "3.0.2",
				"info": {
					"description": "OpenAPI Specification for Testing.",
					"version": "2.0.0",
					"title": "API Endpoints",
					"contact": {
						"name": "API Support",
						"url": "https://docs.example.com/docs/contact-support",
						"email": "support@example.com"
					}
				}
			}
		`,
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.APIRegistry.Get("3bbeunznlboryu0o")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, mockResponse.Body, got, "it returns the API definition string")
	})

	t.Run("when API responds with 404", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    registryEndpoint + "/invalid",
			Status: 404,
			Body:   `{"error": "REGISTRY_NOTFOUND"}`,
		}
		expect := readme.APIErrorResponse{}
		expect.Error = "REGISTRY_NOTFOUND"
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		_, gotResponse, err := api.APIRegistry.Get("invalid")
		got := gotResponse.APIErrorResponse

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 404", "it returns the expected error")
		assert.Equal(
			t,
			expect,
			got,
			"returns expected APIErrorResponse struct",
			"it returns an APIErrorResponse struct",
		)
	})
}

func Test_APIRegistry_Create(t *testing.T) {
	// Arrange
	expect := make(map[string]interface{})

	testDefinition := `{
		"openapi": "3.0.0",
		"info": {
			"version": "1.0.0",
			"title": "Test Pet Store Test again",
			"license": {
				"name": "MIT"
			}
		}
	}
	`
	mockResponse := testutil.APITestResponse{
		URL:    registryEndpoint,
		Status: 201,
		Body:   `{"registryUUID": "abcdefghijklmno", "definition": ` + testDefinition + `}`,
	}
	api := mockResponse.New(t)

	testutil.JsonToStruct(t, testDefinition, &expect)

	t.Run("when called with a definition and no version specified", func(t *testing.T) {
		// Act
		got, _, err := api.APIRegistry.Create(testDefinition)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, "abcdefghijklmno", got.RegistryUUID, "it returns the expected uuid")
		assert.Equal(t, expect, got.Definition, "it returns the expected definition")
	})

	t.Run("when called with a version specified", func(t *testing.T) {
		// Act
		got, apiResponse, err := api.APIRegistry.Create(testDefinition, "1.2.3")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, "1.2.3", apiResponse.Request.Version, "it returns the expected response")
		assert.Equal(t, "abcdefghijklmno", got.RegistryUUID, "it returns the expected uuid")
		assert.Equal(t, expect, got.Definition, "it returns the expected definition")
	})

	t.Run("when API responds with 400", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    registryEndpoint,
			Status: 400,
			Body:   `{"error": "ERROR_SPEC_INVALID"}`,
		}
		api := mockResponse.New(t)

		// Act
		_, apiResponse, err := api.APIRegistry.Create("")

		// Assert
		assert.Equal(t, "ERROR_SPEC_INVALID", apiResponse.APIErrorResponse.Error, "it returns the API error")
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400", "it returns expected application error")
	})
}

func Test_APIRegistry_validateRegistryUUID(t *testing.T) {
	// Arrange
	testCases := []struct {
		value  string
		expect bool
		msg    string
	}{
		{"uuid:3bbeunznlboryu0o", true, "true when uuid matches pattern and is 16 chars"},
		{"uuid:4xypvmzowqax8hqiix0pqtjx", true, "true when uuid matches pattern and is 24 chars"},
		{"uuid:4xypvmzowqax8hqiix0pxtjxt", false, "false when uuid is >= 25 chars"},
		{"uuid:3bbe0nz", false, "false when uuid is too short"},
		{"notauuid", false, "false when uuid does not match pattern"},
		{`{"openapi": "3.0.0"}`, false, "false when api definition provided instead of uuid"},
	}
	for _, tc := range testCases {
		// Act
		// NOTE: `readme.ParseUUID()` refers to a private `parseUUID()` function that has been exported for tests.
		got, _ := readme.ParseUUID(tc.value)

		// Assert
		assert.Equal(t, tc.expect, got)
	}
}
