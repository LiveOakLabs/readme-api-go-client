package testdata

import "github.com/liveoaklabs/readme-api-go-client/readme"

var CustomPages = []readme.CustomPage{
	{
		Algolia: readme.DocAlgolia{
			PublishPending: false,
			RecordCount:    0,
			UpdatedAt:      "2023-01-03T00:29:22.169Z",
		},
		Body:       "This is a test changelog",
		CreatedAt:  "2023-01-03T00:29:22.169Z",
		Fullscreen: false,
		HTML:       "<div class=\"magic-block-textarea\"><p>This is a test changelog</p>\n\n</div>",
		HTMLMode:   false,
		Hidden:     true,
		ID:         "63b376e244ed08009d672b11",
		Metadata: readme.DocMetadata{
			Image:       []any{},
			Title:       "",
			Description: "",
		},
		Revision:  2,
		Slug:      "some-test",
		Title:     "Some Test",
		UpdatedAt: "2023-01-03T00:29:22.169Z",
	},
}
