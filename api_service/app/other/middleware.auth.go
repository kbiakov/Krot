package krot

import "github.com/kataras/iris"

type AuthBody struct {
	AccessToken	string `json:"access_token"`
	UserID		string `json:"user_id"`
}

func checkAuthentication(ctx *iris.Context) {
	body := &AuthBody{}
	if err := ctx.ReadJSON(body); err != nil {
		ctx.JSON(iris.StatusUnauthorized, err)
	}

	if ctx.Param("uid") != body.UserID {
		ctx.JSON(iris.StatusForbidden, "Forbidden for this user id.")
	}

	if !isValidToken(body.AccessToken, body.UserID) {
		ctx.JSON(iris.StatusUnauthorized, "Token doesn't match for this user.")
	}

	ctx.Next()
}

func adminOnly(ctx *iris.Context) {
	// TODO
}

func isValidToken(token string, userID string) bool {
	// TODO
	return true
}
