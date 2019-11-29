package main

import (
	"log"

	"github.com/quintans/go-clean-ddd/internal/framework"
)

func main() {
	err := framework.NewService(":8080")
	log.Fatal(err)
}
