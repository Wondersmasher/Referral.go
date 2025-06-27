package route

import (
	"github.com/Wondersmasher/Referral/controllers"
	"github.com/Wondersmasher/Referral/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAllRoutes(s *gin.Engine) {

	r := s.Group("/")
	r.Use(middleware.AuthMiddleware)
	r.POST("/sign-in", controllers.SignIn)
	r.POST("/sign-out", controllers.SignOut)
	r.POST("/sign-up", controllers.SignUp)
}
