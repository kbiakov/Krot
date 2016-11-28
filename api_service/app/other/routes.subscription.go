package krot

import "github.com/kataras/iris"

func getSubscriptions(ctx *iris.Context) {
	subscriptions, err := GetSubscriptions(ctx.Param("uid"))
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, err)
	}

	ctx.JSON(iris.StatusOK, subscriptions)
}

func createSubscription(ctx *iris.Context) {
	subscription := Subscription{}
	if err := ctx.ReadJSON(&subscription); err != nil {
		ctx.JSON(iris.StatusBadRequest, err)
	}

	if subscription.UserID != ctx.Param("uid") {
		ctx.JSON(iris.StatusForbidden, "Can not subscribe another user.")
	}

	if err := subscription.Subscribe(); err != nil {
		ctx.JSON(iris.StatusInternalServerError, err)
	}

	ctx.JSON(iris.StatusCreated, "Subscription created.")
}

func updateSubscription(ctx *iris.Context) {
	subscription := Subscription{}
	if err := ctx.ReadJSON(&subscription); err != nil {
		ctx.JSON(iris.StatusBadRequest, err)
	}

	if subscription.UserID != ctx.Param("uid") {
		ctx.JSON(iris.StatusForbidden, "Can not subscribe another user.")
	}

	if err := subscription.UpdateSubscription(); err != nil {
		ctx.JSON(iris.StatusInternalServerError, err)
	}

	ctx.JSON(iris.StatusAccepted, "Subscription updated.")
}

func removeSubscription(ctx *iris.Context) {
	// TODO
}

func stopSubscription(ctx *iris.Context)  {
	// TODO
}

func resumeSubscription(ctx *iris.Context)  {
	// TODO
}
