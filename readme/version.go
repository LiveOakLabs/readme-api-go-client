package readme

import (
	"encoding/json"
	"fmt"
)

// VersionEndpoint is the ReadMe API URL endpoint for Version metadata.
const VersionEndpoint = "/version"

// VersionService is an interface for using the version endpoints of the ReadMe.com API.
//
// See: https://docs.readme.com/main/reference/getproject
type VersionService interface {
	// Create a new version within a project.
	//
	// API Reference: https://docs.readme.com/main/reference/createversion
	Create(prams VersionParams) (Version, *APIResponse, error)

	// Delete a version.
	//
	// The version may be provided using either the semver identifier for the project version
	// ('1.0.0') or an API ID for the version prefixed with 'id' ('id:63ac899d11c4680047ec5970').
	//
	// When using the semantic version identifier, use the formatted VersionClean value listed in
	// the response from the GetAll() function for best results.
	//
	// API Reference: https://docs.readme.com/main/reference/deleteversion
	Delete(version string) (bool, *APIResponse, error)

	// Get a single version.
	//
	// The version may be provided using either the semver identifier for the project version
	// ('1.0.0') or an API ID for the version prefixed with 'id' ('id:63ac899d11c4680047ec5970').
	//
	// When using the semantic version identifier, use the formatted VersionClean value listed in
	// the response from the GetAll() function for best results.
	//
	// API Reference: https://docs.readme.com/main/reference/getversion
	Get(version string) (Version, *APIResponse, error)

	// GetAll retrieves a list of versions associated with an API key.
	//
	// API Reference: https://docs.readme.com/main/reference/getversions
	GetAll() ([]VersionSummary, *APIResponse, error)

	// Update an existing version.
	//
	// The version may be provided using either the semver identifier for the project version
	// ('1.0.0') or an API ID for the version prefixed with 'id' ('id:63ac899d11c4680047ec5970').
	//
	// When using the semantic version identifier, use the formatted VersionClean value listed in
	// the response from the GetAll() function for best results.
	//
	// API Reference: https://docs.readme.com/main/reference/updateversion
	Update(version string, params VersionParams) (Version, *APIResponse, error)

	// GetVersion parses a provided string to determine if it it's a semantic version identifier (1.0.0)
	// or an API version identifier (id:63ac899d11c4680047ec5970). If it's an API version identifier,
	// the value is compared with the results from GetAll() to return the semantic version that's used
	// for API requests. If the specified version is already a semantic version string, it will be
	// returned as-is.
	GetVersion(version string) (string, error)
}

// VersionClient handles communication with the Project related methods of the ReadMe.com API.
type VersionClient struct {
	client *Client
}

// Version represents the details of a specific version.
type Version struct {
	Categories   []string `json:"categories"`
	Codename     string   `json:"codename"`
	CreatedAt    string   `json:"createdAt"`
	ForkedFrom   string   `json:"forked_from"`
	ID           string   `json:"_id"`
	IsBeta       bool     `json:"is_beta"`
	IsDeprecated bool     `json:"is_deprecated"`
	IsHidden     bool     `json:"is_hidden"`
	IsStable     bool     `json:"is_stable"`
	Project      string   `json:"project"`
	ReleaseDate  string   `json:"releaseDate"`
	Version      string   `json:"version"`
	VersionClean string   `json:"version_clean"`
}

// VersionSummary represents the response from the ReadMe API when retrieving a list of versions
// with GetAll().
type VersionSummary struct {
	Codename     string `json:"codename"`
	CreatedAt    string `json:"createdAt"`
	ForkedFrom   string `json:"forked_from"`
	ID           string `json:"_id"`
	IsBeta       bool   `json:"is_beta"`
	IsDeprecated bool   `json:"is_deprecated"`
	IsHidden     bool   `json:"is_hidden"`
	IsStable     bool   `json:"is_stable"`
	Version      string `json:"version"`
	VersionClean string `json:"version_clean"`
}

// VersionParams represents the request parameters used when creating or updating a version.
type VersionParams struct {
	// Codename is the dubbed name of version.
	Codename string `json:"codename,omitempty"`
	// From is the semantic Version to use as the base fork.
	// This is *required* when creating or updating a version.
	From string `json:"from"`
	// IsBeta toggles whether the version is a beta release.
	// API default is `false`.
	IsBeta *bool `json:"is_beta"`
	// IsDeprecated toggles whether the version is deprecated. Only allowed when updating.
	// API default is `false`.
	IsDeprecated *bool `json:"is_deprecated"`
	// IsHidden toggles public accessibility.
	// API default is `false`.
	IsHidden *bool `json:"is_hidden"`
	// IsStable toggles whether the version should be the project's main version.
	// API default is `false`.
	IsStable *bool `json:"is_stable"`
	// Version is the number.
	// This is *required* when creating or updating a version.
	Version string `json:"version"`
}

// Ensure the implementation satisfies the expected interfaces.
// This is a compile-time check.
// See: https://golang.org/doc/faq#guarantee_satisfies_interface
var _ VersionService = &VersionClient{}

// GetVersion parses a provided string to determine if it it's a semantic version identifier (1.0.0)
// or an API version identifier (id:63ac899d11c4680047ec5970). If it's an API version identifier,
// the value is compared with the results from GetAll() to return the semantic version that's used
// for API requests. If the specified version is already a semantic version string, it will be
// returned as-is.
func (c VersionClient) GetVersion(version string) (string, error) {
	isID, reqID := ParseID(version)
	if !isID {
		return version, nil
	}

	all, _, err := c.client.Version.GetAll()
	if err != nil {
		return "", fmt.Errorf("unable to get list of versions: %w", err)
	}

	for _, vers := range all {
		if vers.ID == reqID {
			return vers.Version, nil
		}
	}

	return "", fmt.Errorf("no match for version ID %s", reqID)
}

// GetAll retrieves a list of versions associated with an API key.
//
// API Reference: https://docs.readme.com/main/reference/getversions
func (c VersionClient) GetAll() ([]VersionSummary, *APIResponse, error) {
	var versions []VersionSummary

	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "GET",
		Endpoint:     VersionEndpoint,
		UseAuth:      true,
		OkStatusCode: []int{200},
		Response:     &versions,
	})

	return versions, apiResponse, err
}

// Get a single version.
//
// The version may be provided using either the semver identifier for the project version ('1.0.0')
// or an API ID for the version prefixed with 'id' ('id:63ac899d11c4680047ec5970').
//
// When using the semantic version identifier, use the formatted VersionClean value listed in the
// response from the GetAll() function for best results.
//
// API Reference: https://docs.readme.com/main/reference/getversion
func (c VersionClient) Get(version string) (Version, *APIResponse, error) {
	version, err := c.GetVersion(version)
	if err != nil {
		return Version{}, nil, err
	}

	versionResponse := Version{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("%s/%s", VersionEndpoint, version),
		UseAuth:      true,
		OkStatusCode: []int{200},
		Response:     &versionResponse,
	})

	return versionResponse, apiResponse, err
}

// Create a new version within a project.
//
// API Reference: https://docs.readme.com/main/reference/createversion
func (c VersionClient) Create(params VersionParams) (Version, *APIResponse, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return Version{}, nil, fmt.Errorf("unable to parse request: %w", err)
	}

	response := Version{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "POST",
		Endpoint:     VersionEndpoint,
		UseAuth:      true,
		Payload:      payload,
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		OkStatusCode: []int{200},
		Response:     &response,
	})

	return response, apiResponse, err
}

// Update an existing version.
//
// The version may be provided using either the semver identifier for the project version ('1.0.0')
// or an API ID for the version prefixed with 'id' ('id:63ac899d11c4680047ec5970').
//
// When using the semantic version identifier, use the formatted VersionClean value listed in the
// response from the GetAll() function for best results.
//
// API Reference: https://docs.readme.com/main/reference/updateversion
func (c VersionClient) Update(version string, params VersionParams) (Version, *APIResponse, error) {
	version, err := c.GetVersion(version)
	if err != nil {
		return Version{}, nil, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return Version{}, nil, fmt.Errorf("unable to parse request: %w", err)
	}

	response := Version{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "PUT",
		Endpoint:     fmt.Sprintf("%s/%s", VersionEndpoint, version),
		UseAuth:      true,
		Payload:      payload,
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		OkStatusCode: []int{200},
		Response:     &response,
	})

	return response, apiResponse, err
}

// Delete a version.
//
// The version may be provided using either the semver identifier for the project version ('1.0.0')
// or an API ID for the version prefixed with 'id' ('id:63ac899d11c4680047ec5970').
//
// When using the semantic version identifier, use the formatted VersionClean value listed in the
// response from the GetAll() function for best results.
//
// API Reference: https://docs.readme.com/main/reference/deleteversion
func (c VersionClient) Delete(version string) (bool, *APIResponse, error) {
	version, err := c.GetVersion(version)
	if err != nil {
		return false, nil, err
	}

	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "DELETE",
		Endpoint:     fmt.Sprintf("%s/%s", VersionEndpoint, version),
		UseAuth:      true,
		OkStatusCode: []int{200},
	})
	if err != nil {
		return false, apiResponse, err
	}

	return true, apiResponse, nil
}
