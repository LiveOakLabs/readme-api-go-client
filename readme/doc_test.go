package readme_test

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_Doc_Get(t *testing.T) {
	tc := []struct {
		name             string
		input            string
		slug             string
		mockResponse     interface{}
		mockStatus       int
		expect           readme.Doc
		expectErr        bool
		expectErrMsg     string
		setupSearch      bool
		mockSearchResult readme.DocSearchResults
		mockSearchStatus int
		options          readme.RequestOptions
	}{
		{
			name:         "get by slug",
			input:        testdata.Docs[0].Slug,
			slug:         testdata.Docs[0].Slug,
			mockResponse: testdata.Docs[0],
			mockStatus:   200,
			expect:       testdata.Docs[0],
			expectErr:    false,
		},
		{
			name:             "get by ID",
			input:            "id:" + testdata.Docs[0].ID,
			slug:             testdata.Docs[0].Slug,
			mockResponse:     testdata.Docs[0],
			mockStatus:       200,
			expect:           testdata.Docs[0],
			expectErr:        false,
			setupSearch:      true,
			mockSearchResult: testdata.DocSearchResults,
		},
		{
			name:  "invalid ID, no match found",
			input: "id:" + testdata.Docs[0].ID,
			mockSearchResult: readme.DocSearchResults{
				Results: []readme.DocSearchResult{},
			},
			mockSearchStatus: 200,
			expectErr:        true,
			expectErrMsg:     fmt.Sprintf("no doc found matching id %s (is it hidden?)", testdata.Docs[0].ID),
			setupSearch:      true,
		},
		{
			name:         "API returns error",
			input:        testdata.Docs[0].Slug,
			slug:         testdata.Docs[0].Slug,
			mockResponse: readme.APIErrorResponse{Error: "INTERNAL_SERVER_ERROR"},
			mockStatus:   500,
			expectErr:    true,
			expectErrMsg: "ReadMe API Error: 500 on GET",
		},
		{
			name:         "get production doc",
			input:        testdata.Docs[0].Slug,
			slug:         testdata.Docs[0].Slug,
			mockResponse: testdata.Docs[0],
			mockStatus:   200,
			expect:       testdata.Docs[0],
			expectErr:    false,
			options:      readme.RequestOptions{ProductionDoc: true},
		},
		{
			name:             "error searching all docs",
			input:            "id:" + testdata.Docs[0].ID,
			expectErr:        true,
			expectErrMsg:     "ReadMe API Error: 500 on POST",
			setupSearch:      true,
			mockSearchStatus: 500,
		},
		{
			name:         "when no slug is provided",
			input:        "",
			expectErr:    true,
			expectErrMsg: "a doc slug or id must be provided",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the Search request if necessary
			if tt.setupSearch {
				mockStatus := 200
				if tt.mockSearchStatus != 0 {
					mockStatus = tt.mockSearchStatus
				}
				gock.New(TestClient.APIURL).
					Post(readme.DocEndpoint + "/search").
					Reply(mockStatus).
					JSON(tt.mockSearchResult)
			}

			// Mock the Get request
			if tt.mockResponse != nil {
				endpoint := readme.DocEndpoint + "/" + tt.slug
				if tt.options.ProductionDoc {
					endpoint += "/production"
				}
				gock.New(TestClient.APIURL).
					Get(endpoint).
					Reply(tt.mockStatus).
					JSON(tt.mockResponse)
			}
			defer gock.OffAll()

			// Act
			got, _, err := TestClient.Doc.Get(tt.input, tt.options)

			// Assert
			if tt.expectErr {
				assert.Error(t, err, "it returns an error")
				assert.Contains(t, err.Error(), tt.expectErrMsg, "it returns the expected error message")
			} else {
				assert.NoError(t, err, "it does not return an error")
				assert.Equal(t, tt.expect, got, "it returns expected Doc struct")
			}
			assert.True(t, gock.IsDone(), "it makes the expected API call")
		})
	}
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
		expectErr := "ReadMe API Error: 400 on DELETE"

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

		expectErr := "ReadMe API Error: 400 on POST"

		// Act
		_, _, err := TestClient.Doc.Search("something")

		// Assert
		assert.ErrorContains(t, err, expectErr, "it returns the expected error")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}
