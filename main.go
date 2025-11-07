package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"log"
)

type Config struct {
	User int `env:"USER1"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current user is %d\n", cfg.User)
}
