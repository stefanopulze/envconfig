package envconfig

import (
	"fmt"
	"os"
	"testing"
)

type Test struct {
	Host    string `env:"HOST" env-default:"localhost"`
	Port    int    `env:"PORT" env-default:"5432"`
	Enabled bool   `env:"ENABLED" env-default:"true"`
	User    User
	Admin   User `env-prefix:"ADMIN"`
}

type User struct {
	Name string
	Age  int
}

func TestGoEnv(t *testing.T) {
	_ = os.Setenv("HOST", "127.0.0.1")
	_ = os.Setenv("PORT", "5433")
	_ = os.Setenv("ADMIN_NAME", "admin")
	_ = os.Setenv("ADMIN_AGE", "38")
	_ = os.Setenv("NAME", "mario")
	_ = os.Setenv("AGE", "18")

	cfg := new(Test)
	if err := ReadEnv(cfg); err != nil {
		t.Fatal(err)
	}

	t.Log(fmt.Sprintf("%+v", cfg))

	if cfg.Host != "127.0.0.1" {
		t.Error("Invalid host")
	}

	if cfg.Port != 5433 {
		t.Error("Invalid port")
	}

	if cfg.Enabled != true {
		t.Error("Invalid enabled")
	}

	if cfg.User.Name != "mario" {
		t.Error("Invalid user name")
	}

	if cfg.User.Age != 18 {
		t.Error("Invalid user age")
	}

	if cfg.Admin.Name != "admin" {
		t.Error("Invalid admin name")
	}

	if cfg.Admin.Age != 38 {
		t.Error("Invalid admin age")
	}

	if cfg.Admin.Name != "admin" {
		t.Error("Invalid admin name")
	}
}
