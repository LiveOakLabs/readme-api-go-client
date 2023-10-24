package testdata

import (
	"net/http"

	"github.com/liveoaklabs/readme-api-go-client/readme"
)

var APISpecifications = []readme.APISpecification{
	{
		ID:         "0123456789",
		Title:      "Readme Testing",
		LastSynced: "2022-11-29",
		Source:     "0a1b2c3d4e5f",
		Type:       "guide",
		Version:    "abcdef0123456789",
		Category: readme.CategorySummary{
			Title: "TestCatTitle",
			Slug:  "testCatSlug",
			Order: 0,
			Type:  "doc",
			ID:    "00aa11bb22cc33dd44ee55ff",
		},
	},
}

var APISpecResponseVersionEmtpy = readme.APIResponse{
	APIErrorResponse: readme.APIErrorResponse{
		Error:      "VERSION_EMPTY",
		Message:    "string",
		Suggestion: "string",
		Docs:       "https://docs.readme.com/logs/6883d0ee-cf79-447a-826f-a48f7d5bdf5f",
		Help:       "If you need help, email support@readme.io",
		Poem: []string{
			"If you're seeing this error,",
			"Things didn't quite go the way we hoped.",
			"When we tried to process your request,",
			"Maybe trying again it'll work—who knows!",
		},
	},
	HTTPResponse: &http.Response{
		StatusCode: 400,
	},
}

var APISpecResponseSpecFileEmpty = readme.APIResponse{
	APIErrorResponse: readme.APIErrorResponse{
		Error:      "SPEC_FILE_EMPTY",
		Message:    "string",
		Suggestion: "string",
		Docs:       "https://docs.readme.com/logs/6883d0ee-cf79-447a-826f-a48f7d5bdf5f",
		Help:       "If you need help, email support@readme.io",
		Poem: []string{
			"If you're seeing this error,",
			"Things didn't quite go the way we hoped.",
			"When we tried to process your request,",
			"Maybe trying again it'll work—who knows!",
		},
	},
	HTTPResponse: &http.Response{
		StatusCode: 400,
	},
}

var APISpecResponseSpecFileInvalid = readme.APIResponse{
	APIErrorResponse: readme.APIErrorResponse{
		Error:      "ERROR_SPEC_INVALID",
		Message:    "string",
		Suggestion: "string",
		Docs:       "https://docs.readme.com/logs/6883d0ee-cf79-447a-826f-a48f7d5bdf5f",
		Help:       "If you need help, email support@readme.io",
		Poem: []string{
			"If you're seeing this error,",
			"Things didn't quite go the way we hoped.",
			"When we tried to process your request,",
			"Maybe trying again it'll work—who knows!",
		},
	},
	HTTPResponse: &http.Response{
		StatusCode: 400,
	},
}
