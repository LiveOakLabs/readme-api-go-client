package readme

import (
	"encoding/json"
	"errors"
	"fmt"
)

// CategoryEndpoint is the ReadMe API endpoint for categories.
const CategoryEndpoint = "/categories"

// CategoryService is an interface for using the category endpoints of the ReadMe.com API.
//
// API Reference: https://docs.readme.com/main/reference/getcategories
type CategoryService interface {
	// Create a new category in ReadMe.
	//
	// The `response` parameter should be a to a `CategorySaved` or `CategoryVersionSaved`
	// interface. When specifying a `Version` in the `options` parameter, the API responds with
	// `CategoryVersionSaved`. When a version isn't specified, it responds with `CategorySaved`.
	//
	// Without a version:
	//
	//	category := &readme.CategorySaved{}
	//	apiResponse, err := rdme.Category.Create(category, {params...})
	//
	// With a version:
	//
	//	options := readme.RequestOptions{Version: "1.2.0"}
	//	category := &readme.CategoryVersionSaved{}
	//	apiResponse, err := rdme.Category.Create(category, {params...}, options)
	//
	// API Reference: https://docs.readme.com/main/reference/createcategory
	Create(response any, params CategoryParams, options ...RequestOptions) (*APIResponse, error)

	// Delete an existing category in ReadMe.
	//
	// API Reference: https://docs.readme.com/reference/deletecategory
	Delete(slug string, options ...RequestOptions) (bool, *APIResponse, error)

	// Get a single category on ReadMe.com.
	//
	// The `category` parameter may be a slug or category ID prefixed with "id:".
	// The slug is preferred, since it's a more direct request while the ID requires
	// iterating over the list of all categories for a matching ID.
	//
	// API Reference: https://docs.readme.com/reference/getcategory
	Get(category string, options ...RequestOptions) (Category, *APIResponse, error)

	// GetAll retrieves and returns all categories on ReadMe.com.
	//
	// API Reference: https://docs.readme.com/reference/getcategories
	GetAll(options ...RequestOptions) ([]Category, *APIResponse, error)

	// GetDocs a list of docs metadata for a category on ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/getcategorydocs
	GetDocs(slug string, options ...RequestOptions) ([]CategoryDocs, *APIResponse, error)

	// Update an existing category in ReadMe.
	//
	// Note that Update() returns a Category struct type while Create() returns
	// CategorySaved or CategoryVersionSaved.
	//
	// API Reference: https://docs.readme.com/main/reference/updatecategory
	Update(slug string, params CategoryParams, options ...RequestOptions) (Category, *APIResponse, error)
}

// CategoryClient handles communication with the categories related methods of the ReadMe.com API.
type CategoryClient struct {
	client *Client
}

// Ensure the implementation satisfies the expected interfaces.
// This is a compile-time check.
// See: https://golang.org/doc/faq#guarantee_satisfies_interface
var _ CategoryService = &CategoryClient{}

// CategorySummary represents basic information about a category. This is used in the response for
// API specification metadata.
type CategorySummary struct {
	ID    string `json:"id"`
	Order int    `json:"order"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

// Category represents a single category in ReadMe.
type Category struct {
	CategoryType string `json:"categoryType"`
	CreatedAt    string `json:"createdAt"`
	ID           string `json:"id"`
	Order        int    `json:"order"`
	Project      string `json:"project"`
	Reference    bool   `json:"reference"`
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	Version      string `json:"version"`
}

// CategoryParams represents the parameters to create or update a category in ReadMe.
type CategoryParams struct {
	// Title is a *required* short title for the category. This is what will show in the sidebar.
	Title string `json:"title"`
	// Type is tye type of category, which can be "reference" or "guide".
	Type string `json:"type"`
}

// CategorySaved represents the ReadMe API response when a category is created or updated.
type CategorySaved struct {
	CreatedAt string  `json:"createdAt"`
	ID        string  `json:"id"`
	Order     int     `json:"order"`
	Project   string  `json:"project"`
	Reference bool    `json:"reference"`
	Slug      string  `json:"slug"`
	Title     string  `json:"title"`
	Type      string  `json:"type"`
	Version   Version `json:"version"`
}

type CategoryVersionSaved struct {
	CreatedAt string          `json:"createdAt"`
	ID        string          `json:"id"`
	Order     int             `json:"order"`
	Project   string          `json:"project"`
	Reference bool            `json:"reference"`
	Slug      string          `json:"slug"`
	Title     string          `json:"title"`
	Type      string          `json:"type"`
	Version   CategoryVersion `json:"version"`
}

type CategoryVersion struct {
	Categories   []Category                `json:"categories"`
	Codename     string                    `json:"codename"`
	CreatedAt    string                    `json:"createdAt"`
	ForkedFrom   CategoryVersionForkedFrom `json:"forked_from"`
	ID           string                    `json:"id"`
	IsBeta       bool                      `json:"is_beta"`
	IsDeprecated bool                      `json:"is_deprecated"`
	IsHidden     bool                      `json:"is_hidden"`
	IsStable     bool                      `json:"is_stable"`
	Project      string                    `json:"project"`
	ReleaseDate  string                    `json:"releaseDate"`
	Version      string                    `json:"version"`
	VersionClean string                    `json:"version_clean"`
}

type CategoryVersionForkedFrom struct {
	ID      string `json:"_id"`
	Version string `json:"version"`
}

// CategoryDocs represents a document within a category.
type CategoryDocs struct {
	Children []CategoryDocsChildren `json:"children"`
	Hidden   bool                   `json:"hidden"`
	ID       string                 `json:"_id"`
	Order    int                    `json:"order"`
	Slug     string                 `json:"slug"`
	Title    string                 `json:"title"`
}

// CategoryDocsChildren represents a document's children within a category.
type CategoryDocsChildren struct {
	Hidden bool   `json:"hidden"`
	ID     string `json:"_id"`
	Order  int    `json:"order"`
	Slug   string `json:"slug"`
	Title  string `json:"title"`
}

// validChangelogType validates the 'type' field when creating or updating a changelog.
func validCategoryType(categoryType string) bool {
	if categoryType == "guide" || categoryType == "reference" {
		return true
	}

	return false
}

// GetAll retrieves and returns all categories on ReadMe.com.
//
// API Reference: https://docs.readme.com/reference/getcategories
func (c CategoryClient) GetAll(options ...RequestOptions) ([]Category, *APIResponse, error) {
	var err error
	var categories []Category
	var apiResponse *APIResponse
	hasNextPage := false

	// Initialize pagination counter.
	page := 1
	if len(options) > 0 {
		if options[0].Page != 0 {
			page = options[0].Page
		}
	}

	for {
		var paginatedResult []Category

		apiRequest := &APIRequest{
			Method:       "GET",
			Endpoint:     CategoryEndpoint,
			UseAuth:      true,
			OkStatusCode: []int{200},
			Response:     &paginatedResult,
		}
		if len(options) > 0 {
			apiRequest.RequestOptions = options[0]
		}

		apiResponse, hasNextPage, err = c.client.paginatedRequest(apiRequest, page)
		if err != nil {
			return categories, apiResponse, fmt.Errorf("unable to retrieve categories: %w", err)
		}
		categories = append(categories, paginatedResult...)

		if !hasNextPage {
			break
		}

		page = page + 1
	}

	return categories, apiResponse, nil
}

// Get a single category on ReadMe.com.
//
// The `category` parameter may be a slug or category ID prefixed with "id:".
// The slug is preferred, since it's a more direct request while the ID requires
// iterating over the list of all categories for a matching ID.
//
// API Reference: https://docs.readme.com/reference/getcategory
func (c CategoryClient) Get(category string, options ...RequestOptions) (Category, *APIResponse, error) {
	categoryResponse := Category{}

	opts := RequestOptions{}
	if len(options) > 0 {
		opts = options[0]
	}

	isID, paramID := ParseID(category)
	if isID {
		category = ""

		// Get all categories.
		categories, apiResponse, err := c.GetAll(opts)
		if err != nil {
			return categoryResponse, apiResponse, err
		}

		// Find a match by ID and return it.
		for _, cat := range categories {
			if cat.ID == paramID {
				category = cat.Slug
			}
		}
	}

	if category == "" {
		return categoryResponse, nil, errors.New("a category slug or id must be provided")
	}

	apiRequest := &APIRequest{
		Method:         "GET",
		Endpoint:       fmt.Sprintf("%s/%s", CategoryEndpoint, category),
		UseAuth:        true,
		OkStatusCode:   []int{200},
		Response:       &categoryResponse,
		RequestOptions: opts,
	}

	apiResponse, err := c.client.APIRequest(apiRequest)

	return categoryResponse, apiResponse, err
}

// GetDocs a list of docs metadata for a category on ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/getcategorydocs
func (c CategoryClient) GetDocs(slug string, options ...RequestOptions) ([]CategoryDocs, *APIResponse, error) {
	var response []CategoryDocs

	apiRequest := &APIRequest{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("%s/%s/docs", CategoryEndpoint, slug),
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

// Create a new category in ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/createcategory
func (c CategoryClient) Create(
	response any,
	params CategoryParams,
	options ...RequestOptions,
) (*APIResponse, error) {
	if !validCategoryType(params.Type) {
		return nil, fmt.Errorf("type must be 'guide' or 'reference'")
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("unable to parse request: %w", err)
	}

	apiRequest := &APIRequest{
		Method:       "POST",
		Endpoint:     CategoryEndpoint,
		UseAuth:      true,
		Payload:      payload,
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		OkStatusCode: []int{201},
		Response:     &response,
	}

	if len(options) > 0 {
		apiRequest.RequestOptions = options[0]
	}

	apiResponse, err := c.client.APIRequest(apiRequest)

	return apiResponse, err
}

// Update an existing category in ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/updatecategory
func (c CategoryClient) Update(
	slug string,
	params CategoryParams,
	options ...RequestOptions,
) (Category, *APIResponse, error) {
	if !validCategoryType(params.Type) {
		return Category{}, nil, fmt.Errorf("type must be 'guide' or 'reference'")
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return Category{}, nil, fmt.Errorf("unable to parse request: %w", err)
	}

	response := Category{}
	apiRequest := &APIRequest{
		Method:       "PUT",
		Endpoint:     fmt.Sprintf("%s/%s", CategoryEndpoint, slug),
		UseAuth:      true,
		Payload:      payload,
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		OkStatusCode: []int{200},
		Response:     &response,
	}

	if len(options) > 0 {
		apiRequest.RequestOptions = options[0]
	}

	apiResponse, err := c.client.APIRequest(apiRequest)

	return response, apiResponse, err
}

// Delete an existing category in ReadMe.
//
// API Reference: https://docs.readme.com/reference/deletecategory
func (c CategoryClient) Delete(slug string, options ...RequestOptions) (bool, *APIResponse, error) {
	apiRequest := &APIRequest{
		Method:       "DELETE",
		Endpoint:     fmt.Sprintf("%s/%s", CategoryEndpoint, slug),
		UseAuth:      true,
		OkStatusCode: []int{204},
	}
	if len(options) > 0 {
		apiRequest.RequestOptions = options[0]
	}

	apiResponse, err := c.client.APIRequest(apiRequest)
	if err != nil {
		return false, apiResponse, err
	}

	return true, apiResponse, nil
}
