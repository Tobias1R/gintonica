package api_auth

import (
	"context"
	"time"

	"github.com/Tobias1R/gintonica/pkg/localdb"
	sec "github.com/Tobias1R/gintonica/pkg/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	Timestamp primitive.Timestamp `json:"timestamp"`
	Name      string              `json:"name"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	Status    string              `json:"status"`
	LastLogin primitive.Timestamp `json:"lastLogin"`
	IsAdmin   bool                `json:"isAdmin"`
	Groups    []string            `json:"groups"`
}

var UserInterface interface {
	Save() User
	ChangePassword() bool
	Update(changePassword bool) bool
}

func (u User) Save() User {
	u.Password, _ = sec.HashPassword(u.Password)

	c, _ := localdb.Connect()
	defer c.Disconnect(context.TODO())

	uCollection := c.Database("store").Collection("User")
	uCollection.InsertOne(context.TODO(), u)

	return u
}

func (u User) ChangePassword(newPassword string) bool {
	c, _ := localdb.Connect()
	defer c.Disconnect(context.TODO())

	u.Password, _ = sec.HashPassword(newPassword)

	coll := c.Database("store").Collection("User")
	id := u.ID //primitive.ObjectIDFromHex("62b65c3034c47578a007e8ab")
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"password", u.Password}}}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return true
}

func (u User) Update(changePassword bool) bool {
	c, _ := localdb.Connect()
	defer c.Disconnect(context.TODO())

	if changePassword {
		u.Password, _ = sec.HashPassword(u.Password)
	}

	coll := c.Database("store").Collection("User")
	id := u.ID
	filter := bson.D{{"_id", id}}
	// updata := bson.D{
	// 	{"timestamp", primitive.Timestamp{T: uint32(time.Now().Unix())}},
	// 	{"name", u.Name},
	// 	{"email", u.Email},
	// 	{"password", u.Password},
	// 	{"status", u.Status},
	// 	{"lastLogin", primitive.Timestamp{T: uint32(time.Now().Unix())}},
	// }

	update := bson.D{{"$set", u}}
	r, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	println(r)
	return true

}

func ChangeStatus() {}

func CreateSuperUser(name string, email string, password string) (User, error) {
	// Fill struct and create an SuperUser in database
	u := User{
		ID:        primitive.NewObjectID(),
		Timestamp: primitive.Timestamp{T: uint32(time.Now().Unix())},
		Name:      name,
		Email:     email,
		Password:  password,
		Status:    "A",
		LastLogin: primitive.Timestamp{T: uint32(time.Now().Unix())},
		IsAdmin:   true,
		Groups:    []string{"admin"},
	}
	u.Save()
	return u, nil
}
