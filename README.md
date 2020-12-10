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



