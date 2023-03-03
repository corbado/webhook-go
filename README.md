[![Go Reference](https://pkg.go.dev/badge/github.com/corbado/webhook-go.svg)](https://pkg.go.dev/github.com/corbado/webhook-go) ![example workflow](https://github.com/github/docs/actions/workflows/build.yml/badge.svg)


# Description

This Go webhooks library can be used in your backend to handle webhook calls from Corbado.

The library handles webhooks authentication and takes care of the correct formatting of requests and responses. This allows you to focus on the actual implementation of the different actions.

# Installation

For installation, use `go get`:

```
go get -u github.com/corbado/webhook-go
```

# Documentation
To learn how Corbado uses webhooks, please have a look at our [webhooks documentation](https://docs.corbado.com/helpful-guides/webhooks).

# Examples

See [examples](examples/) for a very simple usage of the webhooks library. We provide a [standard HTTP library](examples/standardlib/main.go) example and a [Gin Web Framework](examples/gin/main.go) example.

# Development

This project uses a Makefile where all tasks are configured. `make help` will print out all commands and their function. Some tasks will not work on windows!

### Linting
- Use `make lint-install` to install golangcli
- Use `make lint` to run the linter

### Testing
- Use `make unittest` to run all unittests