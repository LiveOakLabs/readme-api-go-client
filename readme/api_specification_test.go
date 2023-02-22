package readme_test

import (
	"net/http"
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/internal/testutil"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

const (
	// apiSpecEndpoint is the common base URL used throughout tests.
	apiSpecEndpoint = "http://readme-test.local/api/v1/api-specification"
	// apiSpecEndpointPaginated is the common base URL with used throughout tests where pagination is considered.
	apiSpecEndpointPaginated = apiSpecEndpoint + "?perPage=100&page=1"

	// mockAPISpecResponseBody is the raw JSON body of a single APISpecification.
	mockAPISpecResponseBody = `{
		"id": "0123456789",
		"title": "Readme Testing",
		"lastSynced": "2022-11-29",
		"source": "0a1b2c3d4e5f",
		"type": "guide",
		"version": "abcdef0123456789",
		"category": {
			"title": "TestCatTitle",
			"slug": "testCatSlug",
			"order": 0,
			"type": "doc",
			"id": "00aa11bb22cc33dd44ee55ff"
		}
	}`
	// mockAPISpecListResponseBody is the raw JSON body of a list of APISpecifications.
	mockAPISpecListResponseBody = "[" + mockAPISpecResponseBody + "]"

	// mockAPISpecResponseBodyError is the JSON response in the body when an API error occurs.
	mockAPISpecResponseBodyError = `
	{
		"error": "VERSION_EMPTY",
		"message": "string",
		"suggestion": "string",
		"docs": "https://docs.readme.com/logs/6883d0ee-cf79-447a-826f-a48f7d5bdf5f",
		"help": "If you need help, email support@readme.io",
		"poem": [
		  "If you're seeing this error,",
		  "Things didn't quite go the way we hoped.",
		  "When we tried to process your request,",
		  "Maybe trying again it'll work—who knows!"
		]
	}
	`
)

func Test_APISpecification_GetAll(t *testing.T) {
	// Arrange
	var expect []readme.APISpecification

	mockResponse := testutil.APITestResponse{
		URL:     apiSpecEndpointPaginated,
		Status:  200,
		Body:    mockAPISpecListResponseBody,
		Headers: mockPaginatedRequestHeader,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	t.Run("when called no parameters", func(t *testing.T) {
		// Act
		got, _, err := api.APISpecification.GetAll()

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns []APISpecification struct")
	})

	t.Run("when called with RequestOptions parameter", func(t *testing.T) {
		// Act
		expectResponse := readme.RequestOptions{
			Version: "1.2.3",
			Headers: []readme.RequestHeader{{"foo": "bar"}},
		}
		got, gotResponse, err := api.APISpecification.GetAll(expectResponse)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns []APISpecification struct")
		assert.Equal(t, expectResponse.Headers, gotResponse.Request.Headers, "it returns expected response headers")
	})

	t.Run("when API responds with an error", func(t *testing.T) {
		// Arrange
		expect := readme.APIErrorResponse{}
		expect.Error = "VERSION_EMPTY"

		mockResponse := testutil.APITestResponse{
			URL:     apiSpecEndpointPaginated,
			Status:  400,
			Body:    mockAPISpecResponseBodyError,
			Headers: mockPaginatedRequestHeader,
		}
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		_, got, err := api.APISpecification.GetAll()

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400", "it returns expected error")
		assert.Equal(t, expect, got.APIErrorResponse, "it returns API error response")
	})

	t.Run("when API response cannot be parsed", func(t *testing.T) {
		// Arrange
		mockAPISpecResponseBody := `[{"invalid":invalid"}]`

		mockResponse := testutil.APITestResponse{
			URL:     apiSpecEndpointPaginated,
			Status:  200,
			Body:    mockAPISpecResponseBody,
			Headers: mockPaginatedRequestHeader,
		}
		api := mockResponse.New(t)

		// Act
		_, _, err := api.APISpecification.GetAll()

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "unable to parse API response: invalid character", "it returns expected error")
	})

	t.Run("when request is successful and there are no results", func(t *testing.T) {
		// Arrange
		var expect []readme.APISpecification
		mockResponse := testutil.APITestResponse{
			URL:     apiSpecEndpointPaginated,
			Status:  200,
			Body:    `[]`,
			Headers: mockPaginatedRequestHeader,
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.APISpecification.GetAll()

		// Assert
		assert.NoError(t, err, "it does not returns errors")
		assert.Equal(t, expect, got, "it returns empty APISpecification slice")
	})

	t.Run("when API response cannot be parsed", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:     apiSpecEndpointPaginated,
			Status:  200,
			Body:    "invalid",
			Headers: mockPaginatedRequestHeader,
		}
		api := mockResponse.New(t)

		// Act
		_, _, err := api.APISpecification.GetAll()

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "unable to parse API response", "it returns expected error")
	})
}

func Test_APISpecifications_GetAll_Paginated(t *testing.T) {
	t.Run("when pagination has unexpected (page >= (totalCount / perPage))", func(t *testing.T) {
		// Arrange
		var expect []readme.APISpecification
		mockResponse := testutil.APITestResponse{
			URL:    "",
			Status: 200,
			Body:   `[{},{},{}]`,
			Headers: http.Header{
				"Link":          {`</api-specification?page=2>; rel="next", <>; rel="prev", <>; rel="last"`},
				"X-Total-Count": {"6"},
			},
		}
		testutil.JsonToStruct(t, `[{},{},{},{},{},{}]`, &expect)
		api := mockResponse.New(t)

		// Act
		got, _, err := api.APISpecification.GetAll(readme.RequestOptions{PerPage: 3, Page: 1})

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected APISpecification slice containing items from paginated links")
	})

	t.Run("when page >= (totalCount / perPage)", func(t *testing.T) {
		// Arrange
		var expect []readme.APISpecification
		mockInvalidHeader := http.Header{
			"Link":          {`</api-specification?page=16>; rel="next", <>; rel="prev", <>; rel="last"`},
			"X-Total-Count": {"90"},
		}
		mockResponse := testutil.APITestResponse{
			URL:     "",
			Status:  200,
			Body:    `[{},{},{},{}]`,
			Headers: mockInvalidHeader,
		}
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		got, _, err := api.APISpecification.GetAll(readme.RequestOptions{PerPage: 6, Page: 15})

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns the results")
	})

	t.Run("when pagination header cannot be parsed", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:     apiSpecEndpointPaginated,
			Status:  200,
			Body:    mockAPISpecListResponseBody,
			Headers: http.Header{"Link": []string{`rel="next" <> rel="prev", <>; rel="last"`}},
		}
		api := mockResponse.New(t)

		// Act
		_, _, err := api.APISpecification.GetAll()

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "unable to parse link header - invalid format", "it returns expected error")
	})

	t.Run("when pagination header contains invalid count", func(t *testing.T) {
		// Arrange
		var expect []readme.APISpecification
		mockInvalidHeader := http.Header{
			"Link":          []string{`</api-specification?page=2>; rel="next", <>; rel="prev", <>; rel="last"`},
			"X-Total-Count": []string{"x"},
		}
		apiSpecEndpoint := apiSpecEndpoint + "?perPage=5&page=1"

		mockResponse := testutil.APITestResponse{
			URL:     apiSpecEndpoint,
			Status:  200,
			Body:    `[]`,
			Headers: mockInvalidHeader,
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.APISpecification.GetAll(readme.RequestOptions{PerPage: 5, Page: 1})

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "unable to parse 'x-total-count' header", "it returns expected error")
		assert.Equal(t, expect, got, "it returns empty APISpecification slice")
	})
}

func Test_APISpecification_Get(t *testing.T) {
	t.Run("when called with an ID that exists", func(t *testing.T) {
		// Arrange
		expect := readme.APISpecification{}

		mockResponse := testutil.APITestResponse{
			URL:     apiSpecEndpointPaginated,
			Status:  200,
			Body:    mockAPISpecListResponseBody,
			Headers: mockPaginatedRequestHeader,
		}
		// NOTE: api.APISpecification.Get() makes a call to GetAll() behind the
		// scenes and parses the _list_ of specifications, returning a single
		// matching spec. This mocks the GetAll() response that occurs internally.
		testutil.JsonToStruct(t, mockAPISpecResponseBody, &expect)
		api := mockResponse.New(t)

		// Act
		got, _, err := api.APISpecification.Get("0123456789")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns a single APISpecification struct")
	})

	t.Run("when called with an ID that does not exist", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:     apiSpecEndpointPaginated,
			Status:  200,
			Body:    mockAPISpecListResponseBody,
			Headers: mockPaginatedRequestHeader,
		}
		api := mockResponse.New(t)

		// Act
		_, _, err := api.APISpecification.Get("doesnotexist")

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API specification not found", "it returns expected error")
	})

	t.Run("when API response returns an error", func(t *testing.T) {
		// Arrange
		expect := readme.APIErrorResponse{}
		expect.Error = "VERSION_EMPTY"

		mockResponse := testutil.APITestResponse{
			URL:     apiSpecEndpointPaginated,
			Status:  400,
			Body:    mockAPISpecResponseBodyError,
			Headers: mockPaginatedRequestHeader,
		}
		testutil.JsonToStruct(t, mockAPISpecResponseBodyError, &expect)
		api := mockResponse.New(t)

		// Act
		_, got, err := api.APISpecification.Get("0123456789abcdef")

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "unable to retrieve API specifications", "it returns expected error")
		assert.Equal(t, expect, got.APIErrorResponse, "it returns expected API error response")
	})
}

func Test_APISpecification_Create(t *testing.T) {
	// Arrange
	expect := readme.APISpecificationSaved{
		ID:    "0123456789",
		Title: "My Test API",
	}

	mockResponse := testutil.APITestResponse{
		URL:    apiSpecEndpoint,
		Status: 200,
		Body:   `{"_id": "0123456789", "title": "My Test API"}`,
	}
	api := mockResponse.New(t)

	t.Run("when version is not specified and API response is 200", func(t *testing.T) {
		// Act
		got, _, err := api.APISpecification.Create(`{"name": "My Test API OpenAPI Spec"}`)

		// Assert
		assert.NoError(t, err, "it returns no errors")
		assert.Equal(t, expect, got, "it returns expected struct")
	})

	t.Run("when version is specified", func(t *testing.T) {
		// Arrange
		requestOptions := readme.RequestOptions{
			Version: "1.2.3",
		}

		// Act
		_, gotResponse, err := api.APISpecification.Create(`{"name": "My Test API OpenAPI Spec"}`, requestOptions)

		// Assert
		assert.Equal(t, requestOptions.Version, gotResponse.Request.Version, "it response with the requested version")
		assert.NoError(t, err, "it does not return an error")
	})

	t.Run("when registry UUID is specified", func(t *testing.T) {
		// Arrange
		expect := readme.APISpecificationSaved{
			ID:    "0123456789",
			Title: "My Test API",
		}

		mockResponse := testutil.APITestResponse{
			URL:    apiSpecEndpoint,
			Status: 200,
			Body:   `{"_id": "0123456789", "title": "My Test API"}`,
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.APISpecification.Create("uuid:3bbeunznlboryu0oo")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected APISpecificationSaved struct")
	})

	t.Run("when API response cannot be parsed", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    apiSpecEndpoint,
			Status: 200,
			Body:   "",
		}
		api := mockResponse.New(t)

		// Act
		_, _, err := api.APISpecification.Create(`{"name": "My Test API OpenAPI Spec"}`)

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "unable to parse API response", "it returns expected error")
	})

	t.Run("when called with empty definition and API response with 400", func(t *testing.T) {
		// Arrange
		expect := &readme.APIErrorResponse{}

		mockResponse := testutil.APITestResponse{
			URL:    apiSpecEndpoint,
			Status: 400,
			Body:   `{"error": "SPEC_FILE_EMPTY"}`,
		}
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		_, apiErr, err := api.APISpecification.Create("")
		got := apiErr.APIErrorResponse

		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400", "it returns expected error")
		assert.Equal(t, expect.Error, got.Error, "it returns expected APIErrorResponse struct")
	})
}

func Test_APISpecification_Update(t *testing.T) {
	t.Run("when called with an existing ID and definition JSON", func(t *testing.T) {
		// Arrange
		expect := readme.APISpecificationSaved{
			ID:    "0123456789",
			Title: "My Test API",
		}

		mockResponse := testutil.APITestResponse{
			URL:    apiSpecEndpoint + "/" + expect.ID,
			Status: 201,
			Body:   `{"_id": "0123456789", "title": "My Test API"}`,
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.APISpecification.Update(expect.ID, `{"name": "My Test API OpenAPI Spec"}`)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected APISpecificationSaved struct")
	})
}

func Test_APISpecification_Delete(t *testing.T) {
	t.Run("when called with existing ID and API response with success", func(t *testing.T) {
		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    apiSpecEndpoint + "/0123456789",
			Status: 204,
			Body:   "",
		}
		api := mockResponse.New(t)

		// Act
		got, _, err := api.APISpecification.Delete("0123456789")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, true, got, "it returns true")
	})

	t.Run("when called with invalid ID and API response with 400", func(t *testing.T) {
		// Arrange
		expect := readme.APIErrorResponse{Error: "SPEC_ID_INVALID"}
		mockResponse := testutil.APITestResponse{
			URL:    apiSpecEndpoint + "/0123456789",
			Status: 400,
			Body:   `{"error": "SPEC_ID_INVALID"}`,
		}
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		got, gotResponse, err := api.APISpecification.Delete("0123456789")

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400", "it returns expected error")
		assert.Equal(t, false, got, "it returns false")
		assert.Equal(t, expect, gotResponse.APIErrorResponse, "it returns expected APIErrorResponse struct")
	})
}
