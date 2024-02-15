package readme

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

// APISpecificationEndpoint is the ReadMe API Endpoint for API Specifications.
const APISpecificationEndpoint = "/api-specification"

// APISpecificationFormField is the name of the form field for uploading API Specifications.
const APISpecificationFormField = "spec"

// APISpecificationService is an interface for using the API Specification endpoints of the
// ReadMe.com API.
type APISpecificationService interface {
	// Create a new API specification on ReadMe by uploading a specification definition provided as
	// a JSON string or by associating an existing definition in the API registry by providing a
	// registry UUID as a parameter.
	//
	// The `definition` parameter can be a JSON string of the full definition or a string with the
	// UUID of an API registry in ReadMe prefixed with "uuid:". A UUID can be obtained from the
	// APIRegistry.Create() method's response.
	//
	// NOTE: specifying the definition as a UUID is an *undocumented* feature of the API.
	//
	// API Reference: https://docs.readme.com/reference/uploadapispecification
	Create(definition string, options ...RequestOptions) (APISpecificationSaved, *APIResponse, error)

	// Delete an API Specification by ID.
	// It returns true if it successfully deletes an API Specification.
	//
	// API Reference: https://docs.readme.com/reference/deleteapispecification
	Delete(specID string) (bool, *APIResponse, error)

	// Get a single API specification with a provided ID.
	//
	// Requesting a single API specification isn't included in the API. The client uses GetAll() to
	// retrieve the full list and matches the requested ID to the list.
	//
	// An error is returned if the specification wasn't found.
	//
	// API Reference: https://docs.readme.com/reference/getapispecification
	Get(specID string, options ...RequestOptions) (APISpecification, *APIResponse, error)

	// GetAll retrieves and returns all API specifications on ReadMe.com.
	//
	// API Reference: https://docs.readme.com/reference/getapispecification
	GetAll(...RequestOptions) ([]APISpecification, *APIResponse, error)

	// Update an existing API specification on ReadMe by uploading a specification definition
	// provided as a JSON string or by associating an existing definition in the API registry by
	// providing a registry UUID as a parameter.
	//
	// The `definition` parameter can be a JSON string of the full definition or a string with the
	// UUID of an API registry in ReadMe prefixed with "uuid:". A UUID can be obtained from the
	// APIRegistry.Create() method's response.
	//
	// NOTE: specifying the definition as a UUID is an *undocumented* feature of the API.
	//
	// API Reference: https://docs.readme.com/reference/updateapispecification
	Update(specID, definition string) (APISpecificationSaved, *APIResponse, error)

	// UploadDefinition uploads an API specification definition by making a request that submits
	// form data with the specification definition provided as a string.
	// APISpecification.Create() should be used in most cases instead of calling this directly.
	UploadDefinition(method, content, url, version string, response interface{}) (interface{}, *APIResponse, error)
}

// APISpecificationClient handles communication with the API specification related methods of the
// ReadMe.com API.
type APISpecificationClient struct {
	client *Client
}

// Ensure the implementation satisfies the expected interfaces.
var _ APISpecificationService = &APISpecificationClient{}

// APISpecification represents an API specification on ReadMe.com.
type APISpecification struct {
	Category   CategorySummary `json:"category"`
	ID         string          `json:"id"`
	LastSynced string          `json:"lastSynced"`
	Source     string          `json:"source"`
	Title      string          `json:"title"`
	Type       string          `json:"type"`
	Version    string          `json:"version"`
}

// APISpecificationSaved represents a successful response to creating an API specification
// on ReadMe.com.
type APISpecificationSaved struct {
	ID    string `json:"_id"`
	Title string `json:"title"`
}

// GetAll retrieves and returns all API specifications on ReadMe.com.
//
// API Reference: https://docs.readme.com/reference/getapispecification
func (c APISpecificationClient) GetAll(options ...RequestOptions) ([]APISpecification, *APIResponse, error) {
	var specifications []APISpecification
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
		var specPaginatedResult []APISpecification

		apiRequest := &APIRequest{
			Method:       "GET",
			Endpoint:     APISpecificationEndpoint,
			UseAuth:      true,
			OkStatusCode: []int{200},
			Response:     &specPaginatedResult,
		}
		if len(options) > 0 {
			apiRequest.RequestOptions = options[0]
		}

		apiResponse, hasNextPage, err = c.client.paginatedRequest(apiRequest, page)
		if err != nil {
			return specifications, apiResponse, fmt.Errorf("unable to retrieve specifications: %w", err)
		}
		specifications = append(specifications, specPaginatedResult...)

		if !hasNextPage {
			break
		}

		page = page + 1
	}

	return specifications, apiResponse, nil
}

// Get a single API specification with a provided ID.
//
// Requesting a single API specification isn't included in the API. The client uses GetAll() to
// retrieve the full list and matches the requested ID to the list.
//
// An error is returned if the specification wasn't found.
//
// See https://docs.readme.com/reference/getapispecification
func (c APISpecificationClient) Get(specID string, options ...RequestOptions) (APISpecification, *APIResponse, error) {
	specifications, apiResponse, err := c.GetAll(options...)
	if err != nil {
		return APISpecification{}, apiResponse, fmt.Errorf("unable to retrieve API specifications")
	}

	for _, specification := range specifications {
		if specification.ID == specID {
			return specification, apiResponse, nil
		}
	}

	return APISpecification{}, apiResponse, fmt.Errorf("API specification not found")
}

// Create a new API specification on ReadMe by uploading a specification definition provided as a
// JSON string or by associating an existing definition in the API registry by providing a registry
// UUID as a parameter.
//
// The `definition` parameter can be a JSON string of the full definition or a string with the UUID
// of an API registry in ReadMe prefixed with "uuid:". A UUID can be obtained from the
// APIRegistry.Create() method's response.
//
// NOTE: specifying the definition as a UUID is an *undocumented* feature of the API.
//
// See https://docs.readme.com/reference/uploadapispecification
func (c APISpecificationClient) Create(
	definition string,
	options ...RequestOptions,
) (APISpecificationSaved, *APIResponse, error) {
	version := ""

	if len(options) > 0 {
		version = options[0].Version
	}

	created, apiResponse, err := c.createOrUpdateSpec("POST", definition, version)
	if err != nil {
		return APISpecificationSaved{}, apiResponse, err
	}

	return *created, apiResponse, nil
}

// Update an existing API specification on ReadMe by uploading a specification definition provided
// as a JSON string or by associating an existing definition in the API registry by providing a
// registry UUID as a parameter.
//
// The `definition` parameter can be a JSON string of the full definition or a string with the UUID
// of an API registry in ReadMe prefixed with "uuid:". A UUID can be obtained from the
// APIRegistry.Create() method's response.
//
// NOTE: specifying the definition as a UUID is an *undocumented* feature of the API.
//
// API Reference: https://docs.readme.com/reference/updateapispecification
func (c APISpecificationClient) Update(
	specID, definition string,
) (APISpecificationSaved, *APIResponse, error) {
	updated, apiResponse, err := c.createOrUpdateSpec("PUT", definition, "", specID)
	if err != nil {
		return APISpecificationSaved{}, apiResponse, err
	}

	return *updated, apiResponse, nil
}

// Delete an API Specification by ID.
// It returns true if it successfully deletes an API Specification.
//
// API Reference: https://docs.readme.com/reference/deleteapispecification
func (c APISpecificationClient) Delete(specID string) (bool, *APIResponse, error) {
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "DELETE",
		Endpoint:     fmt.Sprintf("%s/%s", APISpecificationEndpoint, specID),
		UseAuth:      true,
		OkStatusCode: []int{204},
	})
	if err != nil {
		return false, apiResponse, err
	}

	return true, apiResponse, nil
}

// createOrUpdateSpec is a private method that handles creating a new specification or updating an
// existing specification. The Create() and Update() methods wrap this with the appropriate
// parameters. The `method` parameter should either be "POST" for creating new specifications or
// "PUT" for updating existing ones.
//
// The `specID` parameter is required if updating.
func (c APISpecificationClient) createOrUpdateSpec(
	method, definition, version string,
	specID ...string,
) (*APISpecificationSaved, *APIResponse, error) {
	var apiResponse *APIResponse
	var err error
	var url string

	// If a specID is provided, update the existing specification.
	if len(specID) == 1 {
		url = fmt.Sprintf("%s/%s", APISpecificationEndpoint, specID[0])
	} else {
		url = APISpecificationEndpoint
	}
	response := &APISpecificationSaved{}

	isUUID, uuid := ParseUUID(definition)
	if isUUID {
		_, apiResponse, err = c.createOrUpdateWithUUID(method, url, uuid, version, response)
	} else {
		_, apiResponse, err = c.UploadDefinition(method, definition, url, version, response)
	}

	return response, apiResponse, err
}

// UploadDefinition uploads an API specification definition by making a request that submits form
// data with the specification definition provided as a string.
//
// APISpecification.Create() should be used in most cases instead of calling this directly.
func (c APISpecificationClient) UploadDefinition(
	method, definition, url, version string,
	response interface{},
) (interface{}, *APIResponse, error) {
	data := strings.NewReader(definition)

	formData := &bytes.Buffer{}
	writer := multipart.NewWriter(formData)

	part, err := writer.CreateFormFile(APISpecificationFormField, "spec.json")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create request form: %w", err)
	}

	_, err = io.Copy(part, data)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to copy data: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to close writer: %w", err)
	}

	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:         method,
		Endpoint:       url,
		UseAuth:        true,
		Headers:        []RequestHeader{{"Content-Type": writer.FormDataContentType()}},
		Payload:        formData.Bytes(),
		OkStatusCode:   []int{200, 201},
		Response:       response,
		RequestOptions: RequestOptions{Version: version},
	})

	return &response, apiResponse, err
}

// createOrUpdateWithUUID handles publishing an API specification using a provided UUID of an
// existing API Registry as an alternative to creating it by uploading the specification.
// APISpecification.Create() should be used in most cases instead of calling this directly.
//
// NOTE: creating an API specification definition using this method is an *undocumented* feature of
// the API.
func (c APISpecificationClient) createOrUpdateWithUUID(
	method, url, uuid, version string,
	response interface{},
) (interface{}, *APIResponse, error) {
	apiRequest := &APIRequest{
		Method:         method,
		Endpoint:       url,
		UseAuth:        true,
		Headers:        []RequestHeader{{"Content-Type": "application/json"}},
		Payload:        []byte(`{"registryUUID": "` + uuid + `"}`),
		OkStatusCode:   []int{200, 201},
		Response:       &response,
		RequestOptions: RequestOptions{Version: version},
	}
	apiResponse, err := c.client.APIRequest(apiRequest)

	return response, apiResponse, err
}
