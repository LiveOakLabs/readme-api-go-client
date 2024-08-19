package readme_test

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_CustomPages_GetAll(t *testing.T) {
	tests := []struct {
		name               string
		reqOpts            readme.RequestOptions
		mockResponseStatus int
		mockHeaders        map[string]string
		mockResponseBody   interface{}
		expectError        bool
		expectedErrorMsg   string
		expectedResult     []readme.CustomPage
		expectedEndpoint   string
	}{
		{
			name: "valid parameters with 200 response",
			reqOpts: readme.RequestOptions{
				PerPage: 100,
				Page:    1,
			},
			mockResponseStatus: 200,
			mockHeaders: map[string]string{
				"Link":          `</custompages?page=2>; rel="next", <>; rel="prev", <>; rel="last"`,
				"x-total-count": "3",
			},
			mockResponseBody: testdata.CustomPages,
			expectError:      false,
			expectedResult:   testdata.CustomPages,
		},
		{
			name: "invalid pagination header with large page number",
			reqOpts: readme.RequestOptions{
				PerPage: 6,
				Page:    16,
			},
			mockResponseStatus: 200,
			mockHeaders: map[string]string{
				"Link":          `</custompages?perPage=6&page=15>; rel="next", <>; rel="prev", <>; rel="last"`,
				"x-total-count": "90",
			},
			mockResponseBody: testdata.CustomPages,
			expectError:      false,
			expectedResult:   testdata.CustomPages,
			expectedEndpoint: "/custompages?perPage=6&page=16",
		},
		{
			name: "invalid x-total-count header",
			reqOpts: readme.RequestOptions{
				PerPage: 6,
				Page:    15,
			},
			mockResponseStatus: 200,
			mockHeaders: map[string]string{
				"Link":          `</custompages?perPage=6&page=15>; rel="next", <>; rel="prev", <>; rel="last"`,
				"x-total-count": "x",
			},
			mockResponseBody: []readme.CustomPage{},
			expectError:      true,
			expectedErrorMsg: "unable to parse 'x-total-count' header:",
			expectedResult:   nil,
		},
		{
			name: "when there are no custom pages",
			reqOpts: readme.RequestOptions{
				PerPage: 6,
				Page:    1,
			},
			mockResponseStatus: 200,
			mockHeaders: map[string]string{
				"Link":          `</custompages?page=1>; rel="next", <>; rel="prev", <>; rel="last"`,
				"x-total-count": "0",
			},
			mockResponseBody: []readme.CustomPage{},
			expectError:      false,
			expectedResult:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			resp := gock.New(TestClient.APIURL).
				Get(readme.CustomPageEndpoint).
				MatchParam("page", fmt.Sprintf("%d", tt.reqOpts.Page)).
				MatchParam("perPage", fmt.Sprintf("%d", tt.reqOpts.PerPage)).
				Reply(tt.mockResponseStatus)
			for key, value := range tt.mockHeaders {
				resp.AddHeader(key, value)
			}
			resp.JSON(tt.mockResponseBody)
			defer gock.Off()

			// Act
			got, apiResponse, err := TestClient.CustomPage.GetAll(tt.reqOpts)

			// Assert
			if tt.expectError {
				assert.Error(t, err, "it returns an error")
				if tt.expectedErrorMsg != "" {
					assert.ErrorContains(t, err, tt.expectedErrorMsg, "it returns the expected error")
				}
			} else {
				assert.NoError(t, err, "it does not return an error")
				assert.Equal(t, tt.expectedResult, got, "it returns expected []CustomPage struct")
			}

			if tt.expectedEndpoint != "" && apiResponse != nil {
				assert.Equal(t, tt.expectedEndpoint, apiResponse.Request.Endpoint, "it requests with the expected pagination query parameters")
			}

			assert.True(t, gock.IsDone(), "it makes the expected API call")
		})
	}
}

func Test_CustomPages_Get(t *testing.T) {
	// Arrange
	expect := testdata.CustomPages[0]
	gock.New(TestClient.APIURL).
		Get(readme.CustomPageEndpoint + "/" + expect.Slug).
		Reply(200).
		JSON(expect)
	defer gock.Off()

	// Act
	got, _, err := TestClient.CustomPage.Get(expect.Slug)

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected []CustomPage struct")
	assert.True(t, gock.IsDone(), "it makes the expected API call")
}

func Test_CustomPages_Create(t *testing.T) {
	// Arrange
	expect := testdata.CustomPages[0]
	gock.New(TestClient.APIURL).
		Post(readme.CustomPageEndpoint).
		Reply(201).
		JSON(expect)
	defer gock.Off()

	createParams := readme.CustomPageParams{
		Body:     expect.Body,
		Title:    expect.Title,
		HTML:     expect.HTML,
		Hidden:   &expect.Hidden,
		HTMLMode: &expect.HTMLMode,
	}

	// Act
	got, _, err := TestClient.CustomPage.Create(createParams)

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected CustomPage struct")
	assert.True(t, gock.IsDone(), "it makes the expected API call")
}

func Test_CustomPages_Update(t *testing.T) {
	// Arrange
	expect := testdata.CustomPages[0]

	gock.New(TestClient.APIURL).
		Put(readme.CustomPageEndpoint + "/" + expect.Slug).
		Reply(200).
		JSON(expect)
	defer gock.Off()

	updateParams := readme.CustomPageParams{
		Body:     expect.Body,
		Title:    expect.Title,
		HTML:     expect.HTML,
		Hidden:   &expect.Hidden,
		HTMLMode: &expect.HTMLMode,
	}

	// Act
	got, _, err := TestClient.CustomPage.Update(expect.Slug, updateParams)

	// Assert
	assert.NoError(t, err, "it does not return an error")
	assert.Equal(t, expect, got, "it returns expected CustomPage struct")
	assert.True(t, gock.IsDone(), "it makes the expected API call")
}

func Test_CustomPages_Delete(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.CustomPageEndpoint + "/foo").
			Reply(204).
			JSON("{}")
		defer gock.Off()

		// Act
		got, _, err := TestClient.CustomPage.Delete("foo")

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.True(t, got, "it returns expected CustomPage struct")
	})

	t.Run("when API responds with error", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Delete(readme.CustomPageEndpoint + "/foo").
			Reply(400).
			JSON("{}")
		defer gock.Off()

		// Act
		got, _, err := TestClient.CustomPage.Delete("foo")

		// Assert
		assert.ErrorContains(t, err, "ReadMe API Error: 400 on DELETE",
			"it return the expected error")
		assert.False(t, got, "it returns false")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}
