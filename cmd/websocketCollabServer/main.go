package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/masterkusok/websocketCollab/internal/routes"
)

func main() {
	app := fiber.New()
	routes.CreateRouting(app)
	log.Fatal(app.Listen(":8080"))
}
