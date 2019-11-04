// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/library"

	jwt "github.com/dgrijalva/jwt-go"
)

func TestToken_Compose(t *testing.T) {
	// setup types
	id := int64(1)
	name := "foo"
	token := "bar"
	hash := "baz"
	u := &library.User{
		ID:    &id,
		Name:  &name,
		Token: &token,
		Hash:  &hash,
	}

	tkn := jwt.New(jwt.SigningMethodHS256)
	claims := tkn.Claims.(jwt.MapClaims)
	claims["active"] = u.Active
	claims["admin"] = u.Admin
	claims["name"] = u.Name
	want, err := tkn.SignedString([]byte(u.GetHash()))
	if err != nil {
		t.Errorf("Unable to create test token: %v", err)
	}

	// run test
	got, err := Compose(u)
	if err != nil {
		t.Errorf("Compose returned err: %v", err)
	}

	if !strings.EqualFold(got, want) {
		t.Errorf("Compose is %v, want %v", got, want)
	}
}

func TestToken_Parse(t *testing.T) {
	// setup types
	id := int64(1)
	name := "foo"
	token := "bar"
	hash := "baz"
	zeroBool := false
	want := &library.User{
		ID:     &id,
		Name:   &name,
		Token:  &token,
		Hash:   &hash,
		Active: &zeroBool,
		Admin:  &zeroBool,
	}

	tkn, err := Compose(want)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateUser(want)

	// run test
	got, err := Parse(tkn, db)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestToken_Parse_Error_NoParse(t *testing.T) {
	// setup types
	id := int64(1)
	name := "foo"
	token := "bar"
	hash := "baz"
	u := &library.User{
		ID:    &id,
		Name:  &name,
		Token: &token,
		Hash:  &hash,
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(u)

	// run test
	got, err := Parse("!@#$%^&*()", db)
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestToken_Parse_Error_InvalidSignature(t *testing.T) {
	// setup types
	id := int64(1)
	name := "foo"
	token := "bar"
	hash := "baz"
	u := &library.User{
		ID:    &id,
		Name:  &name,
		Token: &token,
		Hash:  &hash,
	}

	tkn := jwt.New(jwt.SigningMethodHS512)
	claims := tkn.Claims.(jwt.MapClaims)
	claims["active"] = u.Active
	claims["admin"] = u.Admin
	claims["name"] = u.Name
	token, err := tkn.SignedString([]byte(u.GetHash()))
	if err != nil {
		t.Errorf("Unable to create test token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateUser(u)

	// run test
	got, err := Parse(token, db)
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestToken_Retrieve(t *testing.T) {
	// setup types
	want := "foobar"

	header := fmt.Sprintf("Bearer %s", want)
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", header)

	// run test
	got, err := Retrieve(request)
	if err != nil {
		t.Errorf("Retrieve returned err: %v", err)
	}

	if !strings.EqualFold(got, want) {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}

func TestToken_Retrieve_Error(t *testing.T) {
	// setup types
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// run test
	got, err := Retrieve(request)
	if err == nil {
		t.Errorf("Retrieve should have returned err")
	}

	if len(got) > 0 {
		t.Errorf("Retrieve is %v, want \"\"", got)
	}
}
