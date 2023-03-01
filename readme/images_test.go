package readme_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/liveoaklabs/readme-api-go-client/internal/testutil"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

// imagesTestEndpoint is the endpoint for the images API. Note that this is different from the
// other API endpoints - this one is undocumented and not part of the official API.
const imagesTestEndpoint = "https://dash.readme.com/api/images/image-upload"

func Test_Image_Upload(t *testing.T) {
	testCases := []struct {
		source   string // base64 encoded image
		filename string
	}{
		{
			source:   "iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII",
			filename: "test-01.png",
		},
		{
			source:   "/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAMCAgMCAgMDAwMEAwMEBQgFBQQEBQoHBwYIDAoMDAsKCwsNDhIQDQ4RDgsLEBYQERMUFRUVDA8XGBYUGBIUFRT/wAALCAABAAEBAREA/8QAFAABAAAAAAAAAAAAAAAAAAAACf/EABQQAQAAAAAAAAAAAAAAAAAAAAD/2gAIAQEAAD8AKp//2Q==",
			filename: "test-01.jpg",
		},
		{
			source:   "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII",
			filename: "test-01.gif",
		},
	}

	for _, tc := range testCases {
		// mockImage represents a mock image response struct.
		var mockImage = readme.Image{
			URL:      fmt.Sprintf("https://files.readme.io/c6f07db-%s", tc.filename),
			Filename: tc.filename,
			Width:    512,
			Height:   512,
			Color:    "#000000",
		}

		// mockImageResponse represents a mock image response body from the API when uploading an image.
		var mockImageResponse = []any{
			mockImage.URL,
			mockImage.Filename,
			mockImage.Width,
			mockImage.Height,
			mockImage.Color,
		}

		// Arrange
		mockResponse := testutil.APITestResponse{
			URL:    imagesTestEndpoint,
			Status: 200,
			Body:   testutil.StructToJson(t, mockImageResponse),
		}
		api := mockResponse.New(t)

		// Convert base64 image to byte array
		sampleImage, _ := base64.StdEncoding.DecodeString(tc.source)

		// Act
		got, _, err := api.Image.Upload(sampleImage, tc.filename)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, mockImage, got, "it returns the image struct")
	}
}

func Test_Image_Upload_Invalid(t *testing.T) {
	// Arrange
	mockResponse := testutil.APITestResponse{}
	api := mockResponse.New(t)

	sampleImage := []byte("not-an-image")

	// Act
	_, _, err := api.Image.Upload(sampleImage, "not-an-image.png")

	// Assert
	assert.Error(t, err, "it returns an error when the image type is invalid")
}
