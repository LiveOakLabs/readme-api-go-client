package readme_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_Changelog_Get(t *testing.T) {
	// Arrange
	expect := testdata.Changelogs[0]
	gock.New(TestClient.APIURL).
		Get(readme.ChangelogEndpoint + "/" + expect.Slug).
		Reply(200).
		JSON(testdata.Changelogs[0])
	defer gock.Off()

	// Act
	got, _, err := TestClient.Changelog.Get(expect.Slug)

	// Assert
	assert.NoError(t, err, "it does not returns an error")
	assert.Equal(t, expect, got, "it returns a Changelog struct")
	assert.True(t, gock.IsDone(), "it makes the expected API call")
}

func Test_Changelog_GetAll(t *testing.T) {
	// Arrange
	expect := testdata.Changelogs
	gock.New(TestClient.APIURL).
		Get(readme.ChangelogEndpoint).
		Reply(200).
		JSON(expect)
	defer gock.Off()

	// Act
	got, _, err := TestClient.Changelog.GetAll(readme.RequestOptions{Page: 1})

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns slice of Changelog structs")
	assert.True(t, gock.IsDone(), "it makes the expected API call")
}

func Test_Changelog_Create(t *testing.T) {
	t.Run("when called with valid parameters and API responds with 201", func(t *testing.T) {
		// Arrange
		expect := testdata.Changelogs[0]
		gock.New(TestClient.APIURL).
			Post(readme.ChangelogEndpoint).
			Reply(201).
			JSON(testdata.Changelogs[0])
		defer gock.Off()

		// Act
		hidden := true
		create := readme.ChangelogParams{
			Body:   "This is a test changelog",
			Title:  "Some Test",
			Hidden: &hidden,
			Type:   "fixed",
		}
		got, _, err := TestClient.Changelog.Create(create)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected Changelog response")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when called without a title", func(t *testing.T) {
		// Act
		hidden := true
		_, _, err := TestClient.Changelog.Create(readme.ChangelogParams{
			Type:   "added",
			Body:   "Test without a title",
			Hidden: &hidden,
		})

		// Assert
		assert.ErrorContains(t, err, "title must be provided", "it returns the expected error")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when called with an invalid type", func(t *testing.T) {
		// Act
		hidden := true
		_, _, err := TestClient.Changelog.Create(readme.ChangelogParams{
			Title:  "Test Title",
			Type:   "invalid",
			Body:   "Test with invalid type",
			Hidden: &hidden,
		})

		// Assert
		assert.ErrorContains(t, err,
			"type must be added, fixed, improved, deprecated, or removed",
			"it returns the expected error",
		)
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}

func Test_Changelog_Update(t *testing.T) {
	t.Run("when called with valid parameters and API responds with 200", func(t *testing.T) {
		// Arrange
		expect := testdata.Changelogs[0]
		gock.New(TestClient.APIURL).
			Put(readme.ChangelogEndpoint + "/" + expect.Slug).
			Reply(200).
			JSON(testdata.Changelogs[0])
		defer gock.Off()

		// Act
		hidden := true
		create := readme.ChangelogParams{
			Body:   "This is a test changelog",
			Title:  "Some Test",
			Hidden: &hidden,
			Type:   "fixed",
		}
		got, _, err := TestClient.Changelog.Update("some-test", create)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected Changelog response")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when called without a title", func(t *testing.T) {
		// Act
		hidden := true
		_, _, err := TestClient.Changelog.Update("some-test", readme.ChangelogParams{
			Type:   "added",
			Body:   "Test without a title",
			Hidden: &hidden,
		})

		// Assert
		assert.ErrorContains(t, err, "title must be provided", "it returns the expected error")
	})

	t.Run("when called with an invalid type", func(t *testing.T) {
		// Act
		hidden := true
		_, _, err := TestClient.Changelog.Update("some-test", readme.ChangelogParams{
			Title:  "Test Title",
			Type:   "invalid",
			Body:   "Test with invalid type",
			Hidden: &hidden,
		})

		// Assert
		assert.ErrorContains(t, err,
			"type must be added, fixed, improved, deprecated, or removed",
			"it returns the expected error",
		)
	})
}

func Test_Changelog_Delete(t *testing.T) {
	t.Run("when called with valid parameters and API responds with 204", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.ChangelogEndpoint + "/some-test").
			Reply(204)
		defer gock.Off()

		// Act
		got, _, err := TestClient.Changelog.Delete("some-test")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.True(t, got, "it returns true")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("handles API errors", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.ChangelogEndpoint + "/some-test").
			Reply(404).
			JSON(`{"error":"CHANGELOG_NOTFOUND"}`)
		defer gock.Off()

		// Act
		got, _, err := TestClient.Changelog.Delete("some-test")

		// Assert
		assert.ErrorContains(t, err, "API responded with a non-OK status: 404", "it returns the expected error")
		assert.False(t, got, "it returns false")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}
