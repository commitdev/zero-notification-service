# zero-notification-service

Provide email, push notification, SMS capabilities for applications.

## Overview

This project is schema-first, and is built using [openapi-generator](https://github.com/OpenAPITools/openapi-generator). It requires at least openapi-generator CLI version 5.0.0.

To add or change an API route definition, make a change to the schema in `api/` and then run `make generate`

### Running the server

To run the server:
```
make run
```

To run the server in a docker container
```
make docker
make docker-run
```

### Releasing a new version on GitHub and Brew

We are using a tool called `goreleaser` which you can get from brew if you're on MacOS:
`brew install goreleaser`

After you have the tool, you can follow these steps:
```shell
export GITHUB_TOKEN=<your token with access to write to the repo>
git tag -s -a <version number like v0.0.1> -m "Some message about this release"
git push origin <version number>
goreleaser release
```

This will create a new release in GitHub and automatically collect all the commits since the last release into a changelog.
It will also build binaries for various OSes and attach them to the release.
The configuration for goreleaser is in [.goreleaser.yml](.goreleaser.yml)

Upon a new release being published, a workflow will run to generate a version of the Helm chart.
