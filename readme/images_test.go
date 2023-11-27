package readme_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/stretchr/testify/assert"
)

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
		// Arrange
		mockImage := readme.Image{
			URL:      fmt.Sprintf("https://files.readme.io/c6f07db-%s", tc.filename),
			Filename: tc.filename,
			Width:    512,
			Height:   512,
			Color:    "#000000",
		}

		// mockImageResponse represents a mock image response body from the API
		// when uploading an image.
		mockImageResponse := []any{
			mockImage.URL,
			mockImage.Filename,
			mockImage.Width,
			mockImage.Height,
			mockImage.Color,
		}

		gock.New(readme.ImageAPIURL).
			Post("/").
			Reply(200).
			JSON(mockImageResponse)
		defer gock.Off()

		// Convert base64 image to byte array
		sampleImage, _ := base64.StdEncoding.DecodeString(tc.source)

		// Act
		got, _, err := TestClient.Image.Upload(sampleImage, tc.filename)

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, mockImage, got, "it returns the image struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	}
}

func Test_Image_Upload_Invalid(t *testing.T) {
	// Arrange
	sampleImage := []byte("not-an-image")

	// Act
	_, _, err := TestClient.Image.Upload(sampleImage, "not-an-image.png")

	// Assert
	assert.Error(t, err, "it returns an error when the image type is invalid")
}
