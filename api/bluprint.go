package api

import (
	// GIN
	"github.com/gin-gonic/gin"

	sec "github.com/Tobias1R/gintonica/api/auth"
	v1 "github.com/Tobias1R/gintonica/api/v1"
	mw "github.com/Tobias1R/gintonica/src/middleware"
)

func RegisterAll(router gin.Engine) {
	authorized := router.Group("/")
	authorized.Use(mw.AuthorizeJWT())
	{
		authorized.GET("/products", v1.ProductsList)
		authorized.GET("/products/:category", v1.ProductCategoryList)
	}

	router.POST("/login", sec.JWTAuthenticate)
}
