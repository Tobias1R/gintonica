package localdb

import (
	"os"

	// DOT ENV
	"github.com/joho/godotenv"
	// mongo stuff(LOL)
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient mongo.Client

func getDbUrl() string {
	// Return the connection string for mongoDB
	godotenv.Load()
	url := "mongodb://"
	user := "root"
	password := "example"
	host := string(os.Getenv("MONGO_IP"))
	port := "27017"

	url += user
	url += ":"
	url += password
	url += "@"
	url += host
	url += ":"
	url += port

	return url
}

func Connect() (mongo.Client, error) {
	// Return a connected client or panic
	mongourl := getDbUrl()
	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI(mongourl))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	MongoClient = *client
	return *client, nil
}

func Db() (mongo.Database, error) {
	c, err := Connect()
	if err != nil {
		panic(err)
	}
	db := c.Database("store")
	return *db, nil
}

type ProductDb interface {
	collection() string
	list() []string
}

// Collections

type Product struct {
	ID        primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	Timestamp primitive.Timestamp `json:"timestamp"`
	Category  string              `json:"category"`
	Name      string              `json:"name"`
	Model     string              `json:"model"`
	Price     float64             `json:"price"`
	Brand     string              `json:"brand"`
}

type StorageUnit struct {
	ID        primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	Timestamp primitive.Timestamp `json:"timestamp"`
	SKU       string              `json:"sku"`
	Name      string              `json:"name"`
}

type SkuProduct struct {
	ID        primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	Timestamp primitive.Timestamp `json:"timestamp"`
	Sku       string              `json:"sku"`
	Product   string              `json:"product"`
	Available float64             `json:"available"`
}
