package testdata

import (
	"github.com/liveoaklabs/readme-api-go-client/readme"
)

var Changelogs = []readme.Changelog{
	{
		Metadata: readme.DocMetadata{
			Image:       []any{},
			Title:       "",
			Description: "",
		},
		Algolia: readme.DocAlgolia{
			PublishPending: false,
			RecordCount:    0,
			UpdatedAt:      "2023-01-03T00:29:22.169Z",
		},
		Title:     "Some Test",
		Slug:      "some-test",
		Body:      "This is a test changelog",
		Type:      "added",
		Hidden:    true,
		Revision:  2,
		ID:        "63b376e244ed08009d672b11",
		CreatedAt: "2023-01-03T00:29:22.169Z",
		UpdatedAt: "2023-01-03T00:29:22.169Z",
		HTML:      "<div class=\"magic-block-textarea\"><p>This is a test changelog</p>\n\n</div>",
	},
	{
		Metadata: readme.DocMetadata{
			Image:       []any{},
			Title:       "",
			Description: "",
		},
		Algolia: readme.DocAlgolia{
			PublishPending: false,
			RecordCount:    0,
			UpdatedAt:      "2023-02-10T17:03:56.234Z",
		},
		Title:     "Second Test",
		Slug:      "second-test",
		Body:      "This is a second test changelog",
		Type:      "added",
		Hidden:    true,
		Revision:  8,
		ID:        "4ce276e244ef16009d671ae3",
		CreatedAt: "2023-01-03T00:29:22.169Z",
		UpdatedAt: "2023-02-10T17:03:56.234Z",
		HTML:      "<div class=\"magic-block-textarea\"><p>This is a second test changelog</p>\n\n</div>",
	},
}
