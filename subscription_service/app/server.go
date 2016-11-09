package main

import (
	"golang.org/x/net/context"
	"github.com/nsqio/go-nsq"
	"github.com/bamzi/jobrunner"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"

	"net"
	"fmt"
	"log"

	pb "../rpc"
)

type RpcServer struct {}

func (rpc *RpcServer) Subscribe(ctx context.Context, s *pb.Subscription) (*pb.Response, error) {
	err := (&Subscription{
		UserId: s.UserId,
		Type: uint8(s.Type),
		Url: s.Url,
		Tag: s.Tag,
		PollMs: s.PollMs,
		Status: uint8(s.Status),
	}).Subscribe()

	return &pb.Response{
		Success:err == nil,
		Error:err.Error(),
	}, nil
}

func (rpc *RpcServer) ResumeSubscription(ctx context.Context, sId *pb.SubscriptionId) (*pb.Response, error) {
	return performForId(sId.Id, func(s *Subscription) error {
		return s.ResumeSubscription()
	})
}

func (rpc *RpcServer) StopSubscription(ctx context.Context, sId *pb.SubscriptionId) (*pb.Response, error) {
	return performForId(sId.Id, func(s *Subscription) error {
		return s.StopSubscription()
	})
}

func (rpc *RpcServer) Unsubscribe(ctx context.Context, sId *pb.SubscriptionId) (*pb.Response, error) {
	return performForId(sId.Id, func(s *Subscription) error {
		return s.Unsubscribe()
	})
}

func performForId(id string, handler func (*Subscription) error) (*pb.Response, error) {
	s, err := GetSubscription(id)
	if err != nil {
		return &pb.Response{
			Success: false,
			Error: err.Error(),
		}, nil
	}

	err = handler(s)

	return &pb.Response{
		Success: err == nil,
		Error: err.Error(),
	}, nil
}

var (
	mongo *mgo.Database

	w *nsq.Producer
)

func main() {
	// Connect to MongoDB server
	session, err := mgo.Dial("mymongo:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Hold database context
	mongo = session.DB("krot")

	// Start job runner
	jobrunner.Start()
	defer jobrunner.Stop()

	// Create NSQ producer
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer("lookupd:4150", config)
	defer w.Stop()

	// Start RPC server
	startRpcServer(9020)
}

func startRpcServer(port int)  {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSubscriptionServiceServer(s, &RpcServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
