package main

import (
	"github.com/bamzi/jobrunner"
	"github.com/labstack/echo"
	"net/http"
)

func GetLogs(ctx echo.Context) error {
	ls, err := GetLogsForId(ctx.Param("uid"))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, ls)
}

func GetJobsStatus(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, jobrunner.StatusJson())
}

func DeleteUser(ctx echo.Context) error {
	u, err := GetUserByID(ctx.Param("uid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	if err := u.DeleteUser(); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, "User with all related data was deleted")
}
