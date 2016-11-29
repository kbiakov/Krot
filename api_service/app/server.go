package main

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
)

const (
	Routes_Api = "/api/v1"

	Routes_Auth = Routes_Api + "/auth"
	Routes_SignUp = Routes_Auth + "/signup"
	Routes_Login = Routes_Auth + "/login"
	Routes_Logout = Routes_Auth + "/logout"

	Routes_Users = Routes_Api + "/users/:uid"
	Routes_Logs = Routes_Users + "/logs"
	Routes_Status = Routes_Users + "/status"

	Routes_Receivers = Routes_Users + "/receivers"
	Routes_ReceiverSpec = Routes_Receivers + "/:name"

	Routes_Subscriptions = Routes_Users + "/subscriptions"
	Routes_SubscriptionSpec = Routes_Subscriptions + "/:id"
	Routes_SubscriptionStop = Routes_SubscriptionSpec + "/stop"
	Routes_SubscriptionResume = Routes_SubscriptionSpec + "/resume"
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
	mongo = session.DB("krot")

	// Create necessary indexes
	createIndexes()

	// Init API router
	e := echo.New()
	handleRoutes(e)
	e.Logger.Fatal(e.Start(":5000"))
}

func createIndexes()  {
	if err := mongo.C("users").EnsureIndex(mgo.Index{
		Key: []string{"email"},
		Unique: true,
		DropDups: true,
		Sparse: true,
	}); err != nil {
		panic(err)
	}

	if err := mongo.C("subscriptions").EnsureIndex(mgo.Index{
		Key: []string{"user_id"},
	}); err != nil {
		panic(err)
	}

	if err := mongo.C("logs").EnsureIndex(mgo.Index{
		Key: []string{"subscription_id"},
		Background: true,
	}); err != nil {
		panic(err)
	}
}

func handleRoutes(e *echo.Echo) {
	// Auth
	e.POST(Routes_SignUp, SignUp)
	e.POST(Routes_Login, Login)
	e.GET(Routes_Logout, Logout)

	// Users
	e.GET(Routes_Logs, GetLogs)
	e.GET(Routes_Status, GetJobsStatus)
	e.DELETE(Routes_Users, DeleteUser)

	// Receivers
	e.GET(Routes_Receivers, GetReceivers)
	e.POST(Routes_Receivers, CreateReceiver)
	e.DELETE(Routes_ReceiverSpec, RemoveReceiver)

	// Subscriptions
	e.GET(Routes_Subscriptions, GetSubscriptions)
	e.POST(Routes_Subscriptions, CreateSubscription)
	e.GET(Routes_SubscriptionStop, StopSubscription)
	e.GET(Routes_SubscriptionResume, ResumeSubscription)
	e.DELETE(Routes_SubscriptionSpec, RemoveSubscription)
}
