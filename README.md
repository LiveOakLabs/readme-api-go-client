# Go Client for the ReadMe.com API

This is a Go client library for the [ReadMe.com](https://readme.com) API.

## Getting Started

Import the package and call the [`readme.NewClient()`](docs/README.md#func-newclient)
function with a token provided to set up the API client.

```go
package main

import "github.com/lobliveoaklabs/readme-api-go-client/readme"

const readmeAPIKey string = "rdme_xxxxxxxx..."

func main() {
  client, err := readme.NewClient(readmeAPIKey)
  if err != nil {
    log.Fatal(err)
  }
}
```

## Examples

Using the [`APISpecification.GetAll()`](docs/README.md#func-apispecificationclient-getall)
method to retrieve all API specifications for a project on ReadMe.com:

```go
specs, _, err := client.APISpecifications.GetAll()
if err != nil {
    log.Fatal("Error getting API specifications: ", err)
}

if specs == nil {
    log.Fatal("No results.")
}
```

## Development

* Merge requests should merge into the `main` branch.
* Refer to the [`Makefile`](Makefile) for helpful development tasks.

This project uses a few [tools](readme/tools.go) for validating code quality and functionality:

* [pre-commit](https://pre-commit.com/) for ensuring consistency and code quality before committing (external dependency).
* [golangci-lint](https://golangci-lint.run/) for linting and formatting.
* [gofumpt](https://github.com/mvdan/gofumpt) (is included with golangci-lint).
* [gocover-cobertura](https://github.com/boumenot/gocover-cobertura) for code test coverage reporting.
* [govulncheck](https://github.com/golang/vuln) for detecting vulnerabilities in Go packages.
* [gomarkdoc](https://github.com/princjef/gomarkdoc) for generating the [Markdown docs](docs/README.md).

## References

* [Terraform provider for ReadMe](https://github.com/lobliveoaklabs/terraform-provider-readme):
  Related project using this client library.
