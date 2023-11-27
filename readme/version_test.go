package readme_test

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_Version_GetAll(t *testing.T) {
	t.Run("when called with no parameters", func(t *testing.T) {
		// Arrange
		expect := testdata.VersionSummary
		gock.New(TestClient.APIURL).
			Get(readme.VersionEndpoint).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.Version.GetAll()

		// Assert
		assert.NoError(t, err, "it does not returns an error")
		assert.Equal(t, expect, got, "it returns slice of Version structs")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
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
			expect := readme.APIErrorResponse{Error: tc.errMsg}
			expectErr := fmt.Sprintf("API responded with a non-OK status: %v", tc.status)
			gock.New(TestClient.APIURL).
				Get(readme.VersionEndpoint).
				Reply(tc.status).
				JSON(expect)
			defer gock.Off()

			// Act
			_, apiResponse, err := TestClient.Version.GetAll()

			// Assert
			assert.Error(t, err, "it returns an error")
			assert.ErrorContains(t, err, expectErr,
				"it returns expected API error")
			assert.Equal(t, apiResponse.APIErrorResponse.Error, tc.errMsg,
				"it returns an APIErrorResponse")
			assert.Equal(t, apiResponse.HTTPResponse.StatusCode, tc.status,
				"it returns an HTTPResponse with expected status code",
			)
			assert.True(t, gock.IsDone(), "it makes the expected API call")
		})
	}
}

func Test_Version_Get(t *testing.T) {
	t.Run("when called with existing version number", func(t *testing.T) {
		// Arrange
		expect := testdata.Versions[0]
		gock.New(TestClient.APIURL).
			Get(readme.VersionEndpoint + "/" + expect.Version).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.Version.Get(expect.Version)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns a single Version struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when requested id is invalid", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Get(readme.VersionEndpoint).
			Reply(200).
			JSON(testdata.Versions)
		defer gock.Off()

		// Act
		_, _, err := TestClient.Version.Get("id:invalidinvalidinvalid")

		// Assert
		assert.ErrorContains(t, err, "no match for version ID invalidinvalidinvalid",
			"it returns the expected error")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
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
			expect := readme.APIErrorResponse{Error: tc.errMsg}
			gock.New(TestClient.APIURL).
				Get(readme.VersionEndpoint + "/99.99.99").
				Reply(tc.status).
				JSON(expect)
			defer gock.Off()

			expectErr := fmt.Sprintf("API responded with a non-OK status: %v", tc.status)

			// Act
			_, apiResponse, err := TestClient.Version.Get("99.99.99")

			// Assert
			assert.ErrorContains(t, err, expectErr,
				"it returns expected API error")
			assert.Equal(t, apiResponse.APIErrorResponse.Error, tc.errMsg,
				"it returns an APIErrorResponse")
			assert.Equal(t, apiResponse.HTTPResponse.StatusCode, tc.status,
				"it returns an HTTPResponse with expected status code",
			)
			assert.True(t, gock.IsDone(), "it makes the expected API call")
		})
	}
}

// NOTE: `readme.GetVersion()` refers to a private `getVersion()` function that
// has been exported for tests.
func Test_Version_GetAll_Versions(t *testing.T) {
	t.Run("when looking up version number by ID", func(t *testing.T) {
		// Arrange
		expect := testdata.VersionSummary
		gock.New(TestClient.APIURL).
			Get(readme.VersionEndpoint).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, err := TestClient.Version.GetVersion("id:" + expect[0].ID)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect[0].Version, got, "version matches",
			"it returns the expected version")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when looking up by ID encounters unmarshal error", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Get(readme.VersionEndpoint).
			Reply(200).
			JSON(`[x]`)
		defer gock.Off()

		expect := "unable to parse API response: invalid character"

		// Act
		_, err := TestClient.Version.GetVersion("id:638cf4cfdea3ff0096d1a95a")

		// Assert
		assert.ErrorContains(t, err, expect, "it returns the expected error")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}

func Test_Version_Create_and_Update(t *testing.T) {
	t.Run("when create is successful", func(t *testing.T) {
		// Arrange
		expect := testdata.Versions[1]
		gock.New(TestClient.APIURL).
			Post(readme.VersionEndpoint).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		createParams := readme.VersionParams{
			From:    "1.0.0",
			Version: expect.Version,
		}

		// Act
		got, _, err := TestClient.Version.Create(createParams)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected VersionSaved struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when update is called with a valid version and API response is 200", func(t *testing.T) {
		// Arrange
		expect := testdata.Versions[1]
		gock.New(TestClient.APIURL).
			Put(readme.VersionEndpoint + "/" + expect.Version).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		hidden := true
		updateParams := readme.VersionParams{
			From:     "1.0.0",
			Version:  expect.Version,
			IsHidden: &hidden,
		}

		// Act
		got, _, err := TestClient.Version.Update(expect.VersionClean, updateParams)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected VersionSaved struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when update is called with a non-existing version and the API responds with 400", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Put(readme.VersionEndpoint + "/1.1.x").
			Reply(404).
			JSON(readme.APIErrorResponse{Error: "VERSION_NOTFOUND"})
		defer gock.Off()

		updateParams := readme.VersionParams{
			Version: "1.1.x",
		}

		expect := "API responded with a non-OK status: 404"

		// Act
		_, _, err := TestClient.Version.Update("1.1.x", updateParams)

		// Assert
		assert.ErrorContains(t, err, expect, "it returns the expected error")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}

func Test_Version_Delete(t *testing.T) {
	t.Run("when called with valid existing version and API response is 200", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.VersionEndpoint + "/1.0.0").
			Reply(200).
			JSON("{}")
		defer gock.Off()

		// Act
		got, _, err := TestClient.Version.Delete("1.0.0")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.True(t, got, "it returns true")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when API responds with 400", func(t *testing.T) {
		// Arrange
		expect := readme.APIErrorResponse{Error: "VERSION_CANT_REMOVE_STABLE"}
		gock.New(TestClient.APIURL).
			Delete(readme.VersionEndpoint + "/1.0.0").
			Reply(400).
			JSON(expect)
		defer gock.Off()

		// Act
		got, apiErrorResponse, err := TestClient.Version.Delete("1.0.0")

		// Assert
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400",
			"it returns the expected error")
		assert.Equal(t, apiErrorResponse.APIErrorResponse.Error,
			"VERSION_CANT_REMOVE_STABLE",
			"it returns an APIErrorResponse")
		assert.False(t, got, "it returns false")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}
