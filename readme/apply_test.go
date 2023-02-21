package readme_test

import (
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/internal/testutil"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

const applyTestEndpoint = "http://readme-test.local/api/v1/apply"

func Test_Apply_Get(t *testing.T) {
	// Arrange
	var expect []readme.OpenRole

	mockResponse := testutil.APITestResponse{
		URL:    applyTestEndpoint,
		Status: 200,
		Body: `
			[
				{
				"slug": "front-end-engineer",
				"title": "Front End Engineer",
				"description": "foo bar",
				"pullquote": "Collaborative: Your code reviews and pull requests are detailed, you’re always happy to share knowledge, and you love a good pairing session.",
				"location": "Remote, US",
				"department": "Engineering",
				"url": "https://jobs.ashbyhq.com/readme/43ebc7c3-4653-4037-841e-075ad428a68d/application"
				}
			]
		`,
	}
	testutil.JsonToStruct(t, mockResponse.Body, &expect)
	api := mockResponse.New(t)

	// Act
	got, _, err := api.Apply.Get()

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns open roles")
}

func Test_Apply_Apply(t *testing.T) {
	t.Run("when called with valid params", func(t *testing.T) {
		// Arrange
		expect := readme.ApplyResponse{}

		mockResponse := testutil.APITestResponse{
			URL:    applyTestEndpoint,
			Status: 200,
			Body: `
			{
				"message": "Thanks for applying, Gordon! We'll reach out to you soon!",
				"keyvalues": "https://www.keyvalues.com/readme",
				"careers": "https://readme.com/careers",
				"questions?": "greg@readme.io",
				"poem": [
				"Thanks for applying to work at ReadMe!",
				"Your application is lookin' spiffy",
				"We're going to review it ASAP",
				"And we'll get back to you in a jiffy!"
				]
			}
		`,
		}
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		got, _, err := api.Apply.Apply(readme.Application{
			Name:  "Gordon Ramsay",
			Email: "gordon@example.com",
			Job:   "Front End Engineer",
		})

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns the expected API response message")
	})

	t.Run("when called with invalid params", func(t *testing.T) {
		// Arrange
		expect := readme.APIErrorResponse{}

		mockResponse := testutil.APITestResponse{
			URL:    applyTestEndpoint,
			Status: 400,
			Body: `
			{
				"error": "APPLY_INVALID_NAME",
				"message": "You need to provide a name.",
				"suggestion": "To apply for a job, you need to include your name as a body parameter.",
				"docs": "https://docs.readme.com/main/logs/ddfb22d8-dfc6-43e8-b15f-4f4d9092f27d",
				"help": "https://docs.readme.com/main/reference and include the following link to your API log: 'https://docs.readme.com/main/logs/ddfb22d8-dfc6-43e8-b15f-4f4d9092f27d'.",
				"poem": [
				  "We normally support online anonymity",
				  "Except when you're trying to apply,",
				  "Because if we're gonna work together",
				  "We need to know what you go by!"
				]
			}
		`,
		}
		testutil.JsonToStruct(t, mockResponse.Body, &expect)
		api := mockResponse.New(t)

		// Act
		_, apiResponse, err := api.Apply.Apply(readme.Application{
			Email: "gordon@example.com",
			Job:   "Front End Engineer",
		})

		// Assert
		assert.Error(t, err, "it returns an error")
		assert.Equal(t, expect, apiResponse.APIErrorResponse, "it returns the expected API response message")
	})
}
