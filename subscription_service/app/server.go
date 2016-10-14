package main

import (
	"github.com/bamzi/jobrunner"
	"gopkg.in/mgo.v2"
)

import (
	"flag"
	"fmt"
	"net"

	"golang.org/x/net/context"

	"github.com/nsqio/go-nsq"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "../rpc"
)

var (
	tls	 = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "data/server1.pem", "The TLS cert file")
	keyFile	 = flag.String("key_file", "data/server1.key", "The TLS key file")
)

type RpcServer struct {}

func (s *RpcServer) Subscribe(ctx context.Context, subscription *pb.Subscription) (*pb.Response, error) {
	err := Subscription(subscription).Subscribe()

	return &pb.Response{
		Success:err == nil,
		Error:err,
	}, nil
}

func (s *RpcServer) ResumeSubscription(ctx context.Context, sId *pb.SubscriptionId) (*pb.Response, error) {
	return performForId(sId.Id, func(s *Subscription) error {
		return s.ResumeSubscription()
	})
}

func (c *RpcServer) StopSubscription(ctx context.Context, sId *pb.SubscriptionId) (*pb.Response, error) {
	return performForId(sId.Id, func(s *Subscription) error {
		return s.StopSubscription()
	})
}

func (c *RpcServer) Unsubscribe(ctx context.Context, sId *pb.SubscriptionId) (*pb.Response, error) {
	return performForId(sId.Id, func(s *Subscription) error {
		return s.Unsubscribe()
	})
}

func performForId(id string, handler func (*Subscription) error) (*pb.Response, error) {
	s, err := GetSubscription(id)
	if err != nil {
		return &pb.Response{false, err}, nil
	}

	err = handler(&s)

	return &pb.Response{
		Success: err == nil,
		Error: err,
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

func startRpcServer(port *int)  {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSubscriptionServiceServer(grpcServer, &RpcServer{})
	grpcServer.Serve(listener)
}
