package route

import (
	"github.com/Wondersmasher/Referral/controllers"
	"github.com/Wondersmasher/Referral/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAllRoutes(s *gin.Engine) {

	r := s.Group("/")
	r.Use(middleware.AuthMiddleware)
	r.GET("/", controllers.Test)

}
