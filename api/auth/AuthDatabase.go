package api_auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Tobias1R/gintonica/pkg/localdb"
	sec "github.com/Tobias1R/gintonica/pkg/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	Timestamp primitive.Timestamp `json:"timestamp"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	Status    string              `json:"status"`
	LastLogin primitive.Timestamp `json:"lastLogin"`
	IsAdmin   bool                `json:"isAdmin"`
	Groups    []string            `json:"groups"`
}

var UserInterface interface {
	CreateUser() User
	ChangePassword() bool
	EditUser() User
}

func (u User) CreateUser() User {
	u.Password, _ = sec.HashPassword(u.Password)
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println("PAZUORDI", u.Password)
	c, _ := localdb.Connect()
	//client := localdb.MongoClient
	uCollection := c.Database("store").Collection("User")
	uCollection.InsertOne(context.TODO(), u)
	defer c.Disconnect(context.TODO())
	return u
}

func ChangePassword() {}

func EditUser() {}

func ChangeStatus() {}

func CreateSuperUser(email string, password string) (User, error) {
	//password, _ = sec.HashPassword(password)
	u := User{
		ID:        primitive.NewObjectID(),
		Timestamp: primitive.Timestamp{T: uint32(time.Now().Unix())},
		Email:     email,
		Password:  password,
		Status:    "A",
		LastLogin: primitive.Timestamp{T: uint32(time.Now().Unix())},
		IsAdmin:   false,
		Groups:    []string{"admin"},
	}
	u.CreateUser()
	return u, nil
}
