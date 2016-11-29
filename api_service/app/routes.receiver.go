package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func GetReceivers(ctx echo.Context) error {
	uid := ctx.Param("uid")
	rs, err := getReceivers(uid)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, rs)
}

func CreateReceiver(ctx echo.Context) error {
	r := new(Receiver)
	if err := ctx.Bind(r); err != nil {
		return err
	}

	uid := ctx.Param("uid")
	if err := r.createReceiver(uid); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, r)
}

func RemoveReceiver(ctx echo.Context) error {
	uid := ctx.Param("uid")
	name := ctx.Param("name")
	if err := removeReceiver(uid, name); err != nil {
		return err
	}

	return ctx.JSON(http.StatusAccepted, nil)
}
