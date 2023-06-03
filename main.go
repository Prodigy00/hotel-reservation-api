package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/api"
	"github.com/prodigy00/hotel-reservation-api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	dburi    = "mongodb://localhost:27017"
	dbname   = "hotel-reservation"
	userColl = "users"
)

var errConfig = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5001", "The port of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(errConfig)
	app.Get("/", handleHello)

	v1 := app.Group("/api/v1")

	v1.Get("/users", userHandler.HandleGetUsers)
	v1.Get("/users/:id", userHandler.HandleGetUser)
	v1.Post("/users", userHandler.HandleCreateUser)

	app.Listen(*listenAddr)
}

func handleHello(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"msg": "Hello World!",
	})
}
