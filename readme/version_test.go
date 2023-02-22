package readme_test

import (
	"fmt"
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/internal/testutil"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

const versionTestEndpoint = "http://readme-test.local/api/v1/version"

// mockVersion represents a version of a project.
// This is the response from the API when calling Get().
var mockVersion = readme.Version{
	Codename:     "",
	CreatedAt:    "2022-12-04T19:28:15.190Z",
	ID:           "638cf4cfdea3ff0096d1a95a",
	IsBeta:       false,
	IsDeprecated: false,
	IsHidden:     false,
	IsStable:     true,
	Version:      "1.0.0",
	VersionClean: "1.0.0",
}

// mockVersionSummary represents a list of versions of a project.
// This is the response from the API when calling GetAll().
var mockVersionSummary = []readme.VersionSummary{
	{
		ForkedFrom:   "",
		Codename:     mockVersion.Codename,
		CreatedAt:    mockVersion.CreatedAt,
		ID:           mockVersion.ID,
		IsBeta:       mockVersion.IsBeta,
		IsDeprecated: mockVersion.IsDeprecated,
		IsHidden:     mockVersion.IsHidden,
		IsStable:     mockVersion.IsStable,
		Version:      mockVersion.Version,
		VersionClean: mockVersion.VersionClean,
	},
	{
		ForkedFrom:   "1.0.0",
		Codename:     "",
		CreatedAt:    "2022-12-04T19:28:15.190Z",
		ID:           "639fe99365417e008daa1f55",
		IsBeta:       false,
		IsDeprecated: false,
		IsHidden:     false,
		IsStable:     true,
		Version:      "1.1.0",
		VersionClean: "1.1.0",
	},
}

func Test_Version_GetAll(t *testing.T) {
	t.Run("when called with no parameters", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    versionTestEndpoint,
			Status: 200,
			Body:   testutil.StructToJson(t, mockVersionSummary),
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.Version.GetAll()

		// Assert
		assert.NoError(t, err, "it does not returns an error")
		assert.Equal(t, mockVersionSummary, got, "it returns slice of Version structs")
	})

	// Test API error responses.
	testCases := []struct {
		status int
		errMsg string
	}{
		{400, "APIKEY_EMPTY"},
		{403, "APIKEY_MISMATCH"},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("when API responds with %v", tc.status)
		t.Run(testName, func(t *testing.T) {
			// Arrange
			mockResponse := testutil.APITestResponse{
				URL:    versionTestEndpoint,
				Status: tc.status,
				Body:   `{"error":"` + tc.errMsg + `"}`,
			}
			api := mockResponse.New(t)

			expectErr := fmt.Sprintf("API responded with a non-OK status: %v", tc.status)

			// Act
			_, apiResponse, err := api.Version.GetAll()

			// Assert
			assert.Error(t, err, "it returns an error")
			assert.ErrorContains(t, err, expectErr, "it returns expected API error")
			assert.Equal(t, apiResponse.APIErrorResponse.Error, tc.errMsg, "it returns an APIErrorResponse")
			assert.Equal(
				t,
				apiResponse.HTTPResponse.StatusCode,
				tc.status,
				"it returns an HTTPResponse with expected status code",
			)
		})
	}
}

func Test_Version_Get(t *testing.T) {
	t.Run("when called with existing version number", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    versionTestEndpoint + "/1.0.0",
			Status: 200,
			Body:   testutil.StructToJson(t, mockVersion),
		}
		api := mockResponse.New(t)
		// Act
		got, _, err := api.Version.Get("1.0.0")
		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, mockVersion, got, "it returns a single Version struct")
	})

	t.Run("when requested id is invalid", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    versionTestEndpoint,
			Status: 200,
			Body:   testutil.StructToJson(t, mockVersionSummary),
		}
		api := mockResponse.New(t)
		// Act
		_, _, err := api.Version.Get("id:invalidinvalidinvalid")
		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "no match for version ID invalidinvalidinvalid", "it returns the expected error")
	})

	// Test API error responses.
	testCases := []struct {
		status int
		errMsg string
	}{
		{400, "APIKEY_EMPTY"},
		{403, "APIKEY_MISMATCH"},
		{404, "VERSION_NOTFOUND"},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("when API responds with %v", tc.status)
		t.Run(testName, func(t *testing.T) {
			// Arrange
			mockResponse := testutil.APITestResponse{
				URL:    versionTestEndpoint + "/99.99.99",
				Body:   `{"error":"` + tc.errMsg + `"}`,
				Status: tc.status,
			}
			api := mockResponse.New(t)

			expectErr := fmt.Sprintf("API responded with a non-OK status: %v", tc.status)

			// Act
			_, apiResponse, err := api.Version.Get("99.99.99")

			// Assert
			assert.Error(t, err, "it returns an error")
			assert.ErrorContains(t, err, expectErr, "it returns expected API error")
			assert.Equal(t, apiResponse.APIErrorResponse.Error, tc.errMsg, "it returns an APIErrorResponse")
			assert.Equal(
				t,
				apiResponse.HTTPResponse.StatusCode,
				tc.status,
				"it returns an HTTPResponse with expected status code",
			)
		})
	}
}

// NOTE: `readme.GetVersion()` refers to a private `getVersion()` function that has been exported for tests.
func Test_Version_GetAll_Versions(t *testing.T) {
	mockVersionBody := `
	[
		{
			"version": "1.0.0",
			"version_clean": "1.0.0",
			"codename": "",
			"is_stable": true,
			"is_beta": false,
			"is_hidden": false,
			"is_deprecated": false,
			"_id": "638cf4cfdea3ff0096d1a95a",
			"createdAt": "2022-12-04T19:28:15.190Z"
		}
	]
	`

	t.Run("when looking up version number by ID", func(t *testing.T) {
		// Arrange
		var expect []readme.VersionSummary
		mockResponse := testutil.APITestResponse{
			URL:    versionTestEndpoint,
			Status: 200,
			Body:   mockVersionBody,
		}
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		got, err := readme.GetVersion(api.Version.(*readme.VersionClient), "id:638cf4cfdea3ff0096d1a95a")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, "1.0.0", got, "version matches", "it returns the expected version")
	})

	t.Run("when looking up by ID encounters error", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    versionTestEndpoint,
			Status: 400,
			Body:   "[]",
		}
		api := mockResponse.New(t)

		// Act
		_, err := readme.GetVersion(api.Version.(*readme.VersionClient), "id:638cf4cfdea3ff0096d1a95a")

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "json: cannot unmarshal array into Go value", "it returns the expected error")
	})
}

func Test_Version_Create_and_Update(t *testing.T) {
	// Arrange
	expect := readme.Version{}
	mockResponseBody := `
		{
			"_id": "638cf4cfdea3ff0096d1a95a",
			"categories": [
				"638cf4cfdea3ff0096d1a95c"
			],
			"codename": "A test version",
			"createdAt": "2022-12-04T19:28:15.190Z",
			"forked_from": "638cf4cfdea3ff0096d1a95a",
			"is_beta": false,
			"is_deprecated": false,
			"is_hidden": false,
			"is_stable": true,
			"project": "638cf4cedea3ff0096d1a955",
			"releaseDate": "2022-12-04T19:28:15.190Z",
			"version": "1.1.1",
			"version_clean": "1.1.1"
		}
	`
	testutil.JsonToStruct(t, mockResponseBody, &expect)

	mockVersion := readme.VersionParams{
		From:     "1.0.0",
		Version:  "1.1.0",
		Codename: "A test version",
	}

	t.Run("when create is successful", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    versionTestEndpoint,
			Status: 200,
			Body:   mockResponseBody,
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.Version.Create(mockVersion)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected VersionSaved struct")
	})

	t.Run("when update is called with a valid version and API response is 200", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    versionTestEndpoint + "/1.1.1",
			Status: 200,
			Body:   mockResponseBody,
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.Version.Update("1.1.1", mockVersion)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected VersionSaved struct")
	})

	t.Run("when update is called with a non-existing version and the API responds with 400", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    versionTestEndpoint + "/1.1.x",
			Status: 404,
			Body:   mockResponseBody,
		}
		api := mockResponse.New(t)

		// Act
		_, _, err := api.Version.Update("1.1.x", mockVersion)

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 404", "it returns the expected error")
	})
}

func Test_Version_Delete(t *testing.T) {
	t.Run("when called with valid existing version and API response is 200", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    versionTestEndpoint + "/1.0.0",
			Status: 200,
			Body:   `{"removed":true}`,
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.Version.Delete("1.0.0")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.True(t, got, "it returns true")
	})

	t.Run("when API responds with 400", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    versionTestEndpoint + "/1.0.0",
			Status: 400,
			Body:   `{"error": "VERSION_CANT_REMOVE_STABLE"}`,
		}
		api := mockResponse.New(t)

		// Act
		got, apiErrorResponse, err := api.Version.Delete("1.0.0")

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400", "it returns the expected error")
		assert.Equal(t, apiErrorResponse.APIErrorResponse.Error, "VERSION_CANT_REMOVE_STABLE",
			"it returns an APIErrorResponse")
		assert.False(t, got, "it returns false")
	})
}
