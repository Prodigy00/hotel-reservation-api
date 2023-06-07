package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/db"
	"github.com/prodigy00/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http/httptest"
	"testing"
)

const (
	testdburi = "mongodb://localhost:27017"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, db.TestDbName),
	}
}

func TestUserHandler_HandleCreateUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandleCreateUser)

	params := types.CreateUserParams{
		Email:     "test@foo.com",
		FirstName: "Figaro",
		LastName:  "Pizzaria",
		Password:  "testpwd",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(res.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting encypted password not to be returned")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}

//func TestUserHandler_HandleGetUser(t *testing.T) {
//	tdb := setup(t)
//	defer tdb.teardown(t)
//
//	app := fiber.New()
//	userHandler := NewUserHandler(tdb.UserStore)
//
//	app.Get("/", userHandler.HandleGetUser)
//
//}
