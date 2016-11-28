package auth

import (
	// jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/kataras/iris"
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

func signUp(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)

	if ctx.User.Email == ctx.ExistingUser.Email {
		ctx.JSON(iris.StatusConflict, "User with this email is already exists.")
	}

	if err := ctx.User.CreateUser(); err != nil {
		ctx.JSON(iris.StatusInternalServerError, err)
	}

	//GetJwtMiddlewareInstance().mw.Get(ctx). ???
	// TODO: generate JWT token
	ctx.JSON(iris.StatusCreated, nil)
}

func signIn(c echo.Context) error {
	if !ctx.User.IsValidPassword(ctx.ExistingUser.Password) {
		ctx.JSON(iris.StatusForbidden, "Invalid email address.")
	}

	// TODO: generate JWT token
	ctx.JSON(iris.StatusOK, nil)
}

func logout(c echo.Context) error {
	// TODO: remove user receivers, change subscriptions statuses
	ctx.JSON(iris.StatusNotImplemented, nil)
}
