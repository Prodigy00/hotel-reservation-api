package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/api"
	"github.com/prodigy00/hotel-reservation-api/api/middleware"
	"github.com/prodigy00/hotel-reservation-api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var errConfig = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5001", "The port of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}

	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		hotelHandler = api.NewHotelHandler(store)
		userHandler  = api.NewUserHandler(store)
		authHandler  = api.NewAuthHandler(userStore)
		roomHandler  = api.NewRoomHandler(store)
		app          = fiber.New(errConfig)
		auth         = app.Group("/api")
		v1           = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	)

	app.Get("/", handleHello)

	//auth
	auth.Post("/auth", authHandler.HandleAuth)

	//versioned api routes
	//users
	v1.Get("/users", userHandler.HandleGetUsers)
	v1.Get("/users/:id", userHandler.HandleGetUser)
	v1.Post("/users", userHandler.HandleCreateUser)
	v1.Put("/users/:id", userHandler.HandleUpdateUser)
	v1.Delete("/users/:id", userHandler.HandleDeleteUser)

	//hotels
	v1.Get("/hotels", hotelHandler.HandleGetHotels)
	v1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)
	v1.Get("/hotels/:id", hotelHandler.HandleGetHotel)

	v1.Get("/rooms", roomHandler.HandleGetRooms)
	v1.Post("/rooms/:id/book", roomHandler.HandleBookRoom)
	app.Listen(*listenAddr)
}

func handleHello(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"msg": "Hello World!",
	})
}
