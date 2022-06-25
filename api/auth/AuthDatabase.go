package api_auth

import (
	"context"
	"time"

	"github.com/Tobias1R/gintonica/pkg/localdb"
	sec "github.com/Tobias1R/gintonica/pkg/security"
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
	Update() User
}

func (u User) Save() User {
	u.Password, _ = sec.HashPassword(u.Password)

	c, _ := localdb.Connect()
	defer c.Disconnect(context.TODO())

	uCollection := c.Database("store").Collection("User")
	uCollection.InsertOne(context.TODO(), u)

	return u
}

func ChangePassword() {}

func Update() {}

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
