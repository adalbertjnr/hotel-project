package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/adalbertjnr/hotel-project/db/fixtures"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	fixtures.AddUser(&tdb.Store, "james", "foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
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
	insertedUserEncpw := fixtures.AddUser(&tdb.Store, "james", "foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "james@foo.com",
		Password: "james_foo",
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

	fmt.Println(insertedUserEncpw)
	fmt.Println(authResp.User)
	if !reflect.DeepEqual(insertedUserEncpw, authResp.User) {
		t.Fatal("struct field values not equal")
	}
}
