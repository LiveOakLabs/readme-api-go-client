package readme_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_Category_Get(t *testing.T) {
	// Arrange
	expect := testdata.Categories[0]
	gock.New(TestClient.APIURL).
		Get(readme.CategoryEndpoint + "/" + expect.Slug).
		Reply(200).
		JSON(expect)
	defer gock.Off()

	// Act
	got, _, err := TestClient.Category.Get(expect.Slug, readme.RequestOptions{Version: "1.1.2"})

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected []Category struct")
	assert.True(t, gock.IsDone(), "it makes the expected API call")
}

func Test_Category_GetDocs(t *testing.T) {
	// Arrange
	expect := testdata.CategoryDocs
	gock.New(TestClient.APIURL).
		Get(readme.CategoryEndpoint + "/some-test/docs").
		Reply(200).
		JSON(expect)
	defer gock.Off()

	// Act
	got, _, err := TestClient.Category.GetDocs("some-test", readme.RequestOptions{Version: "1.1.2"})

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected []CategoryDocs struct")
	assert.True(t, gock.IsDone(), "it makes the expected API call")
}

func Test_Category_Create(t *testing.T) {
	t.Run("when called with valid parameters and API responds with 201", func(t *testing.T) {
		// Arrange
		expect := &testdata.CategoryVersionSaved
		gock.New(TestClient.APIURL).
			Post(readme.CategoryEndpoint).
			Reply(201).
			JSON(expect)
		defer gock.Off()

		createOpts := readme.CategoryParams{
			Title: "Test Category",
			Type:  "guide",
		}
		reqOpts := readme.RequestOptions{Version: "1.0.0"}

		// Act
		got := &readme.CategoryVersionSaved{}
		_, err := TestClient.Category.Create(got, createOpts, reqOpts)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected CategoryVersionSaved struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when type is invalid", func(t *testing.T) {
		// Act
		createOpts := readme.CategoryParams{
			Title: "Test Category",
			Type:  "invalid",
		}
		got := &readme.CategoryVersionSaved{}
		_, err := TestClient.Category.Create(got, createOpts)

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "type must be 'guide' or 'reference'",
			"it returns the expected error")
	})
}

func Test_Category_Update(t *testing.T) {
	t.Run("when called with valid parameters and API responds with 200", func(t *testing.T) {
		// Arrange
		expect := testdata.Categories[0]
		gock.New(TestClient.APIURL).
			Put(readme.CategoryEndpoint + "/some-test").
			Reply(200).
			JSON(expect)
		defer gock.Off()

		updateParams := readme.CategoryParams{
			Title: "Test Category",
			Type:  "guide",
		}
		reqOpts := readme.RequestOptions{Version: "1.1.0"}

		// Act
		got, _, err := TestClient.Category.Update("some-test", updateParams, reqOpts)

		// Assert
		assert.NoError(t, err, "it returns no errors")
		assert.Equal(t, expect, got, "it returns expected Category struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when type is invalid", func(t *testing.T) {
		// Arrange
		updateParams := readme.CategoryParams{
			Title: "Test Category",
			Type:  "invalid",
		}

		// Act
		_, _, err := TestClient.Category.Update("some-test", updateParams)

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "type must be 'guide' or 'reference'", "it returns the expected error")
	})
}

func Test_Category_Delete(t *testing.T) {
	t.Run("when called with valid parameters and API responds with 200", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.CategoryEndpoint + "/some-test").
			Reply(204).
			JSON("{}")
		defer gock.Off()

		// Act
		got, _, err := TestClient.Category.Delete("some-test")

		// Assert
		assert.NoError(t, err, "it returns no errors")
		assert.True(t, got, "it returns expected CustomPage struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when API responds with error", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.CategoryEndpoint + "/some-test").
			Reply(400).
			JSON("{}")
		defer gock.Off()

		reqOpts := readme.RequestOptions{Version: "1.1.0"}

		// Act
		got, _, err := TestClient.Category.Delete("some-test", reqOpts)

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "ReadMe API Error: 400 on DELETE", "it returns the expected error")
		assert.False(t, got, "it returns false")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}
