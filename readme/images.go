package readme

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// ImageAPIURL is the base URL for the images endpoint of the ReadMe API.
// This endpoint is used for uploading images to ReadMe and is not part of the documented 'v1' API.
const ImageAPIURL = "https://dash.readme.com/api/images"

// ImagesService is an interface for using the docs endpoints of the ReadMe.com API.
//
// API Reference: https://docs.readme.com/main/reference/getdoc
type ImageService interface {
	// Upload an image to ReadMe.
	Upload(source []byte, filename ...string) (Image, *APIResponse, error)
}

// ImageClient handles uploading images to ReadMe.com.
type ImageClient struct {
	client *Client
}

// Ensure the implementation satisfies the expected interfaces.
// This is a compile-time check.
// See: https://golang.org/doc/faq#guarantee_satisfies_interface
var _ ImageService = &ImageClient{}

// Image represents an image uploaded to ReadMe.
type Image struct {
	URL      string
	Filename string
	Width    int64
	Height   int64
	Color    string
}

// okContentType returns true if the content type is a valid image type.
func okContentType(contentType string) bool {
	return contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/gif"
}

// buildForm builds a multipart form for uploading an image to ReadMe.
func buildForm(filename string, source []byte) ([]byte, string, error) {
	// Create a reader for the data.
	data := strings.NewReader(string(source))

	// Create a new form.
	formData := &bytes.Buffer{}
	writer := multipart.NewWriter(formData)

	// Add the form fields.
	_ = writer.WriteField("name", "image")
	_ = writer.WriteField("filename", filename)

	// Add the file.
	part, err := writer.CreateFormFile("data", filepath.Base(filename))
	if err != nil {
		return nil, "", fmt.Errorf("unable to create request form: %w", err)
	}

	// Copy the data into the form.
	_, err = io.Copy(part, data)
	if err != nil {
		return nil, "", fmt.Errorf("unable to copy data: %w", err)
	}

	// Close the writer.
	err = writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("unable to close writer: %w", err)
	}

	// Get the content type for the form.
	contentType := writer.FormDataContentType()

	return formData.Bytes(), contentType, nil
}

// Upload an image to ReadMe.
func (c ImageClient) Upload(source []byte, filename ...string) (Image, *APIResponse, error) {
	var image Image

	// Validate the image type.
	imageType := http.DetectContentType(source)
	if !okContentType(imageType) {
		return image, nil, fmt.Errorf("invalid image type: %s", imageType)
	}

	// Determine the filename to use.
	upload_filename := "image"
	if len(filename) > 0 {
		upload_filename = filepath.Base(filename[0])
	}

	// Build the form.
	payload, contentType, err := buildForm(upload_filename, source)
	if err != nil {
		return image, nil, fmt.Errorf("unable to build form: %w", err)
	}

	// Make the request.
	var imageResponse []any
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "POST",
		URL:          fmt.Sprintf("%s/image-upload", ImageAPIURL),
		UseAuth:      true,
		Headers:      []RequestHeader{{"Content-Type": contentType}},
		Payload:      payload,
		OkStatusCode: []int{200},
		Response:     &imageResponse,
	})
	if err != nil {
		return image, apiResponse, fmt.Errorf("unable to upload image: %w", err)
	}

	// Map the response to the struct with type assertion.
	image = Image{
		URL:      imageResponse[0].(string),
		Filename: imageResponse[1].(string),
		Width:    int64(imageResponse[2].(float64)),
		Height:   int64(imageResponse[3].(float64)),
		Color:    imageResponse[4].(string),
	}

	return image, apiResponse, err
}
