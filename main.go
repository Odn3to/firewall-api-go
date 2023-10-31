package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"firewall-api-go/router"
	"firewall-api-go/database"
)

// @title Firewall - API - Go
// @version 1.0
// @description Gerencia e lida com o serviço Nftables
// @host 172.23.58.10:8007
// @BasePath /firewall
// @schemes http https
func main() {

	database.ConectaNoBD()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	router.Register(app)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", "8007")))
}