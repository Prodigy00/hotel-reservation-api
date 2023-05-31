package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5001", "The port of the API server")
	flag.Parse()

	app := fiber.New()
	app.Get("/", handleHello)

	v1 := app.Group("/api/v1")

	v1.Get("/user", api.HandleGetUsers)
	v1.Get("/user/:id", api.HandleGetUser)

	app.Listen(*listenAddr)
}

func handleHello(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"msg": "Hello World!",
	})
}
