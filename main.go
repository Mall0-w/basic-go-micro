package main

import (
	"github.com/Mall0-w/basic-go-micro/gateway"
)

func main() {
	router := gateway.CreateRouter()
	router.Run("localhost:8080")
}
