package readme_test

import (
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/internal/testutil"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

const projectEndpoint = "http://readme-test.local/api/v1/"

func Test_Project_Get(t *testing.T) {
	t.Run("when API responds with 200", func(t *testing.T) {
		// Arrange
		expect := readme.Project{}
		mockResponse := testutil.APITestResponse{
			URL:    projectEndpoint,
			Status: 200,
			Body: `
				{
					"name": "Go Testing",
					"subdomain": "foobar",
					"jwtSecret": "123456789abcdef",
					"baseUrl": "https://developer.example.com",
					"plan": "enterprise"
				}
			`,
		}
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		got, _, err := api.Project.Get()

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected Project struct")
	})

	t.Run("when API responds with 401", func(t *testing.T) {
		// Arrange
		expect := readme.APIErrorResponse{}
		expect.Error = "APIKEY_EMPTY"

		mockResponse := testutil.APITestResponse{
			URL:    projectEndpoint,
			Status: 401,
			Body:   `{"error":"APIKEY_EMPTY"}`,
		}
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		_, got, err := api.Project.Get()

		// Assert
		assert.Error(t, err, "it does not return an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 401", "it returns the expected error")
		assert.Equal(t, expect, got.APIErrorResponse, "it returns the API error response")
	})

	t.Run("when API response cannot be parsed", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    projectEndpoint,
			Status: 200,
			Body:   `[{"invalid":invalid"}]`,
		}
		api := mockResponse.New(t)

		// Act
		_, _, err := api.Project.Get()

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "unable to parse API response: invalid character", "it returns the expected error")
	})
}
