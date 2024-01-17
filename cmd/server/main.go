package main

import (
	"github.com/joho/godotenv"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/api/routes"
)

func main() {
	godotenv.Load()
	routes.Setup()
}
