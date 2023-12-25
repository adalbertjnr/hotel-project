package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/adalbertjnr/hotel-project/db"
	"github.com/adalbertjnr/hotel-project/types"
	"github.com/gofiber/fiber/v2"
)

func insertUser(t *testing.T, userStore db.UserStorer) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "james",
		LastName:  "foo",
		Email:     "james@foo.com",
		Password:  "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	newUserWithEncpw, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return newUserWithEncpw
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertUser(t, tdb.UserStorer)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStorer)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "james@foo.com",
		Password: "passwordnotcorrect",
	}

	b, _ := json.Marshal(authParams)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status 400 but got %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatal("expected type error")
	}
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUserEncpw := insertUser(t, tdb.UserStorer)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStorer)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "james@foo.com",
		Password: "password",
	}

	b, _ := json.Marshal(authParams)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Token == "" {
		t.Fatal()
	}
	//set encpw empty because the types.User do not return encpw field as json
	//when try to compare they do not match and the test fail
	//insertedUser has the encpw because of the func insertUser(t, tdb.UserStore) returns it
	//but the authResp type  AuthRespone do not return the encpw field
	//for test just insert in the !reflect...
	//fmt.Println(insertedUserEncpw)
	//fmt.Println(authResp.User)
	insertedUserEncpw.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUserEncpw, authResp.User) {
		t.Fatal("struct field values not equal")
	}
}
