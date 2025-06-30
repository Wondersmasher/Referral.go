package route

import (
	"time"

	"github.com/Wondersmasher/Referral/controllers"
	"github.com/Wondersmasher/Referral/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterAllRoutes(s *gin.Engine) {

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	s.Use(cors.New(corsConfig))

	public := s.Group("/")
	public.POST("/sign-in", controllers.SignIn)
	public.POST("/sign-up", controllers.SignUp)

	protected := s.Group("/")
	protected.Use(middleware.AuthMiddleware)
	protected.POST("/sign-out", controllers.SignOut)
	protected.GET("/referrals/:referredBy", controllers.GetReferrals)
}
