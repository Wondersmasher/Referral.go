package route

import (
	"github.com/Wondersmasher/Referral/controllers"
	"github.com/Wondersmasher/Referral/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAllRoutes(s *gin.Engine) {
	public := s.Group("/")
	public.POST("/sign-in", controllers.SignIn)
	public.POST("/sign-up", controllers.SignUp)

	protected := s.Group("/")
	protected.Use(middleware.AuthMiddleware)
	protected.POST("/sign-out", controllers.SignOut)
	protected.GET("/referrals/:referredBy", controllers.GetReferrals)
}
