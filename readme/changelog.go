package readme

import (
	"encoding/json"
	"fmt"
)

// ChangelogEndpoint is ReadMe API Endpoint for changelogs.
const ChangelogEndpoint = "/changelogs"

// ChangelogService is an interface for using the changelog endpoints of the ReadMe.com API.
//
// API Reference: https://docs.readme.com/main/reference/getchangelogs
type ChangelogService interface {
	// Create a new changelog in ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/createchangelog
	Create(params ChangelogParams) (Changelog, *APIResponse, error)

	// Delete a changelog in ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/deletechangelog
	Delete(slug string) (bool, *APIResponse, error)

	// Get a changelog from ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/getchangelogs
	Get(slug string) (Changelog, *APIResponse, error)

	// GetAll retrieves a list of changelogs from ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/getchangelogs
	GetAll(options ...RequestOptions) ([]Changelog, *APIResponse, error)

	// Update an existing changelog in ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/updatechangelog
	Update(slug string, params ChangelogParams) (Changelog, *APIResponse, error)
}

// ChangelogClient handles communication with the docs related methods of the ReadMe.com API.
type ChangelogClient struct {
	client *Client
}

// Ensure the implementation satisfies the expected interfaces.
// This is a compile-time check.
// See: https://golang.org/doc/faq#guarantee_satisfies_interface
var _ ChangelogService = &ChangelogClient{}

// Changelog represents a Changelog object on ReadMe.
//
// This is the struct for retrieving and creating Changelog.
type Changelog struct {
	Algolia   DocAlgolia  `json:"algolia"`
	Body      string      `json:"body"`
	CreatedAt string      `json:"createdAt"`
	HTML      string      `json:"html"`
	Hidden    bool        `json:"hidden"`
	ID        string      `json:"_id"`
	Metadata  DocMetadata `json:"metadata"`
	Project   string      `json:"project,omitempty"`
	Revision  int         `json:"revision"`
	Slug      string      `json:"slug"`
	Title     string      `json:"title"`
	Type      string      `json:"type"`
	UpdatedAt string      `json:"updatedAt"`
	User      DocUser     `json:"user,omitempty"`
}

// ChangelogParams represents the parameters for creating and updating a Changelog on ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/createchangelog
type ChangelogParams struct {
	// Body content of the changelog.
	// This is *required* when creating or updating a changelog.
	Body string `json:"body"`
	// Hidden toggles the visibility of the changelog.
	// API default is `true`.
	Hidden *bool `json:"hidden"`
	// Title of the changelog.
	// This is *required* when creating or updating a changelog.
	Title string `json:"title"`
	// Type of the changelog.
	Type string `json:"type,omitempty"`
}

// validChangelogType validates the 'type' field when creating or updating a changelog.
func validChangelogType(changelogType string) bool {
	switch changelogType {
	case "added", "fixed", "improved", "deprecated", "removed":
		return true
	}

	return false
}

// GetAll retrieves a list of changelogs from ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/getchangelogs
func (c ChangelogClient) GetAll(options ...RequestOptions) ([]Changelog, *APIResponse, error) {
	var response []Changelog

	apiRequest := &APIRequest{
		Method:       "GET",
		Endpoint:     ChangelogEndpoint,
		UseAuth:      true,
		OkStatusCode: []int{200},
		Response:     &response,
	}

	if len(options) > 0 {
		apiRequest.RequestOptions = options[0]
	}

	apiResponse, err := c.client.APIRequest(apiRequest)

	return response, apiResponse, err
}

// Get retrieves a single changelog from ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/getchangelog
func (c ChangelogClient) Get(slug string) (Changelog, *APIResponse, error) {
	response := Changelog{}
	apiRequest := &APIRequest{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("%s/%s", ChangelogEndpoint, slug),
		UseAuth:      true,
		OkStatusCode: []int{200},
		Response:     &response,
	}

	apiResponse, err := c.client.APIRequest(apiRequest)

	return response, apiResponse, err
}

// Create a new changelog in ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/createchangelog
func (c ChangelogClient) Create(params ChangelogParams) (Changelog, *APIResponse, error) {
	if params.Title == "" {
		return Changelog{}, nil, fmt.Errorf("title must be provided")
	}
	if !validChangelogType(params.Type) {
		return Changelog{}, nil, fmt.Errorf("type must be added, fixed, improved, deprecated, or removed")
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return Changelog{}, nil, fmt.Errorf("unable to parse request: %w", err)
	}

	response := Changelog{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "POST",
		Endpoint:     ChangelogEndpoint,
		UseAuth:      true,
		Payload:      payload,
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		OkStatusCode: []int{201},
		Response:     &response,
	})

	return response, apiResponse, err
}

// Update an existing changelog in ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/updatechangelog
func (c ChangelogClient) Update(slug string, params ChangelogParams) (Changelog, *APIResponse, error) {
	if params.Title == "" {
		return Changelog{}, nil, fmt.Errorf("title must be provided")
	}
	if !validChangelogType(params.Type) {
		return Changelog{}, nil, fmt.Errorf("type must be added, fixed, improved, deprecated, or removed")
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return Changelog{}, nil, fmt.Errorf("unable to parse request: %w", err)
	}

	response := Changelog{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "PUT",
		Endpoint:     fmt.Sprintf("%s/%s", ChangelogEndpoint, slug),
		UseAuth:      true,
		Payload:      payload,
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		OkStatusCode: []int{200},
		Response:     &response,
	})

	return response, apiResponse, err
}

// Delete a changelog in ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/deletechangelog
func (c ChangelogClient) Delete(slug string) (bool, *APIResponse, error) {
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "DELETE",
		Endpoint:     fmt.Sprintf("%s/%s", ChangelogEndpoint, slug),
		UseAuth:      true,
		OkStatusCode: []int{204},
	})
	if err != nil {
		return false, apiResponse, err
	}

	return true, apiResponse, nil
}
