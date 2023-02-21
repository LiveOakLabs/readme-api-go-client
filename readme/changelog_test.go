package readme_test

import (
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/internal/testutil"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

const changelogTestEndpoint = "http://readme-test.local/api/v1/changelogs"

var mockChangelogBody string = `
	{
		"metadata": {
			"image": [],
			"title": "",
			"description": ""
		},
		"algolia": {
			"pendingAlgoliaPublish": false,
			"recordCount": 0,
			"updatedAt": "2023-01-03T00:29:22.169Z"
		},
		"title": "Some Test",
		"slug": "some-test",
		"body": "This is a test changelog",
		"type": "added",
		"hidden": true,
		"revision": 2,
		"_id": "63b376e244ed08009d672b11",
		"createdAt": "2023-01-03T00:29:22.169Z",
		"updatedAt": "2023-01-03T00:29:22.169Z",
		"__v": 0,
		"html": "<div class=\"magic-block-textarea\"><p>This is a test changelog</p>\n\n</div>"
	}
	`

func Test_Changelog_Get(t *testing.T) {
	// Arrange
	var expect readme.Changelog

	mockResponse := testutil.APITestResponse{
		URL:    changelogTestEndpoint + "/some-test",
		Status: 200,
		Body:   mockChangelogBody,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	// Act
	got, _, err := api.Changelog.Get("some-test")

	// Assert
	assert.NoError(t, err, "it does not returns an error")
	assert.Equal(t, expect, got, "it returns a Changelog struct")
}

func Test_Changelog_GetAll(t *testing.T) {
	// Arrange
	var expect []readme.Changelog

	mockResponse := testutil.APITestResponse{
		URL:    changelogTestEndpoint,
		Status: 200,
		Body:   "[" + mockChangelogBody + "]",
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	// Act
	got, _, err := api.Changelog.GetAll(readme.RequestOptions{Page: 1})

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns slice of Changelog structs")
}

func Test_Changelog_Create(t *testing.T) {
	// Arrange
	var expect readme.Changelog

	mockResponse := testutil.APITestResponse{
		URL:    changelogTestEndpoint,
		Status: 201,
		Body:   mockChangelogBody,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	t.Run("when called with valid parameters and API responds with 201", func(t *testing.T) {
		// Act
		hidden := true
		create := readme.ChangelogParams{
			Body:   "This is a test changelog",
			Title:  "Some Test",
			Hidden: &hidden,
			Type:   "fixed",
		}
		got, _, err := api.Changelog.Create(create)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected Changelog response")
	})

	t.Run("when called without a title", func(t *testing.T) {
		// Act
		hidden := true
		_, _, err := api.Changelog.Create(readme.ChangelogParams{
			Type:   "added",
			Body:   "Test without a title",
			Hidden: &hidden,
		})

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "title must be provided", "it returns the expected error")
	})

	t.Run("when called with an invalid type", func(t *testing.T) {
		// Act
		hidden := true
		_, _, err := api.Changelog.Create(readme.ChangelogParams{
			Title:  "Test Title",
			Type:   "invalid",
			Body:   "Test with invalid type",
			Hidden: &hidden,
		})

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(
			t,
			err,
			"type must be added, fixed, improved, deprecated, or removed",
			"it returns the expected error",
		)
	})
}

func Test_Changelog_Update(t *testing.T) {
	// Arrange
	var expect readme.Changelog

	mockResponse := testutil.APITestResponse{
		URL:    changelogTestEndpoint + "/some-test",
		Status: 200,
		Body:   mockChangelogBody,
	}
	testutil.JsonToStruct(t, mockChangelogBody, &expect)
	api := mockResponse.New(t)

	t.Run("when called with valid parameters and API responds with 200", func(t *testing.T) {
		// Act
		hidden := true
		create := readme.ChangelogParams{
			Body:   "This is a test changelog",
			Title:  "Some Test",
			Hidden: &hidden,
			Type:   "fixed",
		}
		got, _, err := api.Changelog.Update("some-test", create)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected Changelog response")
	})

	t.Run("when called without a title", func(t *testing.T) {
		// Act
		hidden := true
		_, _, err := api.Changelog.Update("some-test", readme.ChangelogParams{
			Type:   "added",
			Body:   "Test without a title",
			Hidden: &hidden,
		})

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "title must be provided", "it returns the expected error")
	})

	t.Run("when called with an invalid type", func(t *testing.T) {
		// Act
		hidden := true
		_, _, err := api.Changelog.Update("some-test", readme.ChangelogParams{
			Title:  "Test Title",
			Type:   "invalid",
			Body:   "Test with invalid type",
			Hidden: &hidden,
		})

		// Assert
		assert.Error(t, err, "it returns error")
		assert.ErrorContains(
			t,
			err,
			"type must be added, fixed, improved, deprecated, or removed",
			"it returns the expected error",
		)
	})
}

func Test_Changelog_Delete(t *testing.T) {
	// Arrange
	mockResponse := testutil.APITestResponse{
		URL:    changelogTestEndpoint + "/some-test",
		Status: 204,
		Body:   "",
	}
	api := mockResponse.New(t)

	t.Run("when called with valid parameters and API responds with 204", func(t *testing.T) {
		// Act
		got, _, err := api.Changelog.Delete("some-test")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.True(t, got, "it returns true")
	})

	t.Run("handles API errors", func(t *testing.T) {
		mockResponse := testutil.APITestResponse{
			URL:    changelogTestEndpoint + "/some-test",
			Status: 404,
			Body:   `{"error":"CHANGELOG_NOTFOUND"}`,
		}
		api := mockResponse.New(t)
		got, _, err := api.Changelog.Delete("some-test")
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 404", "it returns the expected error")
		assert.False(t, got, "it returns false")
	})
}
