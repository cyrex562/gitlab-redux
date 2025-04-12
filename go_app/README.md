# GitLab Go Services

This directory contains the Go implementation of GitLab's backend and supporting services.

## Project Structure

```
go/
├── api/        # API definitions, OpenAPI/Swagger specs, JSON schema files, protocol definition files
├── cmd/        # Main applications for this project
├── configs/    # Configuration file templates or default configs
├── internal/   # Private application and library code
├── pkg/        # Library code that's ok to use by external applications
└── scripts/    # Scripts to perform various build, install, analysis, etc operations
```

## Development

### Prerequisites

- Go 1.23 or later
- Make (for using Makefile commands)

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
