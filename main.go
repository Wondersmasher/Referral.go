package main

import (
	"fmt"

	"github.com/Wondersmasher/Referral/env"
	mongodb "github.com/Wondersmasher/Referral/mongoDb"

	// redisCache "github.com/Wondersmasher/Referral/redisCache"
	"github.com/Wondersmasher/Referral/route"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	// Allow all origins (for development)

	env.Env()
	mongodb.InitDb()
	// redisCache.InitRedis()
	route.RegisterAllRoutes(server)

	fmt.Println("Server is running on port", env.PORT)
	server.Run(env.PORT)

}
