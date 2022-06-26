package api

import (
	// GIN
	"github.com/gin-gonic/gin"

	sec "github.com/Tobias1R/gintonica/api/auth"
	v1 "github.com/Tobias1R/gintonica/api/v1"
	mw "github.com/Tobias1R/gintonica/src/middleware"

	// swagger
	_ "github.com/Tobias1R/gintonica/docs/gintonica"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterAll(router gin.Engine) {

	authorized := router.Group("/")
	authorized.Use(mw.AuthorizeJWT())
	{
		authorized.GET("/products", v1.ProductsList)
		authorized.GET("/products/category/:category", v1.ProductCategoryList)
		authorized.GET("/products/view/:id", v1.ProductGet)
		authorized.PATCH("/products/update/:id", v1.ProductUpdate)
		authorized.DELETE("/products/:id", v1.ProductDelete)
		authorized.POST("/products", v1.ProductCreate)
	}
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.POST("/login", sec.JWTAuthenticate)
	router.GET("/tt/:msg", v1.TestQueue)
	router.GET("/task/:taskId", v1.TaskStatus)
	router.GET("/qc", v1.TaskQueue)
}
