# UUID Inspector

Web service to inspect and analyze UUIDs in any format, inspired by Symfony's `uuid:inspect`.

## Features

- Parse UUIDs in multiple formats (canonical, Base58, Base32, hex, URN, braces)
- Display all UUID representations (RFC 4122, Base58, Base32, Hex, URN)
- Extract timestamps from time-based UUIDs (v1, v6, v7)
- Dark theme UI with UnoCSS
- Full test coverage (92.8%)
- **Standard Go Project Layout** (https://github.com/golang-standards/project-layout)

## Quick Start

```bash
# Run locally
make run

# Build
make build

# Test
make test

# Deploy to Fly.io
make deploy
```

## Supported Formats

**Input:**
- `550e8400-e29b-41d4-a716-446655440000` (canonical)
- `550e8400e29b41d4a716446655440000` (no dashes)
- `urn:uuid:550e8400-e29b-41d4-a716-446655440000` (URN)
- `{550e8400-e29b-41d4-a716-446655440000}` (braces)
- `BWBeN28Vb7cMEx7Ym8AUzs` (Base58)
- `KUHIIAHCTNA5JJYWIRTFKRAAAA` (Base32)

**Output:**
- Version & Variant
- All format representations
- Timestamp extraction (v1, v6, v7)

## Project Structure

```
.
├── cmd/
│   └── uuid-inspector/    # Application entrypoint
├── internal/
│   ├── handlers/          # HTTP handlers & templates
│   └── parser/            # UUID parsing logic (92.8% coverage)
├── deploy/
│   ├── Dockerfile
│   └── fly.toml
├── Makefile
└── README.md
```

## Stack

- Go 1.25
- UnoCSS (dark theme)
- Alpine Linux (Docker)
- Fly.io deployment