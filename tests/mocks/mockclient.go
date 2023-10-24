// Package mocks provides a mock implementation of the readme.Client and its
// interfaces.
package mocks

import (
	"testing"

	readme "github.com/liveoaklabs/readme-api-go-client/readme"
)

// MockClient is a mock implementation of the readme.Client interfaces.
// This provides a mock implementation of each service in the API.
type MockClient struct {
	APIRegistry      *MockAPIRegistryService
	APISpecification *MockAPISpecificationService
	Apply            *MockApplyService
	Category         *MockCategoryService
	Changelog        *MockChangelogService
	CustomPage       *MockCustomPageService
	Doc              *MockDocService
	Image            *MockImageService
	Project          *MockProjectService
	Version          *MockVersionService
}

// New provides a mock client for use in tests to provide a mock implementation
// of the readme.Client interface. This instantiates all of the mock services
// as a convenience, rather than requiring the caller to instantiate each
// service individually.
//
// Two variations of the client are returned: the first is the actual client
// that is used in tests, and the second is a mock client that can be used to
// assert expectations on the client.
func New(t *testing.T) (*readme.Client, *MockClient) {
	client, err := readme.NewClient("test", "http://api.example.com/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %s", err)
	}

	mockClient := &MockClient{
		APIRegistry:      NewMockAPIRegistryService(t),
		APISpecification: NewMockAPISpecificationService(t),
		Apply:            NewMockApplyService(t),
		Category:         NewMockCategoryService(t),
		Changelog:        NewMockChangelogService(t),
		CustomPage:       NewMockCustomPageService(t),
		Doc:              NewMockDocService(t),
		Image:            NewMockImageService(t),
		Project:          NewMockProjectService(t),
		Version:          NewMockVersionService(t),
	}

	client.APIRegistry = mockClient.APIRegistry
	client.APISpecification = mockClient.APISpecification
	client.Apply = mockClient.Apply
	client.Category = mockClient.Category
	client.Changelog = mockClient.Changelog
	client.CustomPage = mockClient.CustomPage
	client.Doc = mockClient.Doc
	client.Image = mockClient.Image
	client.Project = mockClient.Project
	client.Version = mockClient.Version

	return client, mockClient
}
