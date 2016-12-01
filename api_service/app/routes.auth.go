package main

import (
	"github.com/labstack/echo"
	"net/http"
	"errors"
)

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

	res, err := newAuthResponse(u)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func Login(ctx echo.Context) error {
	req := new(AuthRequest)
	if err := ctx.Bind(req); err != nil {
		return err
	}

	u, err := GetUser(req.Email)
	if err != nil {
		return err
	}
	if !u.IsValidPassword(req.Password) {
		return echo.ErrUnauthorized
	}

	res, err := newAuthResponse(u)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func Logout(ctx echo.Context) error {
	// TODO: stop all subscriptions
	return ctx.JSON(http.StatusNotImplemented, nil)
}
