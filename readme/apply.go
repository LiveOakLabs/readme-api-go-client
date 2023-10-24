package readme

import (
	"encoding/json"
	"fmt"
)

// ApplyEndpoint is the ReadMe API URL endpoint for listing open positions at ReadMe and applying
// for a them.
const ApplyEndpoint = "/apply"

// ApplyService is an interface for interacting with the Apply endpoints of the ReadMe.com API.
//
// API Reference: https://docs.readme.com/main/reference/getopenroles
type ApplyService interface {
	// Apply for an open role at ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/applytoreadme
	Apply(application Application) (ApplyResponse, *APIResponse, error)

	// Get a list of open roles at ReadMe.
	//
	// API Reference: https://docs.readme.com/main/reference/getopenroles
	Get() ([]OpenRole, *APIResponse, error)
}

// ApplyClient handles communication with the Apply related methods of the ReadMe.com API.
type ApplyClient struct {
	client *Client
}

// Ensure the implementation satisfies the expected interfaces.
var _ ApplyService = &ApplyClient{}

// OpenRole represents an open role at ReadMe.
type OpenRole struct {
	Department  string `json:"department"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Pullquote   string `json:"pullquote"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

// Application represents the parameters used for submitting an application to ReadMe.
type Application struct {
	// CoverLetter is additional information for the application.
	CoverLetter string `json:"coverLetter,omitempty"`
	// DontReallyApply toggles actually applying or just trying out the API. Set
	// this to 'true' to test without applying.
	// API default is `false`.
	DontReallyApply *bool `json:"dontReallyApply"`
	// Email is a valid email we can reach you at.
	// This is *required* when submitting an application.
	Email string `json:"email"`
	// GitHub is a URL for GitHub, Bitbucket, Gitlab or anywhere else your code is hosted!
	GitHub string `json:"github,omitempty"`
	// Job is the job you're looking to apply for.
	// This is *required* when submitting an application.
	Job string `json:"job"`
	// LinkedIn is a link to a LinkedIn profile.
	LinkedIn string `json:"linkedin,omitempty"`
	// Name is your full name.
	// This is *required* when submitting an application.
	Name string `json:"name"`
	// Pronouns is a list pronouns an applicant uses.
	Pronouns string `json:"pronouns,omitempty"`
}

// ApplyResponse represents the API response when an application is submitted.
type ApplyResponse struct {
	Careers   string   `json:"careers"`
	Keyvalues string   `json:"keyvalues"`
	Message   string   `json:"message"`
	Poem      []string `json:"poem"`
	Questions string   `json:"questions?"`
}

// Get a list of open roles at ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/getopenroles
func (c ApplyClient) Get() ([]OpenRole, *APIResponse, error) {
	response := []OpenRole{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "GET",
		Endpoint:     ApplyEndpoint,
		UseAuth:      true,
		OkStatusCode: []int{200},
		Response:     &response,
	})

	return response, apiResponse, err
}

// Apply for an open role at ReadMe.
//
// API Reference: https://docs.readme.com/main/reference/applytoreadme
func (c ApplyClient) Apply(application Application) (ApplyResponse, *APIResponse, error) {
	payload, err := json.Marshal(application)
	if err != nil {
		return ApplyResponse{}, &APIResponse{}, fmt.Errorf("unable to parse application: %w", err)
	}

	response := ApplyResponse{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "POST",
		Endpoint:     ApplyEndpoint,
		UseAuth:      true,
		Payload:      payload,
		Headers:      []RequestHeader{{"Content-Type": "application/json"}},
		OkStatusCode: []int{200},
		Response:     &response,
	})

	return response, apiResponse, err
}
