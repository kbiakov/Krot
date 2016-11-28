package krot

import "github.com/kataras/iris"

func getReceivers(ctx *iris.Context) {
	receivers, err := GetReceivers(ctx.Param("uid"))
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, err)
	}

	ctx.JSON(iris.StatusOK, receivers)
}

func createReceiver(ctx *iris.Context) {
	receiver := &Receiver{}
	if err := ctx.ReadJSON(receiver); err != nil {
		ctx.JSON(iris.StatusBadRequest, err)
	}

	userID := ctx.Param("uid")
	if err := receiver.CreateReceiver(userID); err != nil {
		ctx.JSON(iris.StatusInternalServerError, err)
	}

	ctx.JSON(iris.StatusCreated, "Receiver created.")
}

func removeReceiver(ctx *iris.Context) {
	userID := ctx.Param("uid")
	name := ctx.Param("name")

	if err := RemoveReceiver(userID, name); err != nil {
		if err == ErrReceiverNotFound {
			ctx.JSON(iris.StatusNotFound, "Receiver " + name + " not found")
		}

		ctx.JSON(iris.StatusInternalServerError, err)
	}

	ctx.JSON(iris.StatusOK, "Receiver created.")
}
