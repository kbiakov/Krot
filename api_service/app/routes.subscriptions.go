package main

import (
	"github.com/labstack/echo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"sync"
	"log"

	pb "../rpc"
)

var client *pb.SubscriptionServiceClient

func getRpcClientInstance() pb.SubscriptionServiceClient {
	sync.Once.Do(func() {
		client = &newRpcClient()
	})

	return client
}

func newRpcClient() *pb.SubscriptionServiceClient {
	// set up a connection to the server
	conn, err := grpc.Dial("localhost:9020", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	return &pb.NewSubscriptionServiceClient(conn)
}

func GetSubscriptions(ctx echo.Context) error {
	var ss []pb.Subscription
	uid := ctx.Param("uid")
	query := bson.M{"user_id": uid}

	if err := mongo.C("subscriptions").Find(query).All(&ss); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, ss)
}

func CreateSubscription(ctx echo.Context) error {
	s := new(pb.Subscription)
	if err := ctx.Bind(s); err != nil {
		return err
	}

	c := getRpcClientInstance()
	res, err := c.Subscribe(context.Background(), s)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func performForId(ctx echo.Context, statusOk int, handler func(pb.SubscriptionServiceClient, *pb.SubscriptionId)) error {
	sID := new(pb.SubscriptionId)
	if err := ctx.Bind(sID); err != nil {
		return err
	}

	c := getRpcClientInstance()
	res, err := handler(&c, sID)
	if err != nil {
		return err
	}

	return ctx.JSON(statusOk, res)
}

func StopSubscription(ctx echo.Context) error {
	return performForId(&ctx, http.StatusAccepted,
		func(c pb.SubscriptionServiceClient, sID *pb.SubscriptionId) {
			return c.StopSubscription(context.Background(), sID)
		})
}

func ResumeSubscription(ctx echo.Context) error {
	return performForId(&ctx, http.StatusAccepted,
		func(c pb.SubscriptionServiceClient, sID *pb.SubscriptionId) {
			return c.ResumeSubscription(context.Background(), sID)
		})
}

func RemoveSubscription(ctx echo.Context) error {
	return performForId(&ctx, http.StatusAccepted,
		func(c pb.SubscriptionServiceClient, sID *pb.SubscriptionId) {
			return c.Unsubscribe(context.Background(), sID)
		})
}