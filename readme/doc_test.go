package readme_test

import (
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/internal/testutil"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

const (
	docsEndpoint = "http://readme-test.local/api/v1/docs"
)

var mockDocResponseBody string = `{
		"metadata": {
			"image": [
			  "https://files.readme.io/53c37f9-live-oak-logo-full-color.svg",
			  "live-oak-logo-full-color.svg",
			  950,
			  135,
			  "#1c0e52"
			],
			"title": "Example Doc",
			"description": "This is an example description."
		},
		"algolia": {
			"recordCount": 0,
			"publishPending": false,
			"updatedAt": "2023-01-18T05:22:34.529Z"
		},
		"api": {
			"method": "get",
			"url": "/{something}",
			"auth": "required",
			"params": [
			  {
				"name": "example_path_param",
				"type": "string",
				"enumValues": "",
				"default": "foo",
				"desc": "this is an example path param",
				"required": false,
				"in": "path",
				"ref": "",
				"_id": "63c77f40e790b60061a9a2ed"
			  }
		    ],
			"apiSetting": "63c77e0bf76dee008e0c5197",
			"examples": {
			  "codes": [
				{
				  "code": "",
				  "language": "c"
				}
			  ]
			},
			"results": {
				"codes": [
				  {
					"name": "",
					"code": "{\"message\": \"ok\"}",
					"language": "json",
					"status": 200
				  },
				  {
					"name": "",
					"code": "{}",
					"language": "json",
					"status": 400
				  }
				]
  		    }
		},
		"next": {
			"description": "",
			"pages": [
			  {
				"type": "doc",
				"icon": "file-text-o",
				"name": "The Next Doc",
				"slug": "the-next-doc",
				"deprecated": false,
				"category": "Documentation"
			  }
			]
		},
		"title": "Example Doc",
		"icon": "",
		"updates": [],
		"type": "basic",
		"slug": "example-doc",
		"excerpt": "",
		"body": "A turtle has been here.",
		"order": 999,
		"isReference": false,
		"deprecated": false,
		"hidden": true,
		"sync_unique": "",
		"link_url": "",
		"link_external": false,
		"pendingAlgoliaPublish": false,
		"previousSlug": "",
		"slugUpdatedAt": "2023-01-02T01:44:37.530Z",
		"revision": 2,
		"_id": "63b37e8b65fd5b0057af23f1",
		"category": "63a77777f52b9f006b6bf212",
		"project": "638cf4cedea3ff0096d1a955",
		"createdAt": "2023-01-03T01:02:03.731Z",
		"updatedAt": "2023-01-03T01:02:03.731Z",
		"user": "633c5a54187d2c008e2e074c",
		"version": "63a77777f52b9f006b6bf215",
		"__v": 0,
		"isApi": false,
		"id": "63b37e8b65fd5b0057af23f1",
		"body_html": "<div class=\"magic-block-textarea\"><p>A turtle has been here.</p></div>"
	}`

func Test_Docs_Get(t *testing.T) {
	expect := readme.Doc{}

	mockResponse := testutil.APITestResponse{
		URL:    docsEndpoint + "/example-doc",
		Status: 200,
		Body:   mockDocResponseBody,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	got, _, err := api.Doc.Get("example-doc")

	t.Run("returns no errors when request is valid and response is successful", func(t *testing.T) {
		assert.NoError(t, err)
	})
	t.Run("returns expected doc struct", func(t *testing.T) {
		assert.Equal(t, expect, got)
	})

	t.Run("returns a production doc when requested", func(t *testing.T) {
		mockResponse := testutil.APITestResponse{
			URL:    docsEndpoint + "/example-doc/production",
			Status: 200,
			Body:   mockDocResponseBody,
		}
		api := mockResponse.New(t)

		got, _, _ := api.Doc.Get("example-doc", readme.RequestOptions{ProductionDoc: true})

		assert.Equal(t, expect, got)
	})
}

func Test_Doc_Create(t *testing.T) {
	expect := readme.Doc{}

	hidden := false
	order := 999
	create := readme.DocParams{
		Body:          "&nbsp;",
		Title:         "Some Test Doc",
		Type:          "basic",
		Category:      "63a77777f52b9f006b6bf212",
		Hidden:        &hidden,
		Order:         &order,
		ParentDocSlug: "my-parent-doc",
		Error:         readme.DocErrorObject{Code: "404"},
	}

	mockResponse := testutil.APITestResponse{
		URL:    docsEndpoint,
		Status: 201,
		Body:   mockDocResponseBody,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	got, _, err := api.Doc.Create(create)

	assert.NoError(t, err)
	assert.Equal(t, expect, got)

	t.Run("handles invalid parameters", func(t *testing.T) {
		_, _, err := api.Doc.Create(readme.DocParams{})
		assert.Error(t, err, "returns error when title is empty")

		_, _, err = api.Doc.Create(readme.DocParams{Title: "testing"})
		assert.Error(t, err, "returns error when category is empty")
	})

	t.Run("uses provided RequestOptions", func(t *testing.T) {
		got, _, _ := api.Doc.Create(create, readme.RequestOptions{Version: "1.1.0"})
		assert.Equal(t, expect, got)
	})
}

func Test_Doc_Update(t *testing.T) {
	expect := readme.Doc{}

	mockResponse := testutil.APITestResponse{
		URL:    docsEndpoint + "/some-test-doc",
		Status: 200,
		Body:   mockDocResponseBody,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	hidden := false
	order := 999
	create := readme.DocParams{
		Body:     "&nbsp;",
		Title:    "Some Test Doc",
		Type:     "basic",
		Category: "63a77777f52b9f006b6bf212",
		Hidden:   &hidden,
		Order:    &order,
	}

	got, _, err := api.Doc.Update("some-test-doc", create)

	assert.NoError(t, err)
	assert.Equal(t, expect, got)

	t.Run("handles invalid parameters", func(t *testing.T) {
		_, _, err := api.Doc.Update("some-test-doc", readme.DocParams{})
		assert.Error(t, err, "returns error when title is empty")

		_, _, err = api.Doc.Update("some-test-doc", readme.DocParams{Title: "testing"})
		assert.Error(t, err, "returns error when category is empty")
	})

	t.Run("uses provided RequestOptions", func(t *testing.T) {
		got, _, _ := api.Doc.Update("some-test-doc", create, readme.RequestOptions{Version: "1.1.0"})
		assert.Equal(t, expect, got)
	})
}

func Test_Doc_Delete(t *testing.T) {
	mockResponse := testutil.APITestResponse{
		URL:    docsEndpoint + "/some-test-doc",
		Status: 204,
		Body:   "",
	}
	api := mockResponse.New(t)

	got, _, err := api.Doc.Delete("some-test-doc", readme.RequestOptions{Version: "1.1.0"})
	assert.NoError(t, err)
	assert.True(t, got)

	t.Run("returns error when API response has error", func(t *testing.T) {
		mockResponse := testutil.APITestResponse{
			URL:    docsEndpoint + "/some-test-doc",
			Status: 400,
			Body:   "",
		}
		api := mockResponse.New(t)

		got, _, err := api.Doc.Delete("some-test-doc", readme.RequestOptions{Version: "1.1.0"})
		assert.Error(t, err)
		assert.False(t, got)
	})
}

func Test_Doc_Search(t *testing.T) {
	// Arrange
	// result is what we expect from the search function call.
	// This is also used for the mock body response.
	result := []readme.DocSearchResult{{
		IndexName:    "Page",
		Title:        "Example of manual page",
		Type:         "basic",
		Slug:         "example-of-manual-page",
		IsReference:  false,
		Method:       "get",
		LinkURL:      "",
		Version:      "63b891d3ee384600680ce9f9",
		Project:      "63b891d3ee384600680ce9eb",
		ReferenceID:  "63c7828f56279c006fbd4178",
		Subdomain:    "jbeard-tf-dev-1",
		InternalLink: "docs/example-of-manual-page",
		ObjectID:     "63c7828f56279c006fbd4178-1",
		SnippetResult: readme.DocSearchResultSnippet{
			Title: readme.DocSearchResultSnippetValue{
				Value:      "Example of manual page",
				MatchLevel: "none",
			},
			Excerpt: readme.DocSearchResultSnippetValue{
				Value:      "",
				MatchLevel: "none",
			},
			Body: readme.DocSearchResultSnippetValue{
				Value:      "This is a test page created manually.",
				MatchLevel: "none",
			},
		},
		HighlightResult: readme.DocSearchResultHighlight{
			Title: readme.DocSearchResultHighlightValue{
				Value:      "Example of manual page",
				MatchLevel: "none",
				MatchedWords: []string{
					"manual",
				},
			},
			Excerpt: readme.DocSearchResultHighlightValue{
				Value:        "",
				MatchLevel:   "none",
				MatchedWords: []string{},
			},
			Body: readme.DocSearchResultHighlightValue{
				Value:      "This is a test page created manually.",
				MatchLevel: "none",
				MatchedWords: []string{
					"manual",
				},
			},
		},
		URL: "https://jbeard-tf-dev-1.readme.io/docs/example-of-manual-page",
	}}

	// response represents the JSON body response from the API, which lists results in a "results" key.
	response := readme.DocSearchResults{
		Results: result,
	}

	mockResponse := testutil.APITestResponse{
		URL:    docsEndpoint + "/search?search=manual",
		Status: 200,
		Body:   testutil.StructToJson(t, response),
	}
	api := mockResponse.New(t)

	t.Run("returns expected result when searching and API responds with 200", func(t *testing.T) {
		// Act
		got, _, err := api.Doc.Search("manual")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, result, got)
	})

	t.Run("returns expected result when version is specified", func(t *testing.T) {
		// Act
		_, _, err := api.Doc.Search("manual", readme.RequestOptions{Version: "1.1.0"})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns error when API response has error", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    docsEndpoint + "/search?search=something",
			Status: 400,
			Body:   `{"error":"INVALID", "message":"This is an error"}`,
		}
		api := mockResponse.New(t)

		// Act
		_, _, err := api.Doc.Search("something")

		// Assert
		assert.Error(t, err)
	})
}
