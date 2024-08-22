package main

import (
	"fmt"
	"go-viper-postgres/config"
	"log"
	"os"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

func main() {
	os.Setenv("PORT", "5432")

	cfg, err := config.LoadConfig[Config](".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg)
}
