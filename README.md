# envconfig

A tiny Go library to populate structs from environment variables.
It supports nested structs with prefixes, default values, slices, maps, and custom parsing.
It can also load variables from a .env file without external dependencies.

## Overview

**envconfig** reads environment variables into your Go structs using struct tags.
It aims to be minimal, explicit, and dependency‑free.

Core capabilities:

- Map environment variables to struct fields using tags
- Nested structs with configurable prefixes
- Default values when an env var is missing
- Slices and maps (with configurable separators)
- Built‑in support for time.Duration
- Custom parsing via encoding.TextUnmarshaler or a custom Setter interface
- Optional .env file loader (examples provided)

## Requirements

- Go version: as declared in go.mod: `go 1.25`
- OS: Any platform supported by Go

## Installation

This repository is a Go module.

```bash
go get github.com/stefanopulze/envconfig
```

## Usage

Create a struct that holds your configuration and annotates fields with tags.

```go
package main

import "envconfig"

type Config struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT" env-default:"8080"`
	Database DB     `env-prefix:"DB"`
}

type DB struct {
	URL      string `env:"URL"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
}

func main() {
	cfg := new(Config)
	if err := envconfig.ReadEnv(cfg); err != nil {
		panic(err)
	}
}
```

The library will look for the following variables:

```env
HOST
PORT
DB_URL
DB_USER
DB_PASSWORD
```

### Nested structs and prefixes

Use `env-prefix` on a struct field to prepend a prefix to all its inner fields.

### Default values

Use `env-default:"<value>"` to provide a fallback when the variable is not set.

### Custom parsing

- If a field type implements `encoding.TextUnmarshaler`, its `UnmarshalText` will be used.
- If a field type (or its pointer) implements the `Setter` interface below, `SetValue` will be used:

```text
type Setter interface {
    SetValue(string) error
}
```

### time.Duration

Fields of type `time.Duration` are parsed with Go's duration syntax (e.g., `1s`, `200ms`, `5m`).

## .env file support

This repo includes minimal .env parsing support.
You can load variables from a file before calling `ReadDotEnv`:

```go
package main

import "envconfig"

func main() {
	if err := envconfig.ReadDotEnv("./examples/sample.env"); err != nil {
		panic(err)
	}
}
```

See `examples/dotenv.go` and `examples/sample.env`.

Note: The .env parser is intentionally simple; it supports lines like `KEY=VALUE`, ignores empty lines and lines
starting with `#`, and strips surrounding single or double quotes from the value.

## Tags reference

- `env` — the environment variable name (defaults to the field name, uppercased)
- `env-prefix` — a prefix added to nested struct fields (e.g., `DB_`)
- `env-default` — default value if the environment variable is not set
- `env-separator` — separator for slices and maps (default is `,`)

## Contributing

Open issues and pull requests are welcome.