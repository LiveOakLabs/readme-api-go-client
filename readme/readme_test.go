package readme_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

// mockPaginatedRequestHeader represents common and valid HTTP headers for paginated requests.
var mockPaginatedRequestHeader = http.Header{
	"Link":          {`<>; rel="next", <>; rel="prev", <>; rel="last"`},
	"X-Total-Count": {"20"},
}

// TestNewClient tests the package's main setup function.
func Test_NewClient(t *testing.T) {
	t.Run("when called with default URL and valid token", func(t *testing.T) {
		// Act
		client, err := readme.NewClient("atoken")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, "atoken", client.Token, "it is configured with expected API token")
		assert.Equal(t, readme.ReadmeAPIURL, client.APIURL, "it is configured with default API URL")
	})

	t.Run("when a custom API URL is specified", func(t *testing.T) {
		// Act
		client, err := readme.NewClient("atoken", "http://readme-test.local/v2")

		// Assert
		assert.NoError(t, err, "returns no errors when single custom API URL is provided")
		assert.Equal(t, client.APIURL, "http://readme-test.local/v2",
			"returns expected client configured API URL when custom API URL is provided")
		assert.Equal(t, client.Token, "atoken")
	})

	t.Run("when too many parameters are specified", func(t *testing.T) {
		// Act
		_, err := readme.NewClient("atoken", "one", "toomany")

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.ErrorContains(t, err, "too many values specified for API URL", "it returns the expected error")
	})
}

// TestHasNextPage tests parsing a 'link' header with a 'next' link for pagination.
//
// NOTE: `readme.HasNextPage()` refers to a private `hasNextPage()` function that has been exported for tests.
func Test_HasNextPage(t *testing.T) {
	testValid := []struct {
		expect bool
		value  string
	}{
		{true, `</api-specification?page=2>; rel="next", <>; rel="prev", <>; rel="last"`},
		{false, `<>; rel="next", <>; rel="prev", <>; rel="last"`},
	}
	for _, tc := range testValid {
		testName := fmt.Sprintf("when hasNextPage == %v", tc.expect)
		t.Run(testName, func(t *testing.T) {
			// Arrange
			linksHeader := tc.value

			// Act
			got, err := readme.HasNextPage(linksHeader)

			// Assert
			assert.NoError(t, err, "it does not return an error")
			assert.Equal(t, tc.expect, got, fmt.Sprintf("it returns %v", got))
		})
	}

	testInvalid := []struct {
		value string
	}{
		{"invalid"},
		{`<> rel="next", <>; rel="prev", <>; rel="last"`},
	}
	for _, tc := range testInvalid {
		t.Run("when link header is invalid", func(t *testing.T) {
			// Arrange
			linksHeader := tc.value
			// Act
			expect, err := readme.HasNextPage(linksHeader)
			// Assert
			assert.Error(t, err, "it returns an error")
			assert.ErrorContains(t, err, "unable to parse link header", "it returns the expected error")
			assert.False(t, expect, "it returns false")
		})
	}
}
