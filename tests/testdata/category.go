package testdata

import (
	"net/http"

	"github.com/liveoaklabs/readme-api-go-client/readme"
)

var Categories = []readme.Category{
	{
		CreatedAt: "2022-12-04T19:28:15.240Z",
		ID:        "63a77777f52b9f006b6bf212",
		Order:     10,
		Project:   "go-testing",
		Reference: false,
		Slug:      "documentation",
		Title:     "Documentation",
		Type:      "guide",
		Version:   Versions[0].Version,
	},
	{
		CreatedAt: "2022-12-04T19:28:15.240Z",
		ID:        "6543c5bf91e232000cbdc4cf",
		Order:     20,
		Project:   "go-testing",
		Slug:      "example-category",
		Title:     "Example Category",
		Type:      "guide",
		Version:   Versions[0].Version,
	},
	{
		CreatedAt: "2022-12-04T19:28:15.240Z",
		ID:        "654ac951d35521000d16d325",
		Order:     30,
		Project:   "go-testing",
		Slug:      "swagger-petstore-openapi",
		Title:     "Swagger Petstore - OpenAPI",
		Type:      "reference",
		Version:   Versions[0].Version,
	},
}

var CategoryResponse = &readme.APIResponse{
	Body: []byte(ToJSON(Categories)),
	HTTPResponse: &http.Response{
		StatusCode: 200,
	},
}

var CategoryDocs = []readme.CategoryDocs{
	{
		Hidden: false,
		ID:     "63a77777f52b9f006b6bf212",
		Order:  10,
		Slug:   "documentation",
		Title:  "Documentation",
		Children: []readme.CategoryDocs{
			{
				Hidden: false,
				ID:     "63a77777f52b9f006b6bf213",
				Order:  10,
				Slug:   "child-doc-1",
				Title:  "Child Doc 1",
				Children: []readme.CategoryDocs{
					{
						Hidden: false,
						ID:     "63a77777f52b9f006b6bf214",
						Order:  10,
						Slug:   "child-doc-1-1",
						Title:  "child doc 1.1",
					},
				},
			},
		},
	},
}

var CategorySaved = readme.CategorySaved{
	Title:     Categories[0].Title,
	Slug:      Categories[0].Slug,
	Order:     Categories[0].Order,
	ID:        Categories[0].ID,
	CreatedAt: Categories[0].CreatedAt,
	Reference: Categories[0].Reference,
	Type:      Categories[0].Type,
	Version:   Versions[0],
	Project:   Categories[0].Project,
}

var CategoryVersion = readme.CategoryVersion{
	Categories: []readme.Category{Categories[0]},
	Codename:   Versions[0].Codename,
	CreatedAt:  Versions[0].CreatedAt,
	ForkedFrom: readme.CategoryVersionForkedFrom{
		Version: Versions[0].VersionClean,
		ID:      Versions[0].Version,
	},
	ID:           Versions[0].ID,
	IsBeta:       Versions[0].IsBeta,
	IsDeprecated: Versions[0].IsDeprecated,
	IsHidden:     Versions[0].IsHidden,
	IsStable:     Versions[0].IsStable,
	Project:      Versions[0].Project,
	ReleaseDate:  Versions[0].CreatedAt,
	Version:      Versions[0].Version,
	VersionClean: Versions[0].VersionClean,
}

var CategoryVersionSaved = readme.CategoryVersionSaved{
	Title:     Categories[0].Title,
	Slug:      Categories[0].Slug,
	Order:     Categories[0].Order,
	ID:        Categories[0].ID,
	CreatedAt: Categories[0].CreatedAt,
	Reference: Categories[0].Reference,
	Type:      Categories[0].Type,
	Version:   CategoryVersion,
	Project:   Categories[0].Project,
}
