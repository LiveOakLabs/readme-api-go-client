package readme

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// DocEndpoint is the ReadMe API endpoint for docs.
const DocEndpoint = "/docs"

// DocService is an interface for using the docs endpoints of the ReadMe.com API.
//
// API Reference: https://docs.readme.com/main/reference/getdoc
type DocService interface {
	// Create a new doc in ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/createdoc
	Create(params DocParams, options ...RequestOptions) (Doc, *APIResponse, error)

	// Delete a doc in ReadMe.
	//
	// API Reference: https://docs.readme.com/reference/deletedoc
	Delete(slug string, options ...RequestOptions) (bool, *APIResponse, error)

	// Get a doc from ReadMe.
	//
	// The `doc` parameter may be a slug or doc ID prefixed with "id:".
	// The slug is preferred, since it's a more direct request while the ID requires
	// iterating over the search results for the matching ID.
	//
	// Use the `options` parameter to set `RequestOptions.ProductionDoc` to retrieve a production doc.
	//
	// API References:
	//   - https://docs.readme.com/main/reference/getdoc
	//   - https://docs.readme.com/main/reference/getproductiondoc
	Get(doc string, options ...RequestOptions) (Doc, *APIResponse, error)

	// Search for docs that match the search query parameter.
	//
	// API Reference: https://docs.readme.com/main/reference/searchdocs
	Search(query string, options ...RequestOptions) ([]DocSearchResult, *APIResponse, error)

	// Update an existing doc in ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/updatedoc
	Update(slug string, params DocParams, options ...RequestOptions) (Doc, *APIResponse, error)
}

// DocClient handles communication with the docs related methods of the ReadMe.com API.
type DocClient struct {
	client *Client
}

// Ensure the implementation satisfies the expected interfaces.
// This is a compile-time check.
// See: https://golang.org/doc/faq#guarantee_satisfies_interface
var _ DocService = &DocClient{}

// Doc represents a doc in ReadMe.
type Doc struct {
	Algolia       DocAlgolia     `json:"algolia"`
	API           DocAPI         `json:"api"`
	Body          string         `json:"body"`
	BodyHTML      string         `json:"body_html,omitempty"`
	Category      string         `json:"category"`
	CreatedAt     string         `json:"createdAt"`
	Deprecated    bool           `json:"deprecated"`
	Error         DocErrorObject `json:"error"`
	Excerpt       string         `json:"excerpt"`
	Hidden        bool           `json:"hidden"`
	ID            string         `json:"_id"`
	Icon          string         `json:"icon"`
	IsAPI         bool           `json:"isApi"`
	IsReference   bool           `json:"isReference"`
	LinkExternal  bool           `json:"link_external"`
	LinkURL       string         `json:"link_url"`
	Metadata      DocMetadata    `json:"metadata"`
	Next          DocNext        `json:"next"`
	Order         int            `json:"order"`
	ParentDoc     string         `json:"parentDoc"`
	PreviousSlug  string         `json:"previousSlug"`
	Project       string         `json:"project,omitempty"`
	Revision      int            `json:"revision"`
	Slug          string         `json:"slug"`
	SlugUpdatedAt string         `json:"slugUpdatedAt"`
	Swagger       DocSwagger     `json:"swagger"`
	SyncUnique    string         `json:"sync_unique"`
	Title         string         `json:"title"`
	Type          string         `json:"type"`
	UpdatedAt     string         `json:"updatedAt"`
	// TODO: Verify the data structure and type.
	Updates []any  `json:"updates"`
	User    string `json:"user"`
	Version string `json:"version"`
}

// DocUser represents the corresponding key in the API response for a user
// object on ReadMe. The 'User' field is returned as a string on 'doc' objects,
// but is returned as an object on 'changelog' objects.
type DocUser struct {
	ID string `json:"_id"`
}

// DocAlgolia represents the corresponding 'algolia' key in the API response for a doc.
type DocAlgolia struct {
	PublishPending bool   `json:"publishPending"`
	RecordCount    int    `json:"recordCount"`
	UpdatedAt      string `json:"updatedAt"`
}

// DocMetadata represents a doc's metadata in ReadMe.
// This is used across the different types of docs - changelog, doc (aka guides) and custom pages.
type DocMetadata struct {
	Description string `json:"description"`
	Image       []any  `json:"image"`
	Title       string `json:"title"`
}

// DocAPI represents a doc's API data in ReadMe.
type DocAPI struct {
	APISetting string         `json:"apiSetting"`
	Auth       string         `json:"auth"`
	Examples   DocAPIExamples `json:"examples"`
	Method     string         `json:"method"`
	Params     []DocAPIParams `json:"params"`
	Results    DocAPIResults  `json:"results"`
	URL        string         `json:"url"`
}

// DocAPIExamples represents the "api:examples" object returned for a doc.
type DocAPIExamples struct {
	Codes []DocAPIExamplesCodes `json:"codes"`
}

// DocAPIExamplesCodes represents the "api:examples:codes" object returned for a doc.
type DocAPIExamplesCodes struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

// DocNext represents a doc's "next" doc and pages in ReadMe.
type DocNext struct {
	Description string         `json:"description"`
	Pages       []DocNextPages `json:"pages"`
}

// DocNextPages represents the "next:pages" object returned for a doc.
type DocNextPages struct {
	Category   string `json:"category"`
	Deprecated bool   `json:"deprecated"`
	Icon       string `json:"icon"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Type       string `json:"type"`
}

// DocAPIParams represents the "api:params" object returned for a doc.
type DocAPIParams struct {
	Default    string `json:"default"`
	Desc       string `json:"desc"`
	EnumValues string `json:"enumValues"`
	ID         string `json:"_id"`
	In         string `json:"in"`
	Name       string `json:"name"`
	Ref        string `json:"ref"`
	Required   bool   `json:"required"`
	Type       string `json:"type"`
}

// DocAPIResults represents the "api:results" object returned for a doc.
type DocAPIResults struct {
	Codes []DocAPIResultsCodes `json:"codes"`
}

// DocAPIResultsCodes represents the "api:results:codes" object returned for a doc.
type DocAPIResultsCodes struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
}

// DocSwagger represents the optional "swagger" object returned for a doc.
type DocSwagger struct {
	Path string `json:"path"`
}

// DocParams represents the parameters for creating a doc on ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/createdoc
type DocParams struct {
	// Body content of the page, formatted in ReadMe or GitHub flavored Markdown. Accepts long page
	// content, for example, greater than 100k characters.
	Body string `json:"body,omitempty"`
	// Category ID of the page.
	// This is or 'CategorySlug' is *required* when creating or updating a category.
	Category string `json:"category"`
	// Category Slug of the page.
	// This is or 'Category' is *required* when creating or updating a category.
	CategorySlug string `json:"categorySlug"`
	// Error is an error code for docs with the type set to "error".
	Error DocErrorObject `json:"error,omitempty"`
	// Hidden toggles visibility for the doc.
	// API default is true.
	Hidden *bool `json:"hidden"`
	// Order sets the position of the page in the project sidebar.
	// API default is 999.
	Order *int `json:"order"`
	// ParentDoc is the ID of the parent doc for subpages.
	ParentDoc string `json:"parentDoc,omitempty"`
	// ParentDocSlug is the slug of the parent doc for subpages.
	// This field is an alternative to the ParentDoc field.
	ParentDocSlug string `json:"parentDocSlug,omitempty"`
	// Title of the page.
	// This is *required* when creating or updating a category.
	Title string `json:"title"`
	// Type of the page. The available types all show up under the /docs/ URL path of your docs
	// project (also known as the "guides" section). Can be "basic" (most common), "error" (page
	// desribing an API error), or "link" (page that redirects to an external link).
	Type string `json:"type,omitempty"`
}

// DocErrorObject represents the 'error' key in a doc response.
type DocErrorObject struct {
	Code string `json:"code"`
}

// DocSearchResults represents the response from searching.
type DocSearchResults struct {
	Results []DocSearchResult `json:"results"`
}

// DocSearchResult represents a single item in the list of search results.
type DocSearchResult struct {
	HighlightResult DocSearchResultHighlight `json:"_highlightResult"`
	IndexName       string                   `json:"indexName"`
	InternalLink    string                   `json:"internalLink"`
	IsReference     bool                     `json:"isReference"`
	LinkURL         string                   `json:"link_url"`
	Method          string                   `json:"method"`
	ObjectID        string                   `json:"objectID"`
	Project         string                   `json:"project"`
	ReferenceID     string                   `json:"referenceId"`
	Slug            string                   `json:"slug"`
	SnippetResult   DocSearchResultSnippet   `json:"_snippetResult"`
	Subdomain       string                   `json:"subdomain"`
	Title           string                   `json:"title"`
	Type            string                   `json:"type"`
	URL             string                   `json:"url"`
	Version         string                   `json:"version"`
}

// DocSearchResultHighlight represents the HighlightResult key in a search result item.
type DocSearchResultHighlight struct {
	Title   DocSearchResultHighlightValue `json:"title"`
	Excerpt DocSearchResultHighlightValue `json:"excerpt"`
	Body    DocSearchResultHighlightValue `json:"body"`
}

// DocSearchResultHighlightValue represents the HighlightResult child keys in a search result item.
type DocSearchResultHighlightValue struct {
	Value        string   `json:"value"`
	MatchLevel   string   `json:"matchLevel"`
	MatchedWords []string `json:"matchedWords"`
}

// DocSearchResultSnippet represents the SnippetResult key in a search result item.
type DocSearchResultSnippet struct {
	Title   DocSearchResultSnippetValue `json:"title"`
	Excerpt DocSearchResultSnippetValue `json:"excerpt"`
	Body    DocSearchResultSnippetValue `json:"body"`
}

// DocSearchResultSnippetValue represents the SnippetResult child keys in a search result item.
type DocSearchResultSnippetValue struct {
	Value      string `json:"value"`
	MatchLevel string `json:"matchLevel"`
}

// Get a doc from ReadMe.
//
// The `doc` parameter may be a slug or doc ID prefixed with "id:".
// The slug is preferred, since it's a more direct request while the ID requires
// iterating over the search results for the matching ID.
//
// Use the `options` parameter to set `RequestOptions.ProductionDoc` to retrieve a production doc.
//
// API References:
//   - https://docs.readme.com/main/reference/getdoc
//   - https://docs.readme.com/main/reference/getproductiondoc
func (c DocClient) Get(doc string, options ...RequestOptions) (Doc, *APIResponse, error) {
	response := Doc{}

	opts := RequestOptions{}
	if len(options) > 0 {
		opts = options[0]
	}

	isID, paramID := ParseID(doc)
	if isID {
		// Reset the doc query. A matching search result will set this to its slug.
		doc = ""

		// Search all docs.
		docs, apiResponse, err := c.Search(paramID, opts)
		if err != nil {
			return response, apiResponse, err
		}

		// Find a match by ID and return it.
		for _, docResult := range docs {
			if docResult.ReferenceID == paramID {
				doc = docResult.Slug
			}
		}

		if doc == "" {
			return response, nil, fmt.Errorf("no doc found matching id %s (is it hidden?)", paramID)
		}
	}

	if doc == "" {
		return response, nil, errors.New("a doc slug or id must be provided")
	}

	apiRequest := &APIRequest{
		Method:         "GET",
		Endpoint:       fmt.Sprintf("%s/%s", DocEndpoint, doc),
		UseAuth:        true,
		OkStatusCode:   []int{200},
		Response:       &response,
		RequestOptions: opts,
	}

	if opts.ProductionDoc {
		apiRequest.Endpoint = fmt.Sprintf("%s/production", apiRequest.Endpoint)
	}

	apiResponse, err := c.client.APIRequest(apiRequest)

	return response, apiResponse, err
}

// Create a new doc in ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/createdoc
func (c DocClient) Create(params DocParams, options ...RequestOptions) (Doc, *APIResponse, error) {
	if params.Title == "" {
		return Doc{}, nil, fmt.Errorf("doc title is required")
	}

	if params.Category == "" && params.CategorySlug == "" {
		return Doc{}, nil, fmt.Errorf(("doc category or category slug is required"))
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return Doc{}, nil, fmt.Errorf("unable to marshal request: %w", err)
	}

	response := Doc{}
	apiRequest := &APIRequest{
		Endpoint:     DocEndpoint,
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		Method:       "POST",
		OkStatusCode: []int{201},
		Payload:      payload,
		Response:     &response,
		UseAuth:      true,
	}

	if len(options) > 0 {
		apiRequest.RequestOptions = options[0]
	}

	apiResponse, err := c.client.APIRequest(apiRequest)

	return response, apiResponse, err
}

// Update an existing doc in ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/updatedoc
func (c DocClient) Update(slug string, params DocParams, options ...RequestOptions) (Doc, *APIResponse, error) {
	if params.Title == "" {
		return Doc{}, nil, fmt.Errorf("doc title is required")
	}

	if params.Category == "" && params.CategorySlug == "" {
		return Doc{}, nil, fmt.Errorf(("doc category or category slug is required"))
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return Doc{}, nil, fmt.Errorf("unable to marshal request: %w", err)
	}

	response := Doc{}
	apiRequest := &APIRequest{
		Endpoint:     fmt.Sprintf("%s/%s", DocEndpoint, slug),
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		Method:       "PUT",
		OkStatusCode: []int{200},
		Payload:      payload,
		Response:     &response,
		UseAuth:      true,
	}

	if len(options) > 0 {
		apiRequest.RequestOptions = options[0]
	}

	apiResponse, err := c.client.APIRequest(apiRequest)

	return response, apiResponse, err
}

// Delete a doc in ReadMe.
//
// API Reference: https://docs.readme.com/reference/deletedoc
func (c DocClient) Delete(slug string, options ...RequestOptions) (bool, *APIResponse, error) {
	apiRequest := &APIRequest{
		Method:       "DELETE",
		Endpoint:     fmt.Sprintf("%s/%s", DocEndpoint, slug),
		UseAuth:      true,
		OkStatusCode: []int{204},
	}

	if len(options) > 0 && options[0].Version != "" {
		apiRequest.Version = options[0].Version
	}

	apiResponse, err := c.client.APIRequest(apiRequest)
	if err != nil {
		return false, apiResponse, err
	}

	return true, apiResponse, nil
}

// Search for docs that match the search query parameter.
//
// API Reference: https://docs.readme.com/main/reference/searchdocs
func (c DocClient) Search(query string, options ...RequestOptions) ([]DocSearchResult, *APIResponse, error) {
	results := DocSearchResults{}
	apiRequest := &APIRequest{
		Method:       "POST",
		Endpoint:     fmt.Sprintf("%s/search?search=%s", DocEndpoint, url.QueryEscape(query)),
		UseAuth:      true,
		OkStatusCode: []int{200},
		Response:     &results,
	}

	if len(options) > 0 && options[0].Version != "" {
		apiRequest.Version = options[0].Version
	}

	apiResponse, err := c.client.APIRequest(apiRequest)

	return results.Results, apiResponse, err
}
