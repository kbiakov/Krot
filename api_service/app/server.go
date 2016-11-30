package main

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

const jwtSecret = "secret"

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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
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

	r := e.Group(Routes_Users)
	r.Use(middleware.JWT([]byte(jwtSecret)))

	// Users
	r.GET(Routes_Logs, GetLogs)
	r.GET(Routes_Status, GetJobsStatus)
	r.DELETE(Routes_Users, DeleteUser)

	// Receivers
	r.GET(Routes_Receivers, GetReceivers)
	r.POST(Routes_Receivers, CreateReceiver)
	r.DELETE(Routes_ReceiverSpec, RemoveReceiver)

	// Subscriptions
	r.GET(Routes_Subscriptions, GetSubscriptions)
	r.POST(Routes_Subscriptions, CreateSubscription)
	r.GET(Routes_SubscriptionStop, StopSubscription)
	r.GET(Routes_SubscriptionResume, ResumeSubscription)
	r.DELETE(Routes_SubscriptionSpec, RemoveSubscription)
}
