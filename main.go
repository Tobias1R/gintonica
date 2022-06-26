package main

import (
	//flag
	"flag"
	// DOT ENV
	"context"

	"github.com/joho/godotenv"
	// GIN
	"github.com/gin-gonic/gin"

	// App Libs
	"github.com/Tobias1R/gintonica/api"
	sec "github.com/Tobias1R/gintonica/api/auth"
	"github.com/Tobias1R/gintonica/pkg/localdb"
	"github.com/Tobias1R/gintonica/src/workers"
)

func serve() {
	// Load the .env file in the current directory
	godotenv.Load()
	// default router
	router := gin.Default()
	// api blueprint
	api.RegisterAll(*router)
	// vrum vrum
	router.Run("localhost:8080")
}

func installFixtures() {
	mongoClient, _ := localdb.Connect()
	defer mongoClient.Disconnect(context.TODO())
	localdb.InstallFixtures(&mongoClient)
}

var (
	installFixturesDb *bool
	createSuperUser   *bool
	noServer          *bool
)

func init() {
	installFixturesDb = flag.Bool("installFixturesDb", false, "--installFixturesDb=true")
	createSuperUser = flag.Bool("createSuperUser", false, "--createSuperUser=true")
	noServer = flag.Bool("noServer", false, "--noServer=true")
}

func startMQ() {
	//go mq.Consumer()
	//go mq.Publisher()
	//workers.StartWorkerDirectoryScan()
	const channelName string = "testao"
	w := workers.NewWorker(channelName, workers.TestME)
	t := workers.RunningTask{
		Order:   0,
		Channel: channelName,
		Status:  "PENDING",
		Data:    []byte(""),
	}

	w.Register(&t, channelName)
	go w.Start()
}

// @title Gin Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	flag.Parse()

	var createAdmin bool = *createSuperUser

	if createAdmin {
		// highly secure user
		u, _ := sec.CreateSuperUser("Ozymandias", "email@gmail.com", "password")
		// password change test
		if u.ChangePassword("another") {
			println("password changed")
		}
		// update test
		u.Name = "Jesus Junior"
		// LOGIN INFO
		u.Email = "j@j.sky"
		u.Password = "oddlypuertorican"
		// LOL-GIN INFO
		if u.Update(true) {
			println("name update")
		}
		println("User " + u.Email + " created " + " with password " + u.Password)
	}

	if *installFixturesDb {
		// initial buginganga
		installFixtures()
	}

	if !*noServer {
		// serve or not?
		startMQ()
		defer serve()
	}

}
