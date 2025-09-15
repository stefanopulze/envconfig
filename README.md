# envconfig

Simple golang library to read envs into struct

## How it's work

Create a struct that store your envs

```go
package main

import "envconfig"

type Config struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT" env-default:"8080"`
	Database DB     `env-prefix:"DB"`
}

type DB struct {
	Url      string `env:"URL"`
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

The library tries to load: 
```env
HOST
PORT
DB_URL
DB_USER
DB_PASSWORD
```

## Configuration tag
You can annotate your struct field with annotation
- `env` indicate the name of the env
- `env-prefix` indicate the prefix of the env used when you have a struct
- `env-default` indicate the default value of the env