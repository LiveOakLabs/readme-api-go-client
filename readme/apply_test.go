package readme_test

import (
	"errors"
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_Apply_Get(t *testing.T) {
	// Arrange
	expect := testdata.ApplyOpenRoles
	gock.New(TestClient.APIURL).
		Get("/apply").
		Reply(200).
		JSON(expect)
	defer gock.Off()

	// Act
	got, _, err := TestClient.Apply.Get()

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns open roles")
	assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
}

func Test_Apply_Apply(t *testing.T) {
	t.Run("when called with valid params", func(t *testing.T) {
		// Arrange
		expect := testdata.ApplyResponseSuccess
		gock.New("http://readme-test.local/api/v1").
			Post(readme.ApplyEndpoint).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.Apply.Apply(testdata.ApplyApplication)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns the expected API response message")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when called with invalid params", func(t *testing.T) {
		// Arrange
		expectResponse := testdata.ApplyCreateResponseInvalidName
		expectError := errors.New("API responded with a non-OK status: 400")

		gock.New("http://readme-test.local/api/v1").
			Post(readme.ApplyEndpoint).
			Reply(400).
			JSON(expectResponse.Body)
		defer gock.Off()

		// Act
		application := readme.Application{
			Name:  "",
			Email: "",
			Job:   "Front End Engineer",
		}
		_, got, err := TestClient.Apply.Apply(application)

		// Assert
		assert.Equal(t, expectError, err, "it returns the expected error")
		assert.Equal(t, expectResponse.APIErrorResponse,
			got.APIErrorResponse, "it returns the expected API error response")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})
}
