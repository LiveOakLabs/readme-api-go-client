package readme_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/internal/testutil"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

const categoryEndpoint = "http://readme-test.local/api/v1/categories"

// mockCategory represents a category in a project.
// This is the response from the API when calling Get() and Update() and is
// also referenced by other types of responses throughout the API.
var mockCategory = readme.Category{
	Title:     "Documentation",
	Slug:      "documentation",
	Order:     1,
	Reference: false,
	ID:        "63a77777f52b9f006b6bf212",
	Version:   "63a77777f52b9f006b6bf215",
	Project:   "638cf4cedea3ff0096d1a955",
	CreatedAt: "2022-12-04T19:28:15.240Z",
	Type:      "guide",
}

// mockCategoryVersionCreateResponse represents the response from the API when
// calling Create().
var mockCategoryVersionCreateResponse = readme.CategoryVersionSaved{
	Title:     "Test Category",
	Slug:      "test-category",
	Order:     9999,
	Reference: false,
	ID:        "63b463a1d692480062747ef6",
	Project:   "638cf4cedea3ff0096d1a955",
	Version: readme.CategoryVersion{
		Version:      "1.1.0",
		VersionClean: "1.1.0",
		Codename:     "Some test",
		IsStable:     true,
		IsBeta:       true,
		IsHidden:     false,
		IsDeprecated: false,
		Categories:   []readme.Category{mockCategory},
		ID:           "63a77777f52b9f006b6bf215",
		Project:      "638cf4cedea3ff0096d1a955",
		ReleaseDate:  "2022-12-04T19:28:15.190Z",
		CreatedAt:    "2022-12-04T19:28:15.190Z",
		ForkedFrom: readme.CategoryVersionForkedFrom{
			Version: mockVersion.VersionClean,
			ID:      mockVersion.Version,
		},
	},
	CreatedAt: "2023-01-03T17:19:29.388Z",
	Type:      "guide",
}

func Test_Category_GetAll(t *testing.T) {
	// Arrange
	mockResponse := testutil.APITestResponse{
		URL:     categoryEndpoint + "?perPage=100&page=1",
		Status:  200,
		Body:    fmt.Sprintf("[%s]", testutil.StructToJson(t, mockCategory)),
		Headers: mockPaginatedRequestHeader,
	}
	expect := []readme.Category{mockCategory}
	api := mockResponse.New(t)

	// Act
	got, _, err := api.Category.GetAll(readme.RequestOptions{Page: 1})

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected data")
}

func Test_Category_GetAll_Pagination_Invalid(t *testing.T) {
	var expect []readme.Category

	mockResponse := testutil.APITestResponse{
		URL:    "",
		Status: 200,
		Body:   `[]`,
	}

	t.Run("when pagination header is invalid", func(t *testing.T) {
		// Arrange
		mockResponse.Headers = http.Header{
			"Link":          []string{`</categories?page=16>; rel="next", <>; rel="prev", <>; rel="last"`},
			"X-Total-Count": []string{"90"},
		}
		api := mockResponse.New(t)

		// Act
		got, _, _ := api.Category.GetAll(readme.RequestOptions{PerPage: 6, Page: 15})

		// Asert
		assert.Equal(t, expect, got, "returns nil []Category")
	})

	t.Run("when page >= (totalCount / perPage)", func(t *testing.T) {
		// Arrange
		mockResponse.Headers = http.Header{
			"Link":          []string{`</categories?page=16>; rel="next", <>; rel="prev", <>; rel="last"`},
			"X-Total-Count": []string{"90"},
		}
		api := mockResponse.New(t)

		// Act
		got, apiResponse, _ := api.Category.GetAll(readme.RequestOptions{PerPage: 6, Page: 14})

		// Assert
		assert.Equal(t, "/categories?perPage=6&page=15", apiResponse.Request.Endpoint, "it returns expected endpoint")
		assert.Equal(t, expect, got, "it returns expected []Category response")
	})

	t.Run("when total count header is not a number", func(t *testing.T) {
		// Arrange
		var expect []readme.Category
		mockResponse.Headers = http.Header{
			"Link":          []string{`</categories?page=2>; rel="next", <>; rel="prev", <>; rel="last"`},
			"X-Total-Count": []string{"x"},
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.Category.GetAll(readme.RequestOptions{PerPage: 6, Page: 15})

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "unable to parse 'x-total-count' header", "it returns the expected error")
		assert.Equal(t, expect, got, "it returns nil []Category")
	})
}

func Test_Category_Get(t *testing.T) {
	// Arrange
	mockResponse := testutil.APITestResponse{
		URL:    categoryEndpoint + "/some-test",
		Status: 200,
		Body:   testutil.StructToJson(t, mockCategory),
	}
	api := mockResponse.New(t)

	// Act
	got, _, err := api.Category.Get("some-test", readme.RequestOptions{Version: "1.1.2"})

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, mockCategory, got, "it returns expected []Category struct")
}

func Test_Category_GetDocs(t *testing.T) {
	// Arrange
	var expect []readme.CategoryDocs

	mockResponse := testutil.APITestResponse{
		URL:    categoryEndpoint + "/some-test/docs",
		Status: 200,
		Body: `
			[
				{
					"_id": "63a77777f52b9f006b6bf213",
					"title": "Some test doc",
					"slug": "some-test-doc",
					"order": 999,
					"hidden": true,
					"children": []
				}
			]
		`,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	// Act
	got, _, err := api.Category.GetDocs("some-test", readme.RequestOptions{Version: "1.1.2"})

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected []CategoryDocs struct")
}

func Test_Category_Create(t *testing.T) {
	// Arrange
	mockResponse := testutil.APITestResponse{
		URL:    categoryEndpoint,
		Status: 201,
		Body:   testutil.StructToJson(t, mockCategoryVersionCreateResponse),
	}
	api := mockResponse.New(t)

	t.Run("when called with valid parameters and API responds with 201", func(t *testing.T) {
		// Act
		testCreate := readme.CategoryParams{
			Title: "Test Category",
			Type:  "guide",
		}
		got := &readme.CategoryVersionSaved{}
		_, err := api.Category.Create(got, testCreate, readme.RequestOptions{Version: "1.0.0"})

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, &mockCategoryVersionCreateResponse, got, "it returns expected CategorySaved struct")
	})

	t.Run("when type is invalid", func(t *testing.T) {
		// Act
		testCreate := readme.CategoryParams{
			Title: "Test Category",
			Type:  "invalid",
		}
		got := &readme.CategoryVersionSaved{}
		_, err := api.Category.Create(got, testCreate)

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "type must be 'guide' or 'reference'", "it returns the expected error")
	})
}

func Test_Category_Update(t *testing.T) {
	t.Run("when called with valid parameters and API responds with 200", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    categoryEndpoint + "/some-test",
			Status: 200,
			Body:   testutil.StructToJson(t, mockCategory),
		}
		api := mockResponse.New(t)

		// Act
		testUpdate := readme.CategoryParams{
			Title: "Test Category",
			Type:  "guide",
		}
		got, _, err := api.Category.Update("some-test", testUpdate, readme.RequestOptions{Version: "1.1.0"})

		// Assert
		assert.NoError(t, err, "it returns no errors")
		assert.Equal(t, mockCategory, got, "it returns expected Category struct")
	})

	t.Run("when type is invalid", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    categoryEndpoint + "/some-test",
			Status: 200,
			Body:   testutil.StructToJson(t, mockCategory),
		}
		api := mockResponse.New(t)
		// Act
		testUpdate := readme.CategoryParams{
			Title: "Test Category",
			Type:  "invalid",
		}
		_, _, err := api.Category.Update("some-test", testUpdate)

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "type must be 'guide' or 'reference'", "it returns the expected error")
	})
}

func Test_Category_Delete(t *testing.T) {
	t.Run("when called with valid parameters and API responds with 200", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    categoryEndpoint + "/some-test",
			Status: 204,
			Body:   "{}",
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.Category.Delete("some-test")

		// Assert
		assert.NoError(t, err, "it returns no errors")
		assert.True(t, got, "it returns expected CustomPage struct")
	})

	t.Run("when API responds with error", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    categoryEndpoint + "/some-test",
			Status: 400,
			Body:   "{}",
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.Category.Delete("some-test", readme.RequestOptions{Version: "1.1.0"})

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400", "it returns the expected error")
		assert.False(t, got, "it returns false")
	})
}
