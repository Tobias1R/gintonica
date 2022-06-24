package middleware

import (
	"fmt"
	"net/http"

	"github.com/Tobias1R/gintonica/pkg/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	fmt.Println("AUTH FUNCTION CALL")
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		fmt.Println("HEADER", authHeader)
		if authHeader == "" {
			fmt.Println("Unauthorized, blank header")
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			tokenString := authHeader[len(BEARER_SCHEMA):]
			token, err := service.JWTAuthService().ValidateToken(tokenString)
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				fmt.Println(claims)
			} else {
				fmt.Println(err)
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		}

	}
}
