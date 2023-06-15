# Go Client for the ReadMe.com API

[![Version](https://img.shields.io/github/v/release/liveoaklabs/readme-api-go-client)](https://github.com/liveoaklabs/readme-api-go-client/releases)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/liveoaklabs/readme-api-go-client?tab=doc)

This is a Go client library for the [ReadMe.com](https://readme.com) API.

![Live Oak Bank](.github/readme/lob-logo.png)

This package is developed by [Live Oak Bank](https://liveoakbank.com) and is not affiliated with ReadMe.com.

## Getting Started

Import the package and call the `readme.NewClient()` function with a token provided to set up the API client.

```go
package main

import "github.com/liveoaklabs/readme-api-go-client/readme"

const readmeAPIKey string = "rdme_xxxxxxxx..."

func main() {
  client, err := readme.NewClient(readmeAPIKey)
  if err != nil {
    log.Fatal(err)
  }
}
```

## Examples

Using the `APISpecification.GetAll()` method to retrieve all API specifications for a project on ReadMe.com:

```go
specs, _, err := client.APISpecifications.GetAll()
if err != nil {
    log.Fatal("Error getting API specifications: ", err)
}

if specs == nil {
    log.Fatal("No results.")
}
```

## Reference

Refer to <https://pkg.go.dev/github.com/liveoaklabs/readme-api-go-client> for the Go package documentation.

## Contributing

Refer to [`CONTRIBUTING.md`](CONTRIBUTING.md) for information on contributing to this project.

## Related

The [Terraform provider for ReadMe](https://github.com/liveoaklabs/terraform-provider-readme) uses this client library.

## License

This project is licensed under the MIT License - see the [`LICENSE`](LICENSE) file for details.
