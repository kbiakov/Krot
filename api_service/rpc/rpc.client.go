package rpc

import (
	"github.com/labstack/echo"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"sync"
)

const addr = "localhost:9020"

var once sync.Once

var client SubscriptionServiceClient

func GetRpcClientInstance() SubscriptionServiceClient {
	once.Do(func() {
		client = newRpcClient()
	})

	return client
}

func newRpcClient() SubscriptionServiceClient {
	// set up a connection to the server
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	return NewSubscriptionServiceClient(conn)
}

func PerformForId(ctx echo.Context, handler func(c SubscriptionServiceClient, sID *SubscriptionId) (*Response, error)) error {
	sID := new(SubscriptionId)
	if err := ctx.Bind(sID); err != nil {
		return err
	}

	c := GetRpcClientInstance()
	res, err := handler(c, sID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusAccepted, res)
}
