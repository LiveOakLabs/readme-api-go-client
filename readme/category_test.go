package readme_test

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_Category_Get(t *testing.T) {
	tc := []struct {
		name             string
		input            string
		slug             string
		mockResponse     interface{}
		mockStatus       int
		expect           readme.Category
		expectErr        bool
		expectErrMsg     string
		setupGetAll      bool
		mockGetAllResult []readme.Category
		mockAllStatus    int
	}{
		{
			name:         "valid slug",
			input:        testdata.Categories[0].Slug,
			mockResponse: testdata.Categories[0],
			mockStatus:   200,
			expect:       testdata.Categories[0],
			expectErr:    false,
			slug:         testdata.Categories[0].Slug,
		},
		{
			name:             "valid ID",
			input:            "id:" + testdata.Categories[0].ID,
			mockResponse:     testdata.Categories[0],
			mockStatus:       200,
			expect:           testdata.Categories[0],
			expectErr:        false,
			setupGetAll:      true,
			mockGetAllResult: testdata.Categories,
			slug:             testdata.Categories[0].Slug,
		},
		{
			name:  "invalid slug",
			input: "id:invalid-slug",
			mockResponse: readme.APIResponse{
				APIErrorResponse: readme.APIErrorResponse{
					Error:      "CATEGORY_NOTFOUND",
					Message:    "The category with the slug 'on-existent-id' couldn't be found.",
					Suggestion: "You can list all available categories using GET /api/v1/categories.",
				},
			},
			mockStatus:   404,
			expectErr:    true,
			expectErrMsg: "ReadMe API Error: 404 on GET",
		},
		{
			name:         "invalid ID, no match found",
			input:        "id:" + testdata.Categories[0].ID,
			slug:         testdata.Categories[0].Slug,
			expectErr:    true,
			expectErrMsg: "a category slug or id must be provided",
			setupGetAll:  true,
			mockGetAllResult: []readme.Category{
				testdata.Categories[1], // Does not match the ID
			},
		},
		{
			name:         "API returns error",
			input:        testdata.Categories[0].Slug,
			slug:         testdata.Categories[0].Slug,
			mockResponse: readme.APIErrorResponse{Error: "INTERNAL_SERVER_ERROR"},
			mockStatus:   500,
			expectErr:    true,
			expectErrMsg: "ReadMe API Error: 500 on GET",
		},
		{
			name:          "error retrieving all categories",
			input:         "id:" + testdata.Categories[0].ID,
			expectErr:     true,
			expectErrMsg:  "ReadMe API Error: 500 on GET",
			setupGetAll:   true,
			mockAllStatus: 500,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the GetAll request if necessary
			if tt.setupGetAll {
				mockStatus := 200
				if tt.mockAllStatus != 0 {
					mockStatus = tt.mockAllStatus
				}
				gock.New(TestClient.APIURL).
					Get(readme.CategoryEndpoint).
					Reply(mockStatus).
					AddHeader("x-total-count", fmt.Sprintf("%d", len(tt.mockGetAllResult))).
					AddHeader("link", `<>; rel="next", <>; rel="prev", </categories?perPage=100&page=1>; rel="last"`).
					JSON(tt.mockGetAllResult)
			}
			// Mock the Get request
			if tt.mockResponse != nil {
				gock.New(TestClient.APIURL).
					Get(readme.CategoryEndpoint + "/" + tt.slug).
					Reply(tt.mockStatus).
					JSON(tt.mockResponse)
			}
			defer gock.Off()

			// Act
			got, _, err := TestClient.Category.Get(tt.input, readme.RequestOptions{Version: "1.1.2"})

			// Assert
			if tt.expectErr {
				assert.Error(t, err, "it returns an error")
				assert.Contains(t, err.Error(), tt.expectErrMsg, "it returns the expected error message")
			} else {
				assert.NoError(t, err, "it does not return an error")
				assert.Equal(t, tt.expect, got, "it returns expected Category struct")
			}
			assert.True(t, gock.IsDone(), "it makes the expected API call")
		})
	}
}

func Test_Category_GetAll(t *testing.T) {
	tc := []struct {
		name          string
		expect        []readme.Category
		expectErr     bool
		expectErrMsg  string
		mockStatus    int
		mockHeaders   map[string]string
		mockResponse  any
		expectedPages int
		expectedTotal int
	}{
		{
			name:         "when API responds with 200 and single page",
			expect:       testdata.Categories,
			expectErr:    false,
			expectErrMsg: "",
			mockStatus:   200,
			mockHeaders: map[string]string{
				"x-total-count": "1",
				"link":          `</categories?perPage=100&page=1>; rel="next", </categories?perPage=100&page=1>; rel="prev", </categories?perPage=100&page=1>; rel="last"`,
			},
			mockResponse:  testdata.Categories,
			expectedPages: 1,
			expectedTotal: len(testdata.Categories),
		},
		{
			name:         "when API responds with 400",
			expect:       nil,
			expectErr:    true,
			expectErrMsg: "ReadMe API Error: 400 on GET",
			mockStatus:   400,
			mockHeaders:  nil,
			mockResponse: nil,
		},
		{
			name:         "when API responds with 200 and no categories",
			expect:       nil,
			expectErr:    false,
			expectErrMsg: "",
			mockStatus:   200,
			mockHeaders: map[string]string{
				"x-total-count": "0",
				"Link":          `</categories?perPage=100&page=1>; rel="next", </categories?perPage=100&page=1>; rel="prev", </categories?perPage=100&page=1>; rel="last"`,
			},
			mockResponse:  nil,
			expectedPages: 1,
			expectedTotal: 0,
		},
		{
			name:         "when API paginates across multiple pages",
			expect:       []readme.Category{testdata.Categories[0], testdata.Categories[1]},
			expectErr:    false,
			expectErrMsg: "",
			mockStatus:   200,
			mockHeaders: map[string]string{
				"x-total-count": "2",
				"Link":          `</categories?perPage=1&page=2>; rel="next", <>; rel="prev", </categories?perPage=1&page=2>; rel="last"`,
			},
			mockResponse:  []readme.Category{testdata.Categories[0]},
			expectedPages: 2,
			expectedTotal: 2,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			gock.New(TestClient.APIURL).
				Get(readme.CategoryEndpoint).
				MatchParam("page", "1").
				Reply(tt.mockStatus).
				SetHeaders(tt.mockHeaders).
				JSON(tt.mockResponse)
			defer gock.OffAll()

			// If the test involves pagination, we need to set up the second request
			if tt.name == "when API paginates across multiple pages" {
				gock.New(TestClient.APIURL).
					Get(readme.CategoryEndpoint).
					MatchParam("page", "2").
					Reply(tt.mockStatus).
					SetHeaders(map[string]string{
						"x-total-count": "2",
						"Link":          `</categories?perPage=1&page=1>; rel="prev", <>; rel="next", </categories?perPage=1&page=2>; rel="last"`,
					}).
					JSON([]readme.Category{testdata.Categories[1]})
			}

			// Act
			opts := readme.RequestOptions{
				Version: "1.0.0",
				PerPage: 1,
			}
			got, _, err := TestClient.Category.GetAll(opts)

			// Assert
			if tt.expectErr {
				assert.Error(t, err, "it returns an error")
				assert.ErrorContains(t, err, tt.expectErrMsg, "it returns the expected error")
			} else {
				assert.NoError(t, err, "it does not return an error")
				assert.Equal(t, tt.expect, got, "it returns expected []Category struct")
			}

			assert.True(t, gock.IsDone(), "it makes the expected API call")
		})
	}
}

func Test_Category_GetDocs(t *testing.T) {
	// Arrange
	expect := testdata.CategoryDocs
	gock.New(TestClient.APIURL).
		Get(readme.CategoryEndpoint + "/some-test/docs").
		Reply(200).
		JSON(expect)
	defer gock.OffAll()

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
		defer gock.OffAll()

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
		defer gock.OffAll()

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
		defer gock.OffAll()

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
		defer gock.OffAll()

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
