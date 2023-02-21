package readme_test

import (
	"net/http"
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/internal/testutil"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

const customPagesEndpoint = "http://readme-test.local/api/v1/custompages"

// mockCustomPageBody represents the raw JSON response body for a custom page from the ReadMe.com API.
var mockCustomPageBody string = `
	{
		"metadata": {
		"image": [],
		"title": "",
		"description": ""
		},
		"title": "Some Test Page",
		"slug": "some-test-page",
		"body": "# A test page.",
		"html": "",
		"htmlmode": false,
		"fullscreen": false,
		"hidden": true,
		"pendingAlgoliaPublish": false,
		"revision": 2,
		"_id": "63ae563ec4d05f018b26cf18",
		"createdAt": "2022-12-30T03:08:46.695Z",
		"updatedAt": "2022-12-30T03:09:46.695Z",
		"__v": 0
	}
	`

func Test_CustomPages_GetAll(t *testing.T) {
	t.Run("when called with valid parameters and API responds with 200", func(t *testing.T) {
		// Arrange
		var expect []readme.CustomPage

		mockRequestHeader := http.Header{
			"Link":          []string{`<>; rel="next", <>; rel="prev", <>; rel="last"`},
			"X-Total-Count": []string{"1"},
		}
		mockResponse := testutil.APITestResponse{
			URL:     customPagesEndpoint + "?perPage=100&page=1",
			Status:  200,
			Body:    "[" + mockCustomPageBody + "]",
			Headers: mockRequestHeader,
		}
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		got, _, err := api.CustomPage.GetAll(readme.RequestOptions{PerPage: 100, Page: 1})

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected []CustomPage struct")
	})

	t.Run("when response pagination header is invalid", func(t *testing.T) {
		var expect []readme.CustomPage

		t.Run("when page >= (totalCount / perPage)", func(t *testing.T) {
			// Arrange
			mockResponse := testutil.APITestResponse{
				URL:    "",
				Status: 200,
				Body:   `[]`,
				Headers: http.Header{
					"Link":          []string{`</custompages?page=16>; rel="next", <>; rel="prev", <>; rel="last"`},
					"X-Total-Count": []string{"90"},
				},
			}
			api := mockResponse.New(t)

			// Act
			got, apiResponse, err := api.CustomPage.GetAll(readme.RequestOptions{PerPage: 6, Page: 14})

			// Assert
			assert.NoError(t, err, "it does not return an error")
			assert.Equal(t, "/custompages?perPage=6&page=15", apiResponse.Request.Endpoint,
				"it requests with the expected pagination query parameters")
			assert.Equal(t, expect, got, "it returns expected []CustomPage")
		})

		t.Run("when pagination x-total-count header is invalid", func(t *testing.T) {
			// Arrange
			mockResponse := testutil.APITestResponse{
				URL:    "",
				Status: 200,
				Body:   `[]`,
				Headers: http.Header{
					"Link":          []string{`</custompages?page=2>; rel="next", <>; rel="prev", <>; rel="last"`},
					"X-Total-Count": []string{"x"},
				},
			}
			api := mockResponse.New(t)

			// Act
			got, _, err := api.CustomPage.GetAll(readme.RequestOptions{PerPage: 6, Page: 15})

			// Assert
			assert.Error(t, err, "it returns an error")
			assert.ErrorContains(t, err, "unable to parse 'x-total-count' header:", "it returns the expected error")
			assert.Equal(t, expect, got, "it returns nil []CustomPage")
		})
	})
}

func Test_CustomPages_Get(t *testing.T) {
	// Arrange
	expect := readme.CustomPage{}

	mockResponse := testutil.APITestResponse{
		URL:    customPagesEndpoint + "/some-test-page",
		Status: 200,
		Body:   mockCustomPageBody,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	// Act
	got, _, err := api.CustomPage.Get("some-test-page")

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected []CustomPage struct")
}

func Test_CustomPages_Create(t *testing.T) {
	// Arrange
	expect := readme.CustomPage{}

	mockResponse := testutil.APITestResponse{
		URL:    customPagesEndpoint,
		Status: 201,
		Body:   mockCustomPageBody,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	// Act
	hidden := true
	htmlMode := false
	testCreate := readme.CustomPageParams{
		Body:     "this is a test",
		Title:    "Test Page",
		HTML:     "&nbsp;",
		Hidden:   &hidden,
		HTMLMode: &htmlMode,
	}
	got, _, err := api.CustomPage.Create(testCreate)

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected CustomPage struct")
}

func Test_CustomPages_Update(t *testing.T) {
	// Arrange
	expect := readme.CustomPage{}

	mockResponse := testutil.APITestResponse{
		URL:    customPagesEndpoint + "/foo",
		Status: 200,
		Body:   mockCustomPageBody,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	// Act
	hidden := true
	htmlMode := false
	testCreate := readme.CustomPageParams{
		Body:     "this is a test",
		Title:    "Test Page",
		HTML:     "&nbsp;",
		Hidden:   &hidden,
		HTMLMode: &htmlMode,
	}
	got, _, err := api.CustomPage.Update("foo", testCreate)

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected CustomPage struct")
}

func Test_CustomPages_Delete(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    customPagesEndpoint + "/foo",
			Status: 204,
			Body:   "{}",
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.CustomPage.Delete("foo")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.True(t, got, "it returns expected CustomPage struct")
	})

	t.Run("when API responds with error", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    customPagesEndpoint + "/foo",
			Status: 400,
			Body:   "{}",
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.CustomPage.Delete("foo")

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400", "it return the expected error")
		assert.False(t, got, "it returns false")
	})
}
