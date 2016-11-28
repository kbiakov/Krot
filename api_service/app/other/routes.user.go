package krot

import (
	"github.com/kataras/iris"
	"github.com/bamzi/jobrunner"
)

func deleteUser(ctx *iris.Context) {
	// TODO: do not read user if no related job will be stopped
	user, err := GetUserByID(ctx.Param("uid"))
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, err)
	}

	if err = user.DeleteUser(); err != nil {
		ctx.JSON(iris.StatusBadRequest, err)
	}

	ctx.JSON(iris.StatusOK, "User with all related data was deleted.")
}

func getLogs(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, "") // TODO
}

func getJobsStatus(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, jobrunner.StatusJson())
}
