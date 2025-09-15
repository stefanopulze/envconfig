package main

import (
	"envconfig"
	"fmt"
	"log/slog"
)

type Foo struct {
	Name    string `env:"NAME"`
	Surname string `env:"SURNAME"`
}

func main() {
	if err := envconfig.ReadDotEnv("./examples/sample.env"); err != nil {
		panic(err)
	}

	bar := new(Foo)
	if err := envconfig.ReadEnv(bar); err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("%+v", bar))
}
