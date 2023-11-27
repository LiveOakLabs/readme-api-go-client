package readme

import (
	"fmt"
)

// APIRegistryEndpoint is the ReadMe API URL endpoint for the API registry.
const APIRegistryEndpoint = "/api-registry"

// APIRegistryService is an interface for interacting with the API Registry endpoints of the
// ReadMe.com API.
type APIRegistryService interface {
	// Create a new API registry on ReadMe.
	//
	// The response returns the UUID and the specification definition.
	//
	// NOTE: This is an undocumented endpoint on ReadMe and was discovered by inspecting the ReadMe
	// CLI application. The registry UUID is required for retrieving the remote specification.
	// A typical workflow will be to create the registry with this method and follow-up with a call
	// to APISpecification.Create() with the UUID returned from this response.
	Create(definition string, version ...string) (APIRegistrySaved, *APIResponse, error)

	// Get retrieves an API definition from the ReadMe.com API registry with a provided UUID and
	// returns it as a string.
	//
	// API Reference: https://docs.readme.com/main/reference/getapiregistry
	Get(uuid string) (string, *APIResponse, error)
}

// APIRegistryClient handles communication with the Registry related methods of the ReadMe.com API.
type APIRegistryClient struct {
	client *Client
}

// APIRegistrySaved represents the API response when an API Registry is created.
type APIRegistrySaved struct {
	Definition   map[string]interface{} `json:"definition"`
	RegistryUUID string                 `json:"registryUUID"`
}

// Ensure the implementation satisfies the expected interfaces.
var _ APIRegistryService = &APIRegistryClient{}

// Get retrieves an API definition from the ReadMe.com API registry with a provided UUID and returns
// it as a string.
//
// API Reference: https://docs.readme.com/main/reference/getapiregistry
func (c APIRegistryClient) Get(uuid string) (string, *APIResponse, error) {
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("%s/%s", APIRegistryEndpoint, uuid),
		UseAuth:      true,
		OkStatusCode: []int{200},
	})
	if err != nil {
		return "", apiResponse, err
	}

	spec := string(apiResponse.Body)

	return spec, apiResponse, nil
}

// Create a new API registry on ReadMe.
//
// The response returns the UUID and the specification definition.
//
// NOTE: This is an undocumented endpoint on ReadMe and was discovered by inspecting the ReadMe CLI
// application.
//
// The registry UUID is required for retrieving the remote specification.
// A typical workflow will be to create the registry with this method and follow-up with a call to
// APISpecification.Create() with the UUID returned from this response.
func (c APIRegistryClient) Create(definition string, version ...string) (APIRegistrySaved, *APIResponse, error) {
	var vers string
	if len(version) > 0 {
		vers = version[0]
	}

	response := APIRegistrySaved{}
	_, apiResponse, err := c.client.APISpecification.UploadDefinition(
		"POST",
		definition,
		APIRegistryEndpoint,
		vers,
		&response,
	)

	return response, apiResponse, err // nolint:wrapcheck
}
