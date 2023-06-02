package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/api"
	"github.com/prodigy00/hotel-reservation-api/db"
	"github.com/prodigy00/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	dburi    = "mongodb://localhost:27017"
	dbname   = "hotel-reservation"
	userColl = "users"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	coll := client.Database(dbname).Collection(userColl)

	user := types.User{
		FirstName: "Hagrid",
		LastName:  "Hogwarts",
	}

	res, err := coll.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("created user: ", res)

	var harald types.User
	if err := coll.FindOne(ctx, bson.M{}).Decode(&harald); err != nil {
		log.Fatal(err)
	}

	fmt.Println(harald)

	listenAddr := flag.String("listenAddr", ":5001", "The port of the API server")
	flag.Parse()

	app := fiber.New()
	app.Get("/", handleHello)

	v1 := app.Group("/api/v1")

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	v1.Get("/user", api.HandleGetUsers)
	v1.Get("/user/:id", api.HandleGetUser)

	app.Listen(*listenAddr)
}

func handleHello(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"msg": "Hello World!",
	})
}
