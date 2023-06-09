package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/db"
	"github.com/prodigy00/hotel-reservation-api/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func makeTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "john",
		LastName:  "doe",
		Email:     "johndoe@gmail.com",
		Password:  "supersecretpwd",
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthHandler_HandleAuth(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	createdUser := makeTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "johndoe@gmail.com",
		Password: "supersecretpwd",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", res.StatusCode)
	}

	var authres AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&authres); err != nil {
		t.Fatal(err)
	}

	if authres.Token == "" {
		t.Fatalf("expected the jwt token to be present in the auth rresponse but got %s", authres.Token)
	}

	//set encrypted password to empty. don't return it
	createdUser.EncryptedPassword = ""
	if !reflect.DeepEqual(createdUser, authres.User) {
		t.Fatalf("expected the user to be the created user")
	}
}

func TestAuthHandler_HandleAuthFail(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	makeTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "johndoe@gmail.com",
		Password: "supersecretpwdnotcorrect",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected http status of 200 but got %d", res.StatusCode)
	}

	var gR genericRes
	if err := json.NewDecoder(res.Body).Decode(&gR); err != nil {
		t.Fatal(err)
	}
	if gR.Msg != "invalid credentials" {
		t.Fatalf("expected err msg to be invalid credentials but got %s", gR.Msg)
	}
}
