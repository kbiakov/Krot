package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	handleRoutes(e)
	e.Logger.Fatal(e.Start(":5000"))
}

func createIndexes() {
	if err := mongo.C("users").EnsureIndex(mgo.Index{
		Key:      []string{"email"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	}); err != nil {
		panic(err)
	}

	if err := mongo.C("subscriptions").EnsureIndex(mgo.Index{
		Key: []string{"user_id"},
	}); err != nil {
		panic(err)
	}

	if err := mongo.C("logs").EnsureIndex(mgo.Index{
		Key:        []string{"subscription_id"},
		Background: true,
	}); err != nil {
		panic(err)
	}
}

func handleRoutes(e *echo.Echo) {
	api := e.Group("/api/v1")

	// Auth
	a := e.Group("/auth")
	a.POST("/signup", SignUp)
	a.POST("/login", Login)
	a.GET("/logout", Logout)

	r := api.Group("/users/:uid")
	r.Use(middleware.JWT([]byte(jwtSecret)))
	r.Use(CheckUserAuth)

	// Users
	r.GET("/logs", GetLogs)
	r.GET("/status", GetJobsStatus)
	r.DELETE("", DeleteUser)

	// Receivers
	const Receivers = "/receivers"
	r.GET(Receivers, GetReceivers)
	r.POST(Receivers, CreateReceiver)
	r.DELETE(Receivers+"/:name", RemoveReceiver)

	// Subscriptions
	const Subscriptions = "/subscriptions"
	r.GET(Subscriptions, GetSubscriptions)
	r.POST(Subscriptions, CreateSubscription)

	const SubscriptionId = Subscriptions + "/:id"
	r.GET(SubscriptionId+"/stop", StopSubscription)
	r.GET(SubscriptionId+"/resume", ResumeSubscription)
	r.DELETE(SubscriptionId, RemoveSubscription)
}
