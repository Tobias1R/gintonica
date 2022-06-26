package v1

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	// GIN
	"net/http"

	"github.com/Tobias1R/gintonica/pkg/localdb"
	"github.com/Tobias1R/gintonica/src/mq"
	"github.com/Tobias1R/gintonica/src/workers"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"

	"log"
	"path/filepath"

	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
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

// ProductLIst godoc
// @Summary All products.
// @Description Returns a list of products
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} []localdb.Product{}
// @Router /products [get]
// @securitydefinitions.oauth2.application OAuth2Application
// @in Header
// @Param Authorization header string false "Bearer "
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

// ProductCategory godoc
// @Summary All products from this category.
// @Description Returns a list of products from this category
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} []localdb.Product{}
// @Router /products/category/{category} [get]
// @securitydefinitions.oauth2.application OAuth2Application
// @in Header
// @Param Authorization header string false "Bearer "
// @Param category path string true "The category you want"
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

// ProductGet godoc
// @Summary Retrieve Product document.
// @Description For real dude, it catchs the document that represents the Product.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]localdb.Product{}
// @Router /products/view/{id} [get]
// @securitydefinitions.oauth2.application OAuth2Application
// @in Header
// @Param Authorization header string false "Bearer "
// @Param id path string true "The id"
func ProductGet(c *gin.Context) {
	requestedId := c.Param("id")
	p, err := localdb.GetProduct(requestedId)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "object not found")
	} else {
		c.IndentedJSON(http.StatusOK, p)
	}
}

// ProductUpdate godoc
// @Summary Update Product document.
// @Description For real dude, it catchs the document that represents the Product, and update it.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]localdb.Product{}
// @Router /products/update/{id} [patch]
// @securitydefinitions.oauth2.application OAuth2Application
// @in Header
// @Param Authorization header string false "Bearer "
// @Param id path string true "The id"
// @Param localdb.Product{} body object true "The data"
func ProductUpdate(c *gin.Context) {
	requestedId := c.Param("id")
	p, err := localdb.GetProduct(requestedId)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "object not found")
	} else {
		err1 := c.ShouldBind(&p)
		if err1 != nil {
			c.IndentedJSON(http.StatusBadRequest, string(err1.Error()))
		} else {
			p.Update()
			c.IndentedJSON(http.StatusOK, p)
		}

	}
}

// ProductDelete godoc
// @Summary Delete Product document.
// @Description For real dude, it catchs the document that represents the Product, and update it.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {string} string
// @Router /products/{id} [delete]
// @securitydefinitions.oauth2.application OAuth2Application
// @in Header
// @Param Authorization header string false "Bearer "
// @Param id path string true "The id"
func ProductDelete(c *gin.Context) {
	requestedId := c.Param("id")
	p, err := localdb.GetProduct(requestedId)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "object not found")
	} else {
		err1 := c.ShouldBind(&p)
		if err1 != nil {
			c.IndentedJSON(http.StatusBadRequest, string(err1.Error()))
		} else {
			var msg string
			_, errDel := p.Delete()
			if errDel != nil {
				msg = "Deletion fail"
			} else {
				msg = "Delete OK!"
			}
			c.IndentedJSON(http.StatusOK, msg)
		}

	}
}

// ProductCreate godoc
// @Summary Create Product document.
// @Description Creates a product
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]localdb.Product{}
// @Router /products [post]
// @securitydefinitions.oauth2.application OAuth2Application
// @in Header
// @Param Authorization header string false "Bearer "
// @Param localdb.Product body object true "The data"
func ProductCreate(c *gin.Context) {
	var p localdb.Product
	err1 := c.ShouldBind(&p)
	println(&p)
	if err1 != nil {
		c.IndentedJSON(http.StatusBadRequest, string(err1.Error()))
	} else {
		msg := p.Save()
		c.IndentedJSON(http.StatusOK, msg)
	}
}

func ListDirCelery(c *gin.Context) {
	var dir string = "/home/ozy/Downloads"
	// create redis connection pool
	redisPool := &redis.Pool{
		MaxIdle:     3,                 // maximum number of idle connections in the pool
		MaxActive:   0,                 // maximum number of connections allocated by the pool at a given time
		IdleTimeout: 240 * time.Second, // close connections after remaining idle for this duration
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://redis:eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81@172.20.0.5:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// initialize celery client
	cli, _ := gocelery.NewCeleryClient(
		gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		1,
	)

	// prepare arguments
	taskName := "directoryScan"
	argA := dir
	argB := rand.Intn(10)

	// run task
	asyncResult, err := cli.DelayKwargs(
		taskName,
		map[string]interface{}{
			"directory": argA,
			"scanTime":  argB,
		},
	)
	if err != nil {
		panic(err)
	}

	// get results from backend with timeout
	res, err := asyncResult.Get(10 * time.Second)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, res)

}

func TestQueue(c *gin.Context) {
	msg := c.Param("msg")
	res := mq.Publisher(msg)
	c.IndentedJSON(http.StatusOK, res)
}

func TaskStatus(c *gin.Context) {
	taskId := c.Param("msg")
	rdb := workers.GetRedisClient()
	res, err := rdb.Get(context.Background(), taskId).Result()
	if err != nil {
		log.Printf("Message NOT stored in redis: %s", taskId)
	}
	var t workers.RunningTask
	t.UnmarshalBinary([]byte(res))
	c.IndentedJSON(http.StatusOK, t)
}
