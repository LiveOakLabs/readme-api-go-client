package testdata

import "github.com/liveoaklabs/readme-api-go-client/readme"

var Project = readme.Project{
	Name:      "Go Testing",
	SubDomain: "foobar",
	BaseURL:   "https://developer.example.com",
	Plan:      "enterprise",
	JWTSecret: "123456789abcdef",
}
