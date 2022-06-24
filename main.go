package main

import (

	// DOT ENV

	"github.com/joho/godotenv"

	// GIN
	"github.com/gin-gonic/gin"

	// mongo stuff(LOL)

	"go.mongodb.org/mongo-driver/mongo"

	// my libs
	"github.com/Tobias1R/gintonica/pkg/localdb"

	"github.com/Tobias1R/gintonica/api"
)

var mongoClient mongo.Client

func serve() {
	// Load the .env file in the current directory
	godotenv.Load()

	mongoClient, _ = localdb.Connect()
	//localdb.InstallFixtures(&mongoClient)
	//albumsCollection := client.Database("store").Collection("albums")

	//albumsCollection.InsertMany(context.TODO(), albums)

	router := gin.Default()

	api.RegisterAll(*router)

	router.Run("localhost:8080")
}

func main() {
	//u, _ := sec.CreateSuperUser("email@gmail.com", "password")
	//fmt.Println(u)
	defer serve()
}
