# Changelog

## v0.5.0 - 2024-08-22

### Bug Fixes

* fix: nested category docs (BREAKING CHANGE) @joshbeard (#88)
  * Previously, category docs that were nested under a parent document were not
    returned in the response. This has been fixed.
  * BREAKING CHANGE: The response structure for the `category.GetDocs()` method
    has changed. The `CategoryDocsChildren` struct has been removed and the
    existing `CategoryDocs` struct is now used for the nested category docs.
    The existing `Children` field is present for nested docs, but ReadMe limits
    the depth of the nested docs to 1 level. (Category -> Parent Doc -> Child)
  
* fix: pagination results are incomplete
  * #84 broke pagination because it didn't properly append the results to the
    response. This has been fixed.
  

### Maintenance

- build(deps): bump mvdan.cc/gofumpt from 0.6.0 to 0.7.0 @dependabot (#85)
- category and docs test updates @joshbeard (#88)

## v0.4.0 - 2024-08-19

### Changes

- build(deps): bump github.com/vektra/mockery/v2 from 2.43.2 to 2.44.1 @dependabot (#83)
- build(deps): bump golang.org/x/vuln from 1.1.2 to 1.1.3 @dependabot (#81)
- Pagination: more consistent, add to changelogs (BREAKING CHANGE) @joshbeard (#84)

## v0.3.0 - 2024-08-02

### Changes

- Improve API error responses @joshbeard (#82)

## v0.2.3 - 2024-06-26

This is a maintenance release to update dependencies.

### Changes

- build(deps): bump securego/gosec from 2.19.0 to 2.20.0 @dependabot (#73)
- build(deps): bump github.com/golangci/golangci-lint from 1.57.2 to 1.59.1 @dependabot (#79)
- build(deps): bump golang.org/x/vuln from 1.0.4 to 1.1.2 @dependabot (#80)
- build(deps): bump github.com/vektra/mockery/v2 from 2.42.1 to 2.43.2 @dependabot (#78)
- build(deps): bump golang.org/x/net from 0.22.0 to 0.23.0 @dependabot (#67)
- build(deps): bump github.com/golangci/golangci-lint from 1.57.1 to 1.57.2 @dependabot (#65)
- build(deps): bump github.com/golangci/golangci-lint from 1.56.2 to 1.57.1 @dependabot (#64)
- build(deps): bump github.com/vektra/mockery/v2 from 2.42.0 to 2.42.1 @dependabot (#63)
- build(deps): bump google.golang.org/protobuf from 1.31.0 to 1.33.0 @dependabot (#62)
- build(deps): bump github.com/vektra/mockery/v2 from 2.38.0 to 2.42.0 @dependabot (#59)
- build(deps): bump github.com/stretchr/testify from 1.8.4 to 1.9.0 @dependabot (#61)
- build(deps): bump github.com/golangci/golangci-lint from 1.56.1 to 1.56.2 @dependabot (#60)

## v0.2.2 - 2024-02-15

### Changes

- Remove stray info log line @joshbeard (#58)

## v0.2.1 - 2024-02-14

### Bug Fixes

- changelog type is optional @joshbeard (#57)

### Maintenance

- build(deps): bump github.com/golangci/golangci-lint from 1.55.2 to 1.56.1 @dependabot (#55)
- build(deps): bump securego/gosec from 2.18.2 to 2.19.0 @dependabot (#56)
- build(deps): bump golang.org/x/vuln from 1.0.3 to 1.0.4 @dependabot (#54)
- build(deps): bump golang.org/x/vuln from 1.0.1 to 1.0.3 @dependabot (#51)
- build(deps): bump release-drafter/release-drafter from 5 to 6 @dependabot (#53)
- build(deps): bump mvdan.cc/gofumpt from 0.5.0 to 0.6.0 @dependabot (#52)
- build(deps): bump github.com/segmentio/golines from 0.11.0 to 0.12.2 @dependabot (#49)

## v0.2.0 - 2024-01-19

### Features

- feature: add outbound-ips implementation @joshbeard (#48)

### Maintenance

- build(deps): bump github.com/cloudflare/circl from 1.3.3 to 1.3.7 @dependabot (#47)
- build(deps): bump github.com/go-git/go-git/v5 from 5.7.0 to 5.11.0 @dependabot (#46)
- build(deps): bump golang.org/x/crypto from 0.14.0 to 0.17.0 @dependabot (#45)
- build(deps): bump github/codeql-action from 2 to 3 @dependabot (#44)
- build(deps): bump actions/setup-go from 4 to 5 @dependabot (#43)
- build(deps): bump github.com/vektra/mockery/v2 from 2.37.1 to 2.38.0 @dependabot (#41)
- feature: refactor unit tests; add mocks, test data @joshbeard (#40)

## v0.1.3 - 2023-11-08

### Bug Fixes

- fix: request options and versions @joshbeard (#39)

### Maintenance

- build(deps): bump github.com/golangci/golangci-lint from 1.54.2 to 1.55.2 @dependabot (#38)
- build(deps): bump stefanzweifel/git-auto-commit-action from 4 to 5 @dependabot (#31)
- build(deps): bump securego/gosec from 2.17.0 to 2.18.2 @dependabot (#36)
- build(deps): bump golang.org/x/net from 0.14.0 to 0.17.0 @dependabot (#33)
- build(deps): bump actions/checkout from 3 to 4 @dependabot (#30)
- build(deps): bump github.com/golangci/golangci-lint from 1.54.1 to 1.54.2 @dependabot (#28)
- build(deps): bump golang.org/x/vuln from 1.0.0 to 1.0.1 @dependabot (#29)
- build(deps): bump securego/gosec from 2.16.0 to 2.17.0 @dependabot (#27)
- build(deps): bump github.com/golangci/golangci-lint from 1.53.3 to 1.54.1 @dependabot (#26)
- build(deps): bump golang.org/x/vuln from 0.2.0 to 1.0.0 @dependabot (#25)
- build(deps): bump golang.org/x/vuln from 0.1.0 to 0.2.0 @dependabot (#24)
- build(deps): bump github.com/princjef/gomarkdoc from 1.0.0 to 1.1.0 @dependabot (#23)
- maint: go dependency updates @joshbeard (#21)
- docs: Remove useless badges; add logo @joshbeard (#22)
- build(deps): bump github.com/princjef/gomarkdoc from 0.4.1 to 1.0.0 @dependabot (#20)
- build(deps): bump github.com/golangci/golangci-lint from 1.52.2 to 1.53.2 @dependabot (#19)
- build(deps): bump github.com/stretchr/testify from 1.8.2 to 1.8.3 @dependabot (#17)
- build(deps): bump securego/gosec from 2.15.0 to 2.16.0 @dependabot (#16)
- build(deps): bump github.com/cloudflare/circl from 1.3.2 to 1.3.3 @dependabot (#15)
- build(deps): bump mvdan.cc/gofumpt from 0.4.0 to 0.5.0 @dependabot (#13)
- build(deps): bump golang.org/x/vuln from 0.0.0-20230221181318-b1b4de0d2042 to 0.1.0 @dependabot (#14)
- build(deps): bump github.com/golangci/golangci-lint from 1.52.0 to 1.52.2 @dependabot (#12)
- build(deps): bump github.com/golangci/golangci-lint from 1.51.2 to 1.52.0 @dependabot (#11)
- build(deps): bump actions/setup-go from 3 to 4 @dependabot (#10)

## v0.1.2 - 2023-03-01

### 🚀 Features

- feat: image uploads @joshbeard (#9)

### 🧰 Maintenance

- ci: dependabot for github-actions @joshbeard (#7)
- build(deps): bump github.com/stretchr/testify from 1.8.1 to 1.8.2 @dependabot (#8)

## v0.1.1 - 2023-02-22

### Changes

- docs: remove obsolete links @joshbeard (#6)
- published to new GitHub organization

## v0.1.0 - 2023-02-22

### Changes

- Initial release
