# UUID Inspector

Small Go web service for inspecting UUIDs across multiple representations. It accepts common UUID formats, normalizes them, shows derived encodings, and extracts timestamps from time-based UUID versions.

The project is intentionally compact: HTTP handlers are separated from parsing logic, the parser is covered by tests, and the app can run locally, in Docker, or on Fly.io.

## Features

- Parse UUIDs in canonical, compact hex, URN, braces, Base58, and Base32 formats
- Display RFC 4122, Base58, Base32, hex, and URN representations
- Detect UUID version and variant
- Extract timestamps from UUID v1, v6, and v7
- Dark web UI served by the Go application
- Parser package covered by unit tests
- Docker/Fly.io deployment files included

## Stack

- Go 1.25
- `github.com/google/uuid`
- `github.com/btcsuite/btcutil` for Base58 support
- Docker
- Fly.io

## Quick Start

```bash
git clone https://github.com/dionisvl/go-uuid-inspector.git
cd go-uuid-inspector

make run
```

Build and test:

```bash
make build
make test
```

Deploy to Fly.io:

```bash
make deploy
```

## Supported Input Formats

```text
550e8400-e29b-41d4-a716-446655440000       # canonical
550e8400e29b41d4a716446655440000           # compact hex
urn:uuid:550e8400-e29b-41d4-a716-446655440000
{550e8400-e29b-41d4-a716-446655440000}
BWBeN28Vb7cMEx7Ym8AUzs                     # Base58
KUHIIAHCTNA5JJYWIRTFKRAAAA                 # Base32
```

## Output

For a valid UUID, the service reports:

- UUID version
- RFC variant
- canonical representation
- compact hex representation
- URN representation
- Base58 representation
- Base32 representation
- timestamp, when available for time-based UUIDs

## Project Layout

```text
.
├── cmd/
│   └── uuid-inspector/       # application entrypoint
├── internal/
│   ├── handlers/             # HTTP handlers and templates
│   └── parser/               # UUID parsing and representation logic
├── deploy/
│   ├── Dockerfile
│   └── fly.toml
├── Makefile
└── README.md
```

## Development Notes

The parser is the core of the project. Keep UUID format handling in `internal/parser` and leave HTTP rendering in `internal/handlers`.

Useful checks:

```bash
go test -v -cover ./...
go test ./internal/parser -v
```

## Why This Exists

Symfony has a useful `uuid:inspect` command. This project brings the same kind of inspection workflow to a small standalone Go web service that can be opened in a browser or deployed as a tiny utility.

