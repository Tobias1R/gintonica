package api_auth

import (
	"net/http"

	"github.com/Tobias1R/gintonica/pkg/controller"
	"github.com/Tobias1R/gintonica/pkg/service"

	"github.com/gin-gonic/gin"
)

// login godoc
// @Summary Do user login and return JWT token.
// @Description Login, got it?.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]string
// @Router /login [post]
// @Param email query string true "The email of the citizen"
// @Param password query string true "Citizen's password"
func JWTAuthenticate(ctx *gin.Context) {
	var loginService service.LoginService = service.MongoLoginService()
	var jwtService service.JWTService = service.JWTAuthService()
	var loginController controller.LoginController = controller.LoginHandler(loginService, jwtService)

	token := loginController.Login(ctx)
	if token != "" {
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	} else {
		ctx.JSON(http.StatusUnauthorized, nil)
	}
}
