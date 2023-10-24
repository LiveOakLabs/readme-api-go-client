package readme_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

var apiSpecEndpointPaginated = readme.APISpecificationEndpoint + "?perPage=100&page=1"

func Test_APISpecification_GetAll(t *testing.T) {
	t.Run("when called no parameters", func(t *testing.T) {
		// Arrange
		expect := testdata.APISpecifications
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(200).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "1").
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.APISpecification.GetAll()

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns []APISpecification struct")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when called with RequestOptions parameter", func(t *testing.T) {
		// Arrange
		expect := testdata.APISpecifications
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(200).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=1>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "1").
			JSON(expect)
		defer gock.Off()

		req := readme.RequestOptions{
			Version: "1.2.3",
			Headers: []readme.RequestHeader{{"foo": "bar"}},
		}

		// Act
		got, resp, err := TestClient.APISpecification.GetAll(req)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns []APISpecification struct")
		assert.Equal(t, req.Headers, resp.Request.Headers,
			"it returns expected response headers")
		assert.Equal(t, req.Version, resp.Request.Version,
			"it returns expected response version")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when API responds with an error", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(400).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "1").
			JSON(testdata.APISpecResponseVersionEmtpy.APIErrorResponse)
		defer gock.Off()

		// Act
		_, _, err := TestClient.APISpecification.GetAll()

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "API responded with a non-OK status: 400",
			"it returns expected error")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when API response cannot be parsed", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(200).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "1").
			JSON(`[{"invalid":invalid"}]`)
		defer gock.Off()

		// Act
		_, _, err := TestClient.APISpecification.GetAll()

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "unable to parse API response: invalid character",
			"it returns expected error")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when request is successful and there are no results", func(t *testing.T) {
		// Arrange
		var expectResult []readme.APISpecification
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(200).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "0").
			JSON(`[]`)
		defer gock.Off()

		// Act
		got, _, err := TestClient.APISpecification.GetAll()

		// Assert
		assert.NoError(t, err, "it does not returns errors")
		assert.Equal(t, expectResult, got, "it returns empty APISpecification slice")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when API response cannot be parsed", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(200).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "1").
			JSON(`invalid`)
		defer gock.Off()

		// Act
		_, _, err := TestClient.APISpecification.GetAll()

		// Assert
		assert.ErrorContains(t, err, "unable to parse API response",
			"it returns expected error")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})
}

func Test_APISpecifications_GetAll_Paginated(t *testing.T) {
	t.Run("when pagination has unexpected (page >= (totalCount / perPage))", func(t *testing.T) {
		// Arrange
		expect := testdata.APISpecifications
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(200).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "1").
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.APISpecification.
			GetAll(readme.RequestOptions{PerPage: 3, Page: 1})

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns paginated items")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when page >= (totalCount / perPage)", func(t *testing.T) {
		// Arrange
		expect := testdata.APISpecifications
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(200).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "1").
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.APISpecification.GetAll(readme.RequestOptions{PerPage: 3, Page: 15})

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns the results")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when pagination header cannot be parsed", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(200).
			AddHeader("Link", `invalid`).
			AddHeader("x-total-count", "1").
			JSON(testdata.APISpecifications)
		defer gock.Off()

		expect := "unable to parse link header - invalid format"

		// Act
		_, _, err := TestClient.APISpecification.GetAll()

		// Assert
		assert.ErrorContains(t, err, expect, "it returns expected error")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when pagination header contains invalid count", func(t *testing.T) {
		// Arrange
		var expect []readme.APISpecification
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(200).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "x").
			JSON(expect)
		defer gock.Off()

		expectErr := "unable to parse 'x-total-count' header"

		// Act
		got, _, err := TestClient.APISpecification.GetAll(readme.RequestOptions{PerPage: 5, Page: 1})

		// Assert
		assert.ErrorContains(t, err, expectErr, "it returns expected error")
		assert.Equal(t, expect, got, "it returns empty APISpecification slice")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})
}

func Test_APISpecification_Get(t *testing.T) {
	t.Run("when called with an ID that exists", func(t *testing.T) {
		// Arrange
		expect := testdata.APISpecifications
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint). ///"+expect.ID).
			Reply(200).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "1").
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.APISpecification.Get(expect[0].ID)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect[0], got, "it returns a single APISpecification struct")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when called with an ID that does not exist", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(200).
			AddHeader("Link", `<`+apiSpecEndpointPaginated+`&page=2>; rel="next", <>; rel="prev", <>; rel="last"`).
			AddHeader("x-total-count", "1").
			JSON(testdata.APISpecifications)
		defer gock.Off()

		expect := "API specification not found"

		// Act
		_, _, err := TestClient.APISpecification.Get("doesnotexist")

		// Assert
		assert.ErrorContains(t, err, expect, "it returns expected error")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when API response returns an error", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Get(readme.APISpecificationEndpoint).
			Reply(400).
			JSON(testdata.APISpecResponseVersionEmtpy.APIErrorResponse)
		defer gock.Off()

		expect := "unable to retrieve API specifications"

		// Act
		_, _, err := TestClient.APISpecification.Get("0123456789abcdef")

		// Assert
		assert.ErrorContains(t, err, expect, "it returns expected error")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})
}

func Test_APISpecification_Create(t *testing.T) {
	t.Run("when version is not specified and API response is 200", func(t *testing.T) {
		// Arrange
		expect := readme.APISpecificationSaved{
			ID:    "0123456789",
			Title: "My Test API",
		}

		gock.New(TestClient.APIURL).
			Post(readme.APISpecificationEndpoint).
			Reply(201).
			JSON(expect)

		createParams := `{"name": "My Test API OpenAPI Spec"}`

		// Act
		got, _, err := TestClient.APISpecification.Create(createParams)

		// Assert
		assert.NoError(t, err, "it returns no errors")
		assert.Equal(t, expect, got, "it returns expected struct")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when version is specified", func(t *testing.T) {
		// Arrange
		expect := readme.APISpecificationSaved{
			ID:    "0123456789",
			Title: "My Test API",
		}

		gock.New(TestClient.APIURL).
			Post(readme.APISpecificationEndpoint).
			Reply(201).
			JSON(expect)

		requestOptions := readme.RequestOptions{
			Version: "1.2.3",
		}

		createParams := `{"name": "My Test API OpenAPI Spec"}`

		// Act
		_, gotResponse, err := TestClient.APISpecification.Create(createParams, requestOptions)

		// Assert
		assert.Equal(t, requestOptions.Version, gotResponse.Request.Version,
			"it response with the requested version")
		assert.NoError(t, err, "it does not return an error")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when registry UUID is specified", func(t *testing.T) {
		// Arrange
		expect := readme.APISpecificationSaved{
			ID:    "0123456789",
			Title: "My Test API",
		}

		gock.New(TestClient.APIURL).
			Post(readme.APISpecificationEndpoint).
			Reply(201).
			JSON(expect)

		// Act
		got, _, err := TestClient.APISpecification.Create("uuid:3bbeunznlboryu0oo")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got,
			"it returns expected APISpecificationSaved struct")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when API response cannot be parsed", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Post(readme.APISpecificationEndpoint).
			Reply(200).
			JSON(`invalid`)

		expect := "unable to parse API response"
		createParams := `{"name": "My Test API OpenAPI Spec"}`

		// Act
		_, _, err := TestClient.APISpecification.Create(createParams)

		// Assert
		assert.ErrorContains(t, err, expect, "it returns expected error")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when called with empty definition and API response with 400", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Post(readme.APISpecificationEndpoint).
			Reply(400).
			JSON(testdata.APISpecResponseSpecFileEmpty.APIErrorResponse)

		expect := "SPEC_FILE_EMPTY"
		expectErr := "API responded with a non-OK status: 400"

		// Act
		_, got, err := TestClient.APISpecification.Create("")

		assert.ErrorContains(t, err, expectErr, "it returns expected error")
		assert.Equal(t, expect, got.APIErrorResponse.Error,
			"it returns expected APIErrorResponse struct")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})
}

func Test_APISpecification_Update(t *testing.T) {
	t.Run("when called with an existing ID and definition JSON", func(t *testing.T) {
		// Arrange
		expect := readme.APISpecificationSaved{
			ID:    "0123456789",
			Title: "My Test API",
		}
		gock.New(TestClient.APIURL).
			Put(readme.APISpecificationEndpoint + "/0123456789").
			Reply(201).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.APISpecification.Update(expect.ID,
			`{"name": "My Test API OpenAPI Spec"}`)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got,
			"it returns expected APISpecificationSaved struct")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})
}

func Test_APISpecification_Delete(t *testing.T) {
	t.Run("when called with existing ID and API response with success", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.APISpecificationEndpoint + "/0123456789").
			Reply(204)
		defer gock.Off()

		// Act
		got, _, err := TestClient.APISpecification.Delete("0123456789")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, true, got, "it returns true")
		assert.True(t, gock.IsDone(), "it asserts that all mocks were called")
	})

	t.Run("when called with invalid ID and API response with 400", func(t *testing.T) {
		// Arrange
		expect := testdata.APISpecResponseVersionEmtpy.APIErrorResponse
		gock.New(TestClient.APIURL).
			Delete(readme.APISpecificationEndpoint + "/0123456789").
			Reply(400).
			JSON(expect)
		defer gock.Off()

		expectErr := "API responded with a non-OK status: 400"

		// Act
		got, gotResponse, err := TestClient.APISpecification.Delete("0123456789")

		// Assert
		assert.ErrorContains(t, err, expectErr, "it returns expected error")
		assert.Equal(t, false, got, "it returns false")
		assert.Equal(t, expect, gotResponse.APIErrorResponse,
			"it returns expected APIErrorResponse struct")
		assert.True(t, gock.IsDone(),
			"it asserts that all mocks were called")
	})
}
