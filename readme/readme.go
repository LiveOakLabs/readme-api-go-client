// ReadMe API Client for Go is for performing API operations with ReadMe.com.
//
// Refer to https://docs.readme.com/main/reference/intro/getting-started for more information about
// the ReadMe API.
package readme

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Markdown documentation generation.
//go:generate go run github.com/princjef/gomarkdoc/cmd/gomarkdoc --output ../docs/README.md ./...

const (
	// PaginationHeader is the name of the HTTP response header with pagination links.
	PaginationHeader = "link"

	// ReadmeAPIURL is the default base URL for the ReadMe API.
	ReadmeAPIURL = "https://dash.readme.com/api/v1"

	// TotalCountHeader is the name of the HTTP response header with the total count in results.
	TotalCountHeader = "x-total-count"

	// UserAgent is the name of the HTTP UserAgent when making requests.
	UserAgent = "readme-api-go-client"
)

// IDValidCharacters is a compiled RegEx pattern that matches valid characters in an object ID or
// API Registry UUID.
var IDValidCharacters = regexp.MustCompile("^[0-9a-zA-Z]+$")

// Client sets up the API HTTP client with authentication and exposes the API interfaces.
type Client struct {
	// APIURL is the base URL for the ReadMe API.
	APIURL string
	// HTTPClient is the initialized HTTP client.
	HTTPClient *http.Client
	// Token is the API token for authenticating with ReadMe.
	Token string

	// APIRegistry implements the ReadMe API Registry API for managing API definitions.
	APIRegistry APIRegistryService
	// APISpecification implements the ReadMe API Specification API for managing API specifications.
	APISpecification APISpecificationService
	// Apply implements the ReadMe API Apply API for retrieving and applying for positions at ReadMe.
	Apply ApplyService
	// Category implements the ReadMe Category API for managing categories.
	Category CategoryService
	// Changelog implements the ReadMe Changelog API for managing changelogs.
	Changelog ChangelogService
	// CustomPage implements the ReadMe CustomPage API for managing custom pages.
	CustomPage CustomPageService
	// Doc implements the ReadMe Docs API for managing docs.
	Doc DocService
	// Image implements the ReadMe Image API for uploading images.
	Image ImageService
	// Project implements the ReadMe Project API for retrieving metadata about the project.
	Project ProjectService
	// Version implements the ReadMe Version API for managing versions.
	Version VersionService
}

// RequestHeader represents an HTTP header set on requests.
type RequestHeader map[string]string

// APIRequest represents a request to the ReadMe.com API.
type APIRequest struct {
	// Endpoint is the API endpoint (after the base URL) for the request.
	Endpoint string

	// Headers lists HTTP headers to send in the request, in addition to the implicit headers.
	Headers []RequestHeader

	// Slice of HTTP status codes that are considered 'ok'.
	// Any other status code in the response results in an error.
	OkStatusCode []int

	// Method is the HTTP method to use for the request.
	Method string

	// An optional payload, in bytes, for the request.
	Payload []byte

	// Optional options for a request, including headers, version and pagination options.
	RequestOptions

	// Interface of a struct to map the response body to.
	Response interface{}

	// UseAuth toggles whether the request should use authentication or not.
	UseAuth bool

	// URL is a full URL string to use for the request as an alternative to Endpoint.
	URL string
}

// APIResponse represents the response from a request to the ReadMe API.
type APIResponse struct {
	// APIErrorResponse is a structured error from the ReadMe API when a request results in error.
	APIErrorResponse APIErrorResponse
	// Body is the response body in bytes.
	Body []byte
	// HTTPResponse is the stdlib http.Response type.
	HTTPResponse *http.Response
	// Request is the APIRequest struct used to create the request.
	Request *APIRequest
}

// APIErrorResponse represents the response ReadMe provides in the body of requests that failed.
type APIErrorResponse struct {
	// Docs is a ReadMe Metrics log URL where more information about the request can be retrieved.
	// If metrics URLs are unavailable for the request, this URL will be a URL to the ReadMe API Reference.
	Docs string `json:"docs"`
	// Error is an error code unique to the error received.
	Error string `json:"error"`
	// Help is information on where additional assistance from the ReadMe support team can be obtained.
	Help string `json:"help"`
	// Message is the reason why the error occurred.
	Message string `json:"message"`
	// Poem is a short poem about the error.
	Poem []string `json:"poem"`
	// Suggestion is a helpful suggestion for how to alleviate the error.
	Suggestion string `json:"suggestion"`
}

// RequestOptions is used for specifying options for requests, such as pagination options.
type RequestOptions struct {
	// Headers is a list of additional headers to add to the request.
	Headers []RequestHeader
	// PerPage is the number of items to return in each request when using pagination.
	// The maximum and default is 100.
	PerPage int
	// Page is the page number to request when using pagination.
	Page int
	// ProductionDoc is used by readme.Docs.Get() to indicate whether the requested document is a
	// 'production' doc.
	ProductionDoc bool
	// Version number of a ReadMe project, for example, v3.0. By default the main project version is used.
	Version string
}

// NewClient initializes the API client configuration and returns the HTTP client with an auth token and URL set.
//
// Optionally provide a custom API URL as a second parameter.
func NewClient(token string, apiURL ...string) (*Client, error) {
	client := &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
	client.APIURL = ReadmeAPIURL
	client.Token = token

	if apiURL != nil {
		if len(apiURL) > 1 {
			return nil, fmt.Errorf("unable to configure ReadMe API client: "+
				"too many values specified for API URL (got: %v; expects 1)", len(apiURL))
		}
		client.APIURL = apiURL[0]
	}

	client.APIRegistry = &APIRegistryClient{client: client}
	client.APISpecification = &APISpecificationClient{client: client}
	client.Apply = &ApplyClient{client: client}
	client.Category = &CategoryClient{client: client}
	client.Changelog = &ChangelogClient{client: client}
	client.CustomPage = &CustomPageClient{client: client}
	client.Doc = &DocClient{client: client}
	client.Image = &ImageClient{client: client}
	client.Project = &ProjectClient{client: client}
	client.Version = &VersionClient{client: client}

	return client, nil
}

// APIRequest performs a request to the ReadMe API and handles parsing the response and API errors.
//
// This function is called directly by the receiver functions used to implement each endpoint.
func (c *Client) APIRequest(request *APIRequest) (*APIResponse, error) {
	// Perform the request
	body, httpResponse, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	apiResponse := &APIResponse{
		Body:         body,
		HTTPResponse: &httpResponse,
		Request:      request,
	}

	// Verify the HTTP response from the API.
	apiErrorResponse, err := checkResponseStatus(body, httpResponse.StatusCode, request.OkStatusCode)
	if err != nil {
		apiResponse.APIErrorResponse = apiErrorResponse

		return apiResponse, err
	}

	// Parse the response into the specified interface.
	if request.Response != nil {
		err = json.Unmarshal(body, &request.Response)
		if err != nil {
			return apiResponse, fmt.Errorf("unable to parse API response: %w", err)
		}
	}

	err = httpResponse.Body.Close()
	if err != nil {
		return apiResponse, fmt.Errorf("problem closing HTTP response body")
	}

	return apiResponse, nil
}

// doRequest performs an API request and returns the response or error.
func (c *Client) doRequest(request *APIRequest) ([]byte, http.Response, error) {
	req, err := c.prepareRequest(request)
	if err != nil {
		return nil, http.Response{}, err
	}

	// Perform the request.
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, http.Response{}, fmt.Errorf("unable to make request: %w", err)
	}

	if res.Body == nil {
		return nil, *res, fmt.Errorf("response body is nil in %s request to %s", req.Method, req.URL)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, *res, fmt.Errorf("unable to read response: %w", err)
	}

	err = res.Body.Close()
	if err != nil {
		return nil, *res, fmt.Errorf("problem closing HTTP response body")
	}

	return body, *res, nil
}

// checkResponseStatus compares an HTTP response status code against a slice of 'OK' status codes.
//
// If the response code matches a provided code listed in okCodes, no error is returned.
// If the response code doesn't match, an error and APIErrorResponse is returned.
func checkResponseStatus(body []byte, responseCode int, okCodes []int) (APIErrorResponse, error) {
	var apiErrorResponse APIErrorResponse
	for _, okCode := range okCodes {
		if responseCode == okCode {
			return apiErrorResponse, nil
		}
	}

	err := json.Unmarshal(body, &apiErrorResponse)
	if err != nil {
		return apiErrorResponse, fmt.Errorf("unable to decode API error response: %w", err)
	}

	return apiErrorResponse, fmt.Errorf("API responded with a non-OK status: %v", responseCode)
}

// prepareRequest prepares an http.Request for the ReadMe API.
//
// This sets common headers and prepares an optional payload for the request.
func (c *Client) prepareRequest(request *APIRequest) (*http.Request, error) {
	// Prepare the request.
	if request.URL == "" {
		request.URL = c.APIURL + request.Endpoint
	}
	req, reqErr := http.NewRequest(request.Method, request.URL, nil)

	if request.Payload != nil {
		data := bytes.NewBuffer(request.Payload)
		req, reqErr = http.NewRequest(request.Method, request.URL, data)
	}

	if reqErr != nil {
		return nil, fmt.Errorf("unable to prepare request: %w", reqErr)
	}

	for _, r := range request.Headers {
		for header, value := range r {
			req.Header.Set(header, value)
		}
	}

	if request.UseAuth {
		encodedToken := base64.StdEncoding.EncodeToString([]byte(c.Token))
		authHeader := "Basic " + encodedToken
		req.Header.Set("authorization", authHeader)
	}

	if request.RequestOptions.Version != "" {
		req.Header.Set("x-readme-version", request.RequestOptions.Version)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	return req, nil
}

// paginatedRequest makes a request to the ReadMe API with pagination query parameters set.
//
// An abbreviated *APIRequest struct should be passed, leaving the Headers and Version fields unset.
// These are derived from the RequestOptions field.
//
// This function is intended to be called within a loop and returns the APIResponse struct and a
// boolean indicating if there is a next page indicated in the pagination header.
func (c *Client) paginatedRequest(apiRequest *APIRequest, page int) (*APIResponse, bool, error) {
	// Set default values
	perPage := 100

	// Check for custom values in RequestOptions
	if apiRequest.RequestOptions.PerPage != 0 {
		perPage = apiRequest.RequestOptions.PerPage
	}
	if apiRequest.RequestOptions.Headers != nil {
		apiRequest.Headers = apiRequest.RequestOptions.Headers
	}

	// Add pagination parameters to endpoint
	baseEndpoint := apiRequest.Endpoint
	apiRequest.Endpoint = fmt.Sprintf("%s?perPage=%d&page=%d", baseEndpoint, perPage, page)

	if apiRequest.URL == "" {
		apiRequest.URL = c.APIURL + apiRequest.Endpoint
	}

	// Make API request
	apiResponse, err := c.APIRequest(apiRequest)
	if err != nil {
		return apiResponse, false, fmt.Errorf("unable to make request: %w", err)
	}

	// Check for next page
	hasNextPage, err := HasNextPage(apiResponse.HTTPResponse.Header.Get(PaginationHeader))
	if err != nil {
		return apiResponse, false, fmt.Errorf(
			"unable to check pagination link header '%s': %w; ",
			PaginationHeader,
			err,
		)
	}
	if !hasNextPage {
		return apiResponse, false, nil
	}

	// Get total count of items
	totalCountHeader := apiResponse.HTTPResponse.Header.Get(TotalCountHeader)
	totalCount, err := strconv.Atoi(totalCountHeader)
	if err != nil {
		return apiResponse, false, fmt.Errorf(
			"unable to parse '%s' header: %w; Response: %v",
			TotalCountHeader,
			err,
			apiResponse,
		)
	}

	// Check if current page is last page
	if page >= (totalCount / perPage) {
		return apiResponse, false, nil
	}

	return apiResponse, true, nil
}

// HasNextPage checks if a "next" link is provided in the "links" response header for pagination,
// indicating the request has a next page.
//
// This does a rudimentary parsing of the header value, splitting on the comma-separated links and
// parsing the value of "rel".
//
// A link header looks like:
// </api-specification?page=2>; rel="next", <>; rel="prev", <>; rel="last"
func HasNextPage(links string) (bool, error) {
	// Split links by comma
	parts := strings.Split(links, ",")

	// Return error if invalid format
	if len(parts) < 3 {
		return false, fmt.Errorf("unable to parse link header - invalid format: "+
			"'%s'; expected "+`'<>; rel="next", <>; rel="prev", <>; rel="last"'`, links)
	}

	// Check for "rel=next" in parts
	for _, part := range parts {
		rel := strings.Split(part, ";")
		if len(rel) != 2 {
			return false, fmt.Errorf("unable to parse link header - invalid format: "+
				"'%s'; expected "+`'<>; rel="next", <>; rel="prev", <>; rel="last"'`, links)
		}
		if rel[1] == " rel=\"next\"" && rel[0] != "<>" {
			return true, nil
		}
	}

	// Return false if "rel=next" is not found
	return false, nil
}

// ValidateID is a helper script for parseID() and parseUUID() that checks a string to determine if
// it appears to be a valid ReadMe API object ID or Registry UUID.
func ValidateID(id, prefix string, min_len, max_len int) (bool, string) {
	if !strings.HasPrefix(id, prefix+":") {
		return false, ""
	}

	parts := strings.Split(id, ":")

	if len(parts[1]) < min_len || len(parts[1]) > max_len {
		return false, ""
	}

	return IDValidCharacters.MatchString(parts[1]), parts[1]
}

// ParseUUID checks a string to determine if it appears to be a valid ReadMe API Registry UUID.
//
// The provided parameter should be a ReadMe API Registry UUID prefixed with "uuid:".
//
// NOTE: The min and max lengths aren't certain or documented in the API. The UUID length varies.
func ParseUUID(uuid string) (bool, string) {
	return ValidateID(uuid, "uuid", 10, 24)
}

// ParseID checks a string to determine if it appears to be a valid ReadMe API object ID.
//
// The provided parameter should be a ReadMe API object ID prefixed with "id:".
//
// NOTE: The min and max lengths aren't certain or documented in the API.
func ParseID(id string) (bool, string) {
	return ValidateID(id, "id", 20, 24)
}
