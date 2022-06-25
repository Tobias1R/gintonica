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
)

func serve() {
	// Load the .env file in the current directory
	godotenv.Load()

	router := gin.Default()

	api.RegisterAll(*router)

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

func main() {
	flag.Parse()

	var createAdmin bool = *createSuperUser

	if createAdmin {
		// highly secure user
		u, _ := sec.CreateSuperUser("Ozymandias", "email@gmail.com", "password")
		println("User " + u.Name + " created")
	}

	if *installFixturesDb {
		// initial buginganga
		installFixtures()
	}

	if !*noServer {
		// serve or not?
		defer serve()
	}

}
