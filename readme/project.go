package readme

// ProjectEndpoint is the ReadMe API URL endpoint for Project metadata.
const ProjectEndpoint = "/"

// ProjectService is an interface for using the projects endpoints of the ReadMe.com API.
//
// API Reference: https://docs.readme.com/main/reference/getproject
type ProjectService interface {
	// Get retrieves project metadata from the ReadMe.com API.
	//
	// API Reference: https://docs.readme.com/main/reference/getproject
	Get() (Project, *APIResponse, error)
}

// ProjectClient handles communication with the Project related methods of the ReadMe.com API.
type ProjectClient struct {
	client *Client
}

// Project represents the response from the ReadMe API when returning project metadata.
type Project struct {
	// BaseURL is the base URL for the project. If the project is not running under a custom domain,
	// it will be https://projectSubdomain.readme.io, otherwise it can either be https://example.com
	// or, in the case of an enterprise child project https://example.com/projectSubdomain.
	BaseURL string `json:"baseUrl"`
	// JWTSecret is the JSON Web Token used for the project.
	JWTSecret string `json:"jwtSecret"`
	// Name of the project.
	Name string `json:"name"`
	// Plan is the subscription plan of the project on ReadMe.
	Plan string `json:"plan"`
	// SubDomain for the project on ReadMe.com.
	SubDomain string `json:"subdomain"`
}

// Ensure the implementation satisfies the expected interfaces.
var _ ProjectService = &ProjectClient{}

// Get retrieves project metadata from the ReadMe.com API.
//
// API Reference: https://docs.readme.com/main/reference/getproject
func (c ProjectClient) Get() (Project, *APIResponse, error) {
	project := Project{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "GET",
		Endpoint:     ProjectEndpoint,
		UseAuth:      true,
		OkStatusCode: []int{200},
		Response:     &project,
	})

	return project, apiResponse, err
}
