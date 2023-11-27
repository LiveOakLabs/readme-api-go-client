package testdata

import "github.com/liveoaklabs/readme-api-go-client/readme"

// Version represents a version of a project.
// This is the response from the API when calling Get().
var Versions = []readme.Version{
	{
		Codename:     "",
		CreatedAt:    "2022-12-04T19:28:15.190Z",
		ID:           "638cf4cfdea3ff0096d1a95a",
		IsBeta:       false,
		IsDeprecated: false,
		IsHidden:     false,
		IsStable:     true,
		Version:      "1.0.0",
		VersionClean: "1.0.0",
	},
	{
		Codename:     "",
		CreatedAt:    "2023-12-04T13:04:32.056Z",
		ForkedFrom:   "638cf4cfdea3ff0096d1a95a",
		ID:           "4e8cf4cfdea3ff0096d1a95a",
		IsBeta:       false,
		IsDeprecated: false,
		IsHidden:     false,
		IsStable:     false,
		Version:      "1.1.0",
		VersionClean: "1.1.0",
	},
}

// VersionSummary represents a list of versions of a project.
// This is the response from the API when calling GetAll().
var VersionSummary = []readme.VersionSummary{
	{
		ForkedFrom:   Versions[0].ForkedFrom,
		Codename:     Versions[0].Codename,
		CreatedAt:    Versions[0].CreatedAt,
		ID:           Versions[0].ID,
		IsBeta:       Versions[0].IsBeta,
		IsDeprecated: Versions[0].IsDeprecated,
		IsHidden:     Versions[0].IsHidden,
		IsStable:     Versions[0].IsStable,
		Version:      Versions[0].Version,
		VersionClean: Versions[0].VersionClean,
	},
	{
		ForkedFrom:   Versions[1].ForkedFrom,
		Codename:     Versions[1].Codename,
		CreatedAt:    Versions[1].CreatedAt,
		ID:           Versions[1].ID,
		IsBeta:       Versions[1].IsBeta,
		IsDeprecated: Versions[1].IsDeprecated,
		IsHidden:     Versions[1].IsHidden,
		IsStable:     Versions[1].IsStable,
		Version:      Versions[1].Version,
		VersionClean: Versions[1].VersionClean,
	},
}
