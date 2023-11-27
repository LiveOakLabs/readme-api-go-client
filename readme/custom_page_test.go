package readme_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_CustomPages_GetAll(t *testing.T) {
	t.Run("when called with valid parameters and API responds with 200", func(t *testing.T) {
		// Arrange
		expect := testdata.CustomPages
		gock.New(TestClient.APIURL).
			Get(readme.CustomPageEndpoint).
			Reply(200).
			AddHeader("Link", `</custompages?page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "3").
			JSON(expect)
		defer gock.Off()

		reqOpts := readme.RequestOptions{PerPage: 100, Page: 1}

		// Act
		got, _, err := TestClient.CustomPage.GetAll(reqOpts)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected []CustomPage struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when response pagination header is invalid", func(t *testing.T) {
		t.Run("when page >= (totalCount / perPage)", func(t *testing.T) {
			// Arrange
			expect := testdata.CustomPages
			gock.New(TestClient.APIURL).
				Get(readme.CustomPageEndpoint).
				MatchParam("page", "16").
				MatchParam("perPage", "6").
				Reply(200).
				AddHeader("Link", `</custompages?perPage=6&page=15>; rel="next", <>; rel="prev", <>; rel="last"`).
				AddHeader("x-total-count", "90").
				JSON(expect)
			defer gock.Off()

			reqOpts := readme.RequestOptions{PerPage: 6, Page: 16}

			// Act
			got, apiResponse, err := TestClient.CustomPage.GetAll(reqOpts)

			// Assert
			assert.NoError(t, err, "it does not return an error")
			assert.Equal(t, "/custompages?perPage=6&page=16",
				apiResponse.Request.Endpoint,
				"it requests with the expected pagination query parameters")
			assert.Equal(t, expect, got, "it returns expected []CustomPage")
			assert.True(t, gock.IsDone(), "it makes the expected API call")
		})

		t.Run("when pagination x-total-count header is invalid", func(t *testing.T) {
			// Arrange
			var expect []readme.CustomPage
			gock.New(TestClient.APIURL).
				Get(readme.CustomPageEndpoint).
				MatchParam("page", "15").
				MatchParam("perPage", "6").
				Reply(200).
				AddHeader("Link", `</custompages?perPage=6&page=15>; rel="next", <>; rel="prev", <>; rel="last"`).
				AddHeader("x-total-count", "x").
				JSON(expect)
			defer gock.Off()

			reqOpts := readme.RequestOptions{PerPage: 6, Page: 15}

			// Act
			got, _, err := TestClient.CustomPage.GetAll(reqOpts)

			// Assert
			assert.Error(t, err, "it returns an error")
			assert.ErrorContains(t, err, "unable to parse 'x-total-count' header:",
				"it returns the expected error")
			assert.Equal(t, expect, got, "it returns nil []CustomPage")
			assert.True(t, gock.IsDone(), "it makes the expected API call")
		})
	})
}

func Test_CustomPages_Get(t *testing.T) {
	// Arrange
	expect := testdata.CustomPages[0]
	gock.New(TestClient.APIURL).
		Get(readme.CustomPageEndpoint + "/" + expect.Slug).
		Reply(200).
		JSON(expect)
	defer gock.Off()

	// Act
	got, _, err := TestClient.CustomPage.Get(expect.Slug)

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected []CustomPage struct")
	assert.True(t, gock.IsDone(), "it makes the expected API call")
}

func Test_CustomPages_Create(t *testing.T) {
	// Arrange
	expect := testdata.CustomPages[0]
	gock.New(TestClient.APIURL).
		Post(readme.CustomPageEndpoint).
		Reply(201).
		JSON(expect)
	defer gock.Off()

	createParams := readme.CustomPageParams{
		Body:     expect.Body,
		Title:    expect.Title,
		HTML:     expect.HTML,
		Hidden:   &expect.Hidden,
		HTMLMode: &expect.HTMLMode,
	}

	// Act
	got, _, err := TestClient.CustomPage.Create(createParams)

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected CustomPage struct")
	assert.True(t, gock.IsDone(), "it makes the expected API call")
}

func Test_CustomPages_Update(t *testing.T) {
	// Arrange
	expect := testdata.CustomPages[0]

	gock.New(TestClient.APIURL).
		Put(readme.CustomPageEndpoint + "/" + expect.Slug).
		Reply(200).
		JSON(expect)
	defer gock.Off()

	updateParams := readme.CustomPageParams{
		Body:     expect.Body,
		Title:    expect.Title,
		HTML:     expect.HTML,
		Hidden:   &expect.Hidden,
		HTMLMode: &expect.HTMLMode,
	}

	// Act
	got, _, err := TestClient.CustomPage.Update(expect.Slug, updateParams)

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected CustomPage struct")
	assert.True(t, gock.IsDone(), "it makes the expected API call")
}

func Test_CustomPages_Delete(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.CustomPageEndpoint + "/foo").
			Reply(204).
			JSON("{}")
		defer gock.Off()

		// Act
		got, _, err := TestClient.CustomPage.Delete("foo")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.True(t, got, "it returns expected CustomPage struct")
	})

	t.Run("when API responds with error", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.CustomPageEndpoint + "/foo").
			Reply(400).
			JSON("{}")
		defer gock.Off()

		// Act
		got, _, err := TestClient.CustomPage.Delete("foo")

		// Assert
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400",
			"it return the expected error")
		assert.False(t, got, "it returns false")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}
