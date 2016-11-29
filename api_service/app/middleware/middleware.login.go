package krot

import "github.com/kataras/iris"

type LoginContext struct {
	iris.Context

	User, ExistingUser *User
}

func validateUserBody(ctx *LoginContext) {
	user := &User{}
	if err := ctx.ReadJSON(user); err != nil {
		ctx.JSON(iris.StatusBadRequest, err)
	}

	existingUser, err := GetUser(user.Email)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, err)
	}

	ctx.User = user
	ctx.ExistingUser = existingUser
	ctx.Next()
}
