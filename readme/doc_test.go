package readme_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_Docs_Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// Arrange
		expect := testdata.Docs[0]
		gock.New(TestClient.APIURL).
			Get(readme.DocEndpoint + "/" + expect.Slug).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.Doc.Get(expect.Slug)

		// Assert
		assert.Equal(t, expect, got, "it returns expected Doc struct")
		assert.NoError(t, err, "it does not return an error")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("returns a production doc when requested", func(t *testing.T) {
		// Arrange
		expect := testdata.Docs[0]
		gock.New(TestClient.APIURL).
			Get(readme.DocEndpoint + "/" + expect.Slug + "/production").
			Reply(200).
			JSON(expect)
		defer gock.Off()

		reqOpts := readme.RequestOptions{ProductionDoc: true}

		// Act
		got, _, _ := TestClient.Doc.Get(expect.Slug, reqOpts)

		// Assert
		assert.Equal(t, expect, got)
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}

func Test_Doc_Create(t *testing.T) {
	expect := testdata.Docs[0]
	createParams := readme.DocParams{
		Body:          expect.Body,
		Category:      expect.Category,
		Error:         expect.Error,
		Hidden:        &expect.Hidden,
		Order:         &expect.Order,
		ParentDocSlug: "foo",
		Title:         expect.Title,
		Type:          expect.Type,
	}

	t.Run("happy path creation", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Post(readme.DocEndpoint).
			Reply(201).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.Doc.Create(createParams)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected Doc struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("handles invalid parameters", func(t *testing.T) {
		// Arrange
		_, _, err := TestClient.Doc.Create(readme.DocParams{})
		assert.Error(t, err, "returns error when title is empty")

		_, _, err = TestClient.Doc.Create(readme.DocParams{Title: "testing"})
		assert.Error(t, err, "returns error when category is empty")
	})

	t.Run("uses provided RequestOptions", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Post(readme.DocEndpoint).
			Reply(201).
			JSON(expect)
		defer gock.Off()

		reqOpts := readme.RequestOptions{Version: "1.1.0"}

		// Act
		got, _, _ := TestClient.Doc.Create(createParams, reqOpts)

		// Assert
		assert.Equal(t, expect, got, "it returns expected Doc struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}

func Test_Doc_Update(t *testing.T) {
	expect := testdata.Docs[0]
	create := readme.DocParams{
		Body:     expect.Body,
		Title:    expect.Title,
		Type:     expect.Type,
		Category: expect.Category,
		Hidden:   &expect.Hidden,
		Order:    &expect.Order,
	}
	t.Run("happy path update", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Put(readme.DocEndpoint + "/" + expect.Slug).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.Doc.Update(expect.Slug, create)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expect, got)
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("uses provided RequestOptions", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Put(readme.DocEndpoint + "/" + expect.Slug).
			Reply(200).
			JSON(expect)
		defer gock.Off()

		reqOpts := readme.RequestOptions{Version: "1.1.0"}

		// Act
		got, _, err := TestClient.Doc.Update(expect.Slug, create, reqOpts)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expect, got)
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("handles invalid parameters", func(t *testing.T) {
		_, _, err := TestClient.Doc.Update(expect.Slug, readme.DocParams{})
		assert.Error(t, err, "returns error when title is empty")

		updateParams := readme.DocParams{Title: "updated title"}
		_, _, err = TestClient.Doc.Update(expect.Slug, updateParams)
		assert.Error(t, err, "returns error when category is empty")
	})
}

func Test_Doc_Delete(t *testing.T) {
	t.Run("successfully deletes a doc", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.DocEndpoint + "/some-test-doc").
			Reply(204).
			JSON("{}")
		defer gock.Off()

		reqOpts := readme.RequestOptions{Version: "1.1.0"}

		// Act
		got, _, err := TestClient.Doc.Delete("some-test-doc", reqOpts)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.True(t, got, "it returns true")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("returns error when API response has error", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.DocEndpoint + "/some-test-doc").
			Reply(400).
			JSON(readme.APIErrorResponse{Error: "INVALID"})
		defer gock.Off()

		reqOpts := readme.RequestOptions{Version: "1.1.0"}
		expectErr := "API responded with a non-OK status: 400"

		// Act
		got, _, err := TestClient.Doc.Delete("some-test-doc", reqOpts)

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, expectErr, "it returns the expected error")
		assert.False(t, got, "it returns false")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}

func Test_Doc_Search(t *testing.T) {
	t.Run("returns expected result when searching and API responds with 200", func(t *testing.T) {
		// Arrange
		expect := testdata.DocSearchResults
		gock.New(TestClient.APIURL).
			Post(readme.DocEndpoint+"/search").
			MatchParam("search", "turtle").
			Reply(200).
			JSON(expect)
		defer gock.Off()
		// Act
		got, _, err := TestClient.Doc.Search("turtle")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expect.Results, got)
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("returns expected result when version is specified", func(t *testing.T) {
		// Arrange
		expect := testdata.DocSearchResults
		gock.New(TestClient.APIURL).
			Post(readme.DocEndpoint + "/search").
			Reply(200).
			JSON(expect)
		defer gock.Off()

		requestOpts := readme.RequestOptions{Version: "1.1.0"}

		// Act
		got, _, err := TestClient.Doc.Search("turtle", requestOpts)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expect.Results, got)
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("returns error when API response has error", func(t *testing.T) {
		// Arrange
		expect := readme.APIErrorResponse{
			Error:   "INVALID",
			Message: "This is an error",
		}
		gock.New(TestClient.APIURL).
			Post(readme.DocEndpoint + "/search").
			Reply(400).
			JSON(expect)
		defer gock.Off()

		expectErr := "API responded with a non-OK status: 400"

		// Act
		_, _, err := TestClient.Doc.Search("something")

		// Assert
		assert.ErrorContains(t, err, expectErr, "it returns the expected error")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}
