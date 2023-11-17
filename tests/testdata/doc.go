package testdata

import "github.com/liveoaklabs/readme-api-go-client/readme"

var Docs = []readme.Doc{
	{
		Metadata: readme.DocMetadata{
			Image:       []any{},
			Title:       "Example Doc",
			Description: "This is an example description.",
		},
		Algolia: readme.DocAlgolia{
			PublishPending: false,
			RecordCount:    0,
			UpdatedAt:      "2023-01-18T05:22:34.529Z",
		},
		API: readme.DocAPI{
			Method:     "get",
			URL:        "/{something}",
			Auth:       "required",
			APISetting: "63c77e0bf76dee008e0c5197",
			Examples: readme.DocAPIExamples{
				Codes: []readme.DocAPIExamplesCodes{
					{
						Code:     "",
						Language: "c",
					},
				},
			},
			Params: []readme.DocAPIParams{
				{
					Name:       "example_path_param",
					Type:       "string",
					EnumValues: "",
					Default:    "foo",
					Desc:       "this is an example path param",
					Required:   false,
					In:         "path",
					Ref:        "",
					ID:         "63c77f40e790b60061a9a2ed",
				},
			},
			Results: readme.DocAPIResults{
				Codes: []readme.DocAPIResultsCodes{
					{
						Name:     "",
						Code:     "{\"message\": \"ok\"}",
						Language: "json",
						Status:   200,
					},
					{
						Name:     "",
						Code:     "{}",
						Language: "json",
						Status:   400,
					},
				},
			},
		},
		Next: readme.DocNext{
			Description: "",
			Pages: []readme.DocNextPages{
				{
					Type:       "doc",
					Icon:       "file-text-o",
					Name:       "The Next Doc",
					Slug:       "the-next-doc",
					Deprecated: false,
					Category:   "Documentation",
				},
			},
		},
		Title:         "Example Doc",
		Icon:          "",
		Updates:       []any{},
		Type:          "basic",
		Slug:          "example-doc",
		Excerpt:       "",
		Body:          "A turtle has been here.",
		Order:         999,
		IsReference:   false,
		Deprecated:    false,
		Hidden:        false,
		SyncUnique:    "",
		LinkURL:       "",
		LinkExternal:  false,
		PreviousSlug:  "",
		SlugUpdatedAt: "2023-01-02T01:44:37.530Z",
		Revision:      2,
		ID:            "63b37e8b65fd5b0057af23f1",
		Category:      Categories[0].ID,
		Project:       Project.Name,
		CreatedAt:     "2023-01-03T01:02:03.731Z",
		UpdatedAt:     "2023-01-03T01:02:03.731Z",
		User:          "633c5a54187d2c008e2e074c",
		Version:       Versions[0].Version,
		IsAPI:         false,
		BodyHTML:      "<div class=\"magic-block-textarea\"><p>A turtle has been here.</p></div>",
	},
	{
		Metadata: readme.DocMetadata{
			Image:       []any{},
			Title:       "Example Doc 2",
			Description: "This is an example description.",
		},
		Algolia: readme.DocAlgolia{
			PublishPending: false,
			RecordCount:    0,
			UpdatedAt:      "2023-01-18T05:22:34.529Z",
		},
		API: readme.DocAPI{},
		Next: readme.DocNext{
			Description: "",
			Pages: []readme.DocNextPages{
				{
					Type:       "doc",
					Icon:       "file-text-o",
					Name:       "The Next Doc",
					Slug:       "the-next-doc",
					Deprecated: false,
					Category:   "Documentation",
				},
			},
		},
		Title:         "Example Doc 2",
		Icon:          "",
		Updates:       []any{},
		Type:          "basic",
		Slug:          "example-doc-2",
		Excerpt:       "",
		Body:          "A turtle has been here.",
		Order:         999,
		IsReference:   false,
		Deprecated:    false,
		Hidden:        true,
		SyncUnique:    "",
		LinkURL:       "",
		LinkExternal:  false,
		PreviousSlug:  "",
		SlugUpdatedAt: "2023-01-02T01:44:37.530Z",
		Revision:      2,
		ID:            "63b37e8b65fd5b0057af23f2",
		Category:      Categories[0].ID,
		Project:       Project.Name,
		CreatedAt:     "2023-01-03T01:02:03.731Z",
		UpdatedAt:     "2023-01-03T01:02:03.731Z",
		User:          "633c5a54187d2c008e2e074c",
		Version:       Versions[0].Version,
		IsAPI:         false,
		BodyHTML:      "<div class=\"magic-block-textarea\"><p>A turtle has been here.</p></div>",
	},
}

var DocSearchResult = []readme.DocSearchResult{
	{
		ObjectID: Docs[0].ID,
		HighlightResult: readme.DocSearchResultHighlight{
			Body: readme.DocSearchResultHighlightValue{
				Value:      Docs[0].Body,
				MatchLevel: "none",
				MatchedWords: []string{
					"turtle",
				},
			},
			Title: readme.DocSearchResultHighlightValue{
				Value:      Docs[0].Title,
				MatchLevel: "none",
				MatchedWords: []string{
					"example",
				},
			},
		},
		Slug:        Docs[0].Slug,
		Title:       Docs[0].Title,
		Version:     Docs[0].Version,
		ReferenceID: Docs[0].ID,
	},
}

var DocSearchResults = readme.DocSearchResults{
	Results: DocSearchResult,
}
