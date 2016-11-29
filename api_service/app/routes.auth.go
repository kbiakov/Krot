package main

import (
	// jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"errors"
)

/*
func generateUserToken(uid string) {
	token, err := jwt.ParseFromRequest(
		req,
		func(token *jwt.Token) (interface{}, error) {
			return authBackend.PublicKey, nil
		})
}
*/

func SignUp(ctx echo.Context) error {
	u := new(User)
	if err := ctx.Bind(u); err != nil {
		return err
	}

	existingUser, err := GetUser(u.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("User with this email is already exists")
	}

	if err := u.CreateUser(); err != nil {
		return err
	}

	// TODO: generate JWT token
	// GetJwtMiddlewareInstance().mw.Get(ctx). ???
	return ctx.JSON(http.StatusCreated, u)
}

func Login(ctx echo.Context) error {
	// TODO: check uid, email & password
	// TODO: generate JWT token
	return ctx.JSON(http.StatusOK, nil)
}

func Logout(ctx echo.Context) error {
	// TODO: remove receivers, change subscription statuses
	return ctx.JSON(http.StatusNotImplemented, nil)
}
