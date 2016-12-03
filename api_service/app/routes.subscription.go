package main

import (
	"github.com/labstack/echo"
	"net/http"
	"context"

	pb "../rpc"
)

func GetSubscriptions(ctx echo.Context) error {
	ss, err := GetSubscriptionsForUserID(ctx.Param("uid"))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, ss)
}

func CreateSubscription(ctx echo.Context) error {
	s := new(pb.Subscription)
	if err := ctx.Bind(s); err != nil {
		return err
	}

	c := pb.GetRpcClientInstance()
	res, err := c.Subscribe(context.Background(), s)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func StopSubscription(ctx echo.Context) error {
	return pb.PerformForId(ctx, func(c pb.SubscriptionServiceClient, sID *pb.SubscriptionId) (*pb.Response, error) {
		return c.StopSubscription(context.Background(), sID)
	})
}

func ResumeSubscription(ctx echo.Context) error {
	return pb.PerformForId(ctx, func(c pb.SubscriptionServiceClient, sID *pb.SubscriptionId) (*pb.Response, error) {
		return c.ResumeSubscription(context.Background(), sID)
	})
}

func RemoveSubscription(ctx echo.Context) error {
	return pb.PerformForId(ctx, func(c pb.SubscriptionServiceClient, sID *pb.SubscriptionId) (*pb.Response, error) {
		return c.Unsubscribe(context.Background(), sID)
	})
}
