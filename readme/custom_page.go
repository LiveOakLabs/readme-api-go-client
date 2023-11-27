package readme

import (
	"encoding/json"
	"fmt"
)

// CustomPageEndpoint is the ReadMe API Endpoint for custom pages.
const CustomPageEndpoint = "/custompages"

// CustomPageService is an interface for using the custom page endpoints of the ReadMe API.
type CustomPageService interface {
	// Create a new custom page in ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/createcustompage
	Create(params CustomPageParams) (CustomPage, *APIResponse, error)

	// Delete a custom page in ReadMe.
	//
	// API Reference: https://docs.readme.com/reference/deletecustompages
	Delete(slug string) (bool, *APIResponse, error)

	// Get a single custom page's data from ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/getcustompage
	Get(slug string) (CustomPage, *APIResponse, error)

	// GetAll retrieves a list of custom pages and their data from ReadMe.
	//
	// Pagination options may be specified with the `options` parameter.
	//
	// API Reference: https://docs.readme.com/main/reference/getcustompages
	GetAll(options ...RequestOptions) ([]CustomPage, *APIResponse, error)

	// Update an existing custom page in ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/updatecustompage
	Update(slug string, params CustomPageParams) (CustomPage, *APIResponse, error)
}

// CustomPageClient handles communication with the custom page related methods of the ReadMe API.
type CustomPageClient struct {
	client *Client
}

// Ensure the implementation satisfies the expected interfaces.
var _ CustomPageService = &CustomPageClient{}

// CustomPage represents a custom page in ReadMe.
type CustomPage struct {
	Algolia    DocAlgolia  `json:"algolia"`
	Body       string      `json:"body"`
	CreatedAt  string      `json:"createdAt"`
	Fullscreen bool        `json:"fullscreen"`
	HTML       string      `json:"html"`
	Hidden     bool        `json:"hidden"`
	HTMLMode   bool        `json:"htmlmode"`
	ID         string      `json:"_id"`
	Metadata   DocMetadata `json:"metadata"`
	Revision   int         `json:"revision"`
	Slug       string      `json:"slug"`
	Title      string      `json:"title"`
	UpdatedAt  string      `json:"updatedAt"`
}

// CustomPageParams represents the parameters to create or update a custom page in ReadMe.
type CustomPageParams struct {
	// Body formatted in Markdown (displayed by default).
	Body string `json:"body,omitempty"`
	// Hidden toggles the visibility of the custom page.
	// API default is `true`.
	Hidden *bool `json:"hidden"`
	// Body formatted in HTML (sanitized, only displayed if HTMLMode is true).
	HTML string `json:"html,omitempty"`
	// HTMLMode toggles if html should be displayed. Body will be displayed if false.
	// API default is `false`.
	HTMLMode *bool `json:"htmlmode"`
	// Title of the custom page.
	// This is *required* when creating or updating a custom page.
	Title string `json:"title"`
}

// GetAll retrieves a list of custom pages and their data from ReadMe.
//
// Pagination options may be specified with the `options` parameter.
//
// API Reference: https://docs.readme.com/main/reference/getcustompages
func (c CustomPageClient) GetAll(options ...RequestOptions) ([]CustomPage, *APIResponse, error) {
	var customPages []CustomPage
	var apiResponse *APIResponse
	var err error
	hasNextPage := false

	// Initialize pagination counter.
	page := 1
	if len(options) > 0 {
		if options[0].Page != 0 {
			page = options[0].Page
		}
	}

	for {
		var paginatedResult []CustomPage

		apiRequest := &APIRequest{
			Method:       "GET",
			Endpoint:     CustomPageEndpoint,
			UseAuth:      true,
			OkStatusCode: []int{200},
			Response:     &paginatedResult,
		}
		if len(options) > 0 {
			apiRequest.RequestOptions = options[0]
		}

		apiResponse, hasNextPage, err = c.client.paginatedRequest(apiRequest, page)
		if err != nil {
			return customPages, apiResponse, fmt.Errorf("unable to retrieve custom pages: %w", err)
		}
		customPages = append(customPages, paginatedResult...)

		if !hasNextPage {
			break
		}

		page = page + 1
	}

	return customPages, apiResponse, nil
}

// Get a single custom page's data from ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/getcustompage
func (c CustomPageClient) Get(slug string) (CustomPage, *APIResponse, error) {
	customPage := CustomPage{}

	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("%s/%s", CustomPageEndpoint, slug),
		UseAuth:      true,
		OkStatusCode: []int{200},
		Response:     &customPage,
	})

	return customPage, apiResponse, err
}

// Create a new custom page in ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/createcustompage
func (c CustomPageClient) Create(params CustomPageParams) (CustomPage, *APIResponse, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return CustomPage{}, nil, fmt.Errorf("unable to marshal request: %w", err)
	}

	response := CustomPage{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "POST",
		Endpoint:     CustomPageEndpoint,
		UseAuth:      true,
		Payload:      payload,
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		OkStatusCode: []int{201},
		Response:     &response,
	})

	return response, apiResponse, err
}

// Update an existing custom page in ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/updatecustompage
func (c CustomPageClient) Update(slug string, params CustomPageParams) (CustomPage, *APIResponse, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return CustomPage{}, nil, fmt.Errorf("unable to marshal request: %w", err)
	}

	response := CustomPage{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "PUT",
		Endpoint:     fmt.Sprintf("%s/%s", CustomPageEndpoint, slug),
		UseAuth:      true,
		Payload:      payload,
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		OkStatusCode: []int{200},
		Response:     &response,
	})

	return response, apiResponse, err
}

// Delete a custom page in ReadMe.
//
// API Reference: https://docs.readme.com/reference/deletecustompages
func (c CustomPageClient) Delete(slug string) (bool, *APIResponse, error) {
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "DELETE",
		Endpoint:     fmt.Sprintf("%s/%s", CustomPageEndpoint, slug),
		UseAuth:      true,
		OkStatusCode: []int{204},
	})
	if err != nil {
		return false, apiResponse, err
	}

	return true, apiResponse, nil
}
