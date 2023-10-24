package readme_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_APIRegistry_Get(t *testing.T) {
	t.Run("when called with an existing uuid", func(t *testing.T) {
		// Arrange
		expect := testdata.ToJSON(testdata.APIDefinition) + "\n"
		uuid := "3bbeunznlboryu0o"
		gock.New(TestClient.APIURL).
			Get(readme.APIRegistryEndpoint + "/" + uuid).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.APIRegistry.Get(uuid)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns the expected body")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when API responds with 404", func(t *testing.T) {
		// Arrange
		expect := testdata.APISpecResponseVersionEmtpy
		gock.New(TestClient.APIURL).
			Get(readme.APIRegistryEndpoint + "/invalid").
			Reply(404).
			JSON(expect)
		defer gock.Off()

		// Act
		got, gotResponse, err := TestClient.APIRegistry.Get("invalid")

		// Assert
		assert.ErrorContains(t, err, "API responded with a non-OK status: 404",
			"it returns the expected error")
		assert.Equal(t, string(expect.Body), got,
			"it returns the expected body")
		assert.Equal(t, 404, gotResponse.HTTPResponse.StatusCode,
			"it returns the expected status code")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}

func Test_APIRegistry_Create(t *testing.T) {
	t.Run("when called with a definition and no version specified", func(t *testing.T) {
		// Arrange
		expect := testdata.APIRegistrySaved("abcdefghijklmno")
		expectJSON := testdata.ToJSON(expect.Definition)
		gock.New(TestClient.APIURL).
			Post(readme.APIRegistryEndpoint).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.APIRegistry.Create(expectJSON)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns the expected definition")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when called with a version specified", func(t *testing.T) {
		// Arrange
		expect := testdata.APIRegistrySaved("abcdefghijklmno")
		gock.New(TestClient.APIURL).
			Post(readme.APIRegistryEndpoint).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, apiResponse, err := TestClient.APIRegistry.
			Create(testdata.ToJSON(expect), "1.2.3")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, "1.2.3", apiResponse.Request.Version,
			"it returns the expected response")
		assert.Equal(t, "abcdefghijklmno", got.RegistryUUID,
			"it returns the expected uuid")
		assert.Equal(t, expect, got, "it returns the expected definition")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when API responds with 400", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Post(readme.APIRegistryEndpoint).
			Reply(400).
			JSON(testdata.APISpecResponseSpecFileInvalid.APIErrorResponse)
		defer gock.Off()

		// Act
		_, apiResponse, err := TestClient.APIRegistry.Create("")

		// Assert
		assert.Equal(t, "ERROR_SPEC_INVALID", apiResponse.APIErrorResponse.Error,
			"it returns the API error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400",
			"it returns expected application error")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
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
		// NOTE: `readme.ParseUUID()` refers to a private `parseUUID()`
		// function that has been exported for tests.
		got, _ := readme.ParseUUID(tc.value)

		// Assert
		assert.Equal(t, tc.expect, got)
	}
}
