package main

import (
	"github.com/bamzi/jobrunner"
	"gopkg.in/mgo.v2"

	"../rpc"
)

var mongo *mgo.Database

func main() {
	// Connect to MongoDB server
	session, err := mgo.Dial("mymongo:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Hold database context
	mongo = session.DB("krot").C("subscriptions")

	// Start job runner
	jobrunner.Start()
	defer jobrunner.Stop()

	// Start RPC server
	rpc.StartServer(9020)
}
