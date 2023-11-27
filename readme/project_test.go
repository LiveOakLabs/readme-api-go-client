package readme_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_Project_Get(t *testing.T) {
	t.Run("when API responds with 200", func(t *testing.T) {
		// Arrange
		expect := testdata.Project
		gock.New(TestClient.APIURL).
			Get("/").
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.Project.Get()

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected Project struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when API responds with 401", func(t *testing.T) {
		// Arrange
		expect := readme.APIErrorResponse{Error: "APIKEY_EMPTY"}
		gock.New(TestClient.APIURL).
			Get("/").
			Reply(401).
			JSON(expect)
		defer gock.Off()

		// Act
		_, got, err := TestClient.Project.Get()

		// Assert
		assert.ErrorContains(t, err, "API responded with a non-OK status: 401",
			"it returns the expected error")
		assert.Equal(t, expect, got.APIErrorResponse,
			"it returns the API error response")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when API response cannot be parsed", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Get("/").
			Reply(200).
			JSON(`[{"invalid":invalid"}]`)
		defer gock.Off()

		// Act
		_, _, err := TestClient.Project.Get()

		// Assert
		assert.ErrorContains(t, err,
			"unable to parse API response: invalid character",
			"it returns the expected error")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}
