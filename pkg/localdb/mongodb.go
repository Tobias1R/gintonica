package localdb

import (
	"os"
	"time"

	// DOT ENV
	"github.com/joho/godotenv"
	// mongo stuff(LOL)
	"context"

	"go.mongodb.org/mongo-driver/bson"
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

type ProductInterface interface {
	Save() string   // objectid
	Update() string // object id
	Delete() bool
	GetProduct(id string) (Product, error)
}

func GetProduct(id string) (Product, error) {
	var p Product
	c, _ := Connect()
	defer c.Disconnect(context.TODO())

	coll := c.Database("store").Collection("Product")
	vid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": vid}

	err := coll.FindOne(context.TODO(), filter).Decode(&p)
	if err != nil {
		return p, err
	}

	return p, nil

}

func (p Product) Save() string {
	p.ID = primitive.NewObjectID()
	p.Timestamp = primitive.Timestamp{T: uint32(time.Now().Unix())}
	c, _ := Connect()
	defer c.Disconnect(context.TODO())

	uCollection := c.Database("store").Collection("Product")
	uCollection.InsertOne(context.TODO(), p)

	return p.ID.String()
}

func (p Product) Update() string {
	c, _ := Connect()
	defer c.Disconnect(context.TODO())

	coll := c.Database("store").Collection("Product")
	id := p.ID
	filter := bson.D{{"_id", id}}
	p.Timestamp = primitive.Timestamp{T: uint32(time.Now().Unix())}
	update := bson.D{{"$set", p}}
	r, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	println(r)
	return p.ID.String()
}

func (p Product) Delete() (bool, error) {
	c, _ := Connect()
	defer c.Disconnect(context.TODO())

	coll := c.Database("store").Collection("Product")
	id := p.ID
	filter := bson.D{{"_id", id}}
	r, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	println(r)
	return true, nil
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
