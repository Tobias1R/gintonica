package v1

import (
	"context"
	"fmt"
	"os"

	// GIN
	"net/http"

	"github.com/Tobias1R/gintonica/pkg/localdb"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"

	"log"
	"path/filepath"
)

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func ListDir(c *gin.Context) {

	var files []string

	root := "/home/ozy/Downloads"
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	data := []bson.M{{"files": files}}
	fmt.Println(data)
	c.IndentedJSON(http.StatusOK, data[0])

}

func ProductsList(c *gin.Context) {
	client, _ := localdb.Connect()
	defer client.Disconnect(context.TODO())
	productsCollection := client.Database("store").Collection("Product")
	// pipeline := []bson.M{
	// 	{"$convert": bson.M{"timestamp": bson.M{"to": "date"}}},
	// }
	cursor, err := productsCollection.Find(context.TODO(), bson.M{})
	// convert the cursor result to bson
	var results []localdb.Product
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, results)
}

func ProductCategoryList(c *gin.Context) {
	requestedCategory := c.Param("category")
	client, _ := localdb.Connect()
	defer client.Disconnect(context.TODO())
	productsCollection := client.Database("store").Collection("Product")
	cursor, err := productsCollection.Find(context.TODO(), bson.D{{"category", requestedCategory}})
	// convert the cursor result to bson
	var results []localdb.Product
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, results)
}
