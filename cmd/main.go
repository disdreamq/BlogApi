package main

import (
	"fmt"

	"github.com/disdreamq/BlogApi/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	fmt.Printf("%p\n", cfg)
}
