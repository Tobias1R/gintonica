package service

import (
	"context"
	"log"

	"github.com/Tobias1R/gintonica/pkg/localdb"
	"github.com/Tobias1R/gintonica/pkg/security"
	"go.mongodb.org/mongo-driver/bson"
)

var AUTH_BACKEND string = "MONGO"

type LoginService interface {
	LoginUser(email string, password string) bool
}
type loginInformation struct {
	email    string
	password string
}

func StaticLoginService() LoginService {
	return &loginInformation{
		email:    "test@testing.test",
		password: "testing",
	}
}

func (info *loginInformation) LoginUser(pemail string, ppassword string) bool {
	if AUTH_BACKEND == "MONGO" {

		c, _ := localdb.Connect()
		defer c.Disconnect(context.TODO())

		users := c.Database("store").Collection("User")
		u := users.FindOne(context.TODO(), bson.M{"email": pemail})
		JSONData := struct {
			Email    string `bson:"email"`
			Password string `bson:"password"`
		}{}
		decodeError := u.Decode(&JSONData)
		if decodeError != nil {
			log.Println("Decode error: ", decodeError)
			return false
		}

		return security.CheckPasswordHash(ppassword, JSONData.Password)

	}
	return info.email == pemail && info.password == ppassword
}

func MongoLoginService() LoginService {

	AUTH_BACKEND = "MONGO"
	return &loginInformation{
		email:    "mongo",
		password: "mongo",
	}
}
